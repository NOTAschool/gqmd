package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Store struct {
	db     *sql.DB
	dbPath string
}

type Status struct {
	DBPath         string
	TotalDocs      int
	Collections    int
	HasVectorIndex bool
}

type Collection struct {
	ID        int64
	Name      string
	Path      string
	Pattern   string
	CreatedAt string
}

type Document struct {
	ID         int64
	Collection string
	Path       string
	Title      string
	Hash       string
	CreatedAt  string
	ModifiedAt string
	Active     bool
}

func getDBPath() (string, error) {
	cacheDir := os.Getenv("XDG_CACHE_HOME")
	if cacheDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		cacheDir = filepath.Join(home, ".cache")
	}
	dbDir := filepath.Join(cacheDir, "gqmd")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(dbDir, "index.sqlite"), nil
}

func Open() (*Store, error) {
	dbPath, err := getDBPath()
	if err != nil {
		return nil, fmt.Errorf("get db path: %w", err)
	}
	return OpenPath(dbPath)
}

func OpenPath(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	s := &Store{db: db, dbPath: dbPath}
	if err := s.init(); err != nil {
		db.Close()
		return nil, fmt.Errorf("init db: %w", err)
	}
	return s, nil
}

func (s *Store) init() error {
	// Collections table
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		pattern TEXT DEFAULT '**/*.md',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return err
	}

	// Content-addressable storage
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS content (
		hash TEXT PRIMARY KEY,
		doc TEXT NOT NULL,
		created_at TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	// Documents table
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS documents (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		collection TEXT NOT NULL,
		path TEXT NOT NULL,
		title TEXT NOT NULL,
		hash TEXT NOT NULL,
		created_at TEXT NOT NULL,
		modified_at TEXT NOT NULL,
		active INTEGER NOT NULL DEFAULT 1,
		FOREIGN KEY (hash) REFERENCES content(hash) ON DELETE CASCADE,
		UNIQUE(collection, path)
	)`)
	if err != nil {
		return err
	}

	// Indexes
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_documents_collection ON documents(collection, active)`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_documents_hash ON documents(hash)`)
	if err != nil {
		return err
	}

	// FTS5 virtual table
	_, err = s.db.Exec(`
	CREATE VIRTUAL TABLE IF NOT EXISTS documents_fts USING fts5(
		filepath, title, body,
		tokenize='porter unicode61'
	)`)
	if err != nil {
		return err
	}

	// Vector embeddings table
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS embeddings (
		hash TEXT NOT NULL,
		chunk_idx INTEGER NOT NULL DEFAULT 0,
		model TEXT NOT NULL,
		dimensions INTEGER NOT NULL,
		vector BLOB NOT NULL,
		created_at TEXT NOT NULL,
		PRIMARY KEY (hash, chunk_idx)
	)`)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) GetStatus() (*Status, error) {
	status := &Status{DBPath: s.dbPath}

	row := s.db.QueryRow("SELECT COUNT(*) FROM collections")
	if err := row.Scan(&status.Collections); err != nil {
		return nil, err
	}

	row = s.db.QueryRow("SELECT COUNT(*) FROM documents WHERE active = 1")
	if err := row.Scan(&status.TotalDocs); err != nil {
		return nil, err
	}

	status.HasVectorIndex = false
	return status, nil
}

// Collection management

func (s *Store) AddCollection(name, path, pattern string) error {
	if pattern == "" {
		pattern = "**/*.md"
	}
	_, err := s.db.Exec(
		`INSERT INTO collections (name, path, pattern) VALUES (?, ?, ?)`,
		name, path, pattern,
	)
	return err
}

func (s *Store) ListCollections() ([]Collection, error) {
	rows, err := s.db.Query(
		`SELECT id, name, path, pattern, created_at FROM collections ORDER BY name`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var c Collection
		if err := rows.Scan(&c.ID, &c.Name, &c.Path, &c.Pattern, &c.CreatedAt); err != nil {
			return nil, err
		}
		collections = append(collections, c)
	}
	return collections, rows.Err()
}

func (s *Store) GetCollection(name string) (*Collection, error) {
	row := s.db.QueryRow(
		`SELECT id, name, path, pattern, created_at FROM collections WHERE name = ?`,
		name,
	)
	var c Collection
	if err := row.Scan(&c.ID, &c.Name, &c.Path, &c.Pattern, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (s *Store) RemoveCollection(name string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Delete documents in collection
	_, err = tx.Exec(`DELETE FROM documents WHERE collection = ?`, name)
	if err != nil {
		return err
	}

	// Delete collection
	result, err := tx.Exec(`DELETE FROM collections WHERE name = ?`, name)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("collection %q not found", name)
	}

	return tx.Commit()
}

// Document indexing

func (s *Store) IndexDocument(collection, path, title, content, hash string) error {
	now := nowISO()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert content (ignore if exists)
	_, err = tx.Exec(
		`INSERT OR IGNORE INTO content (hash, doc, created_at) VALUES (?, ?, ?)`,
		hash, content, now,
	)
	if err != nil {
		return err
	}

	// Upsert document
	_, err = tx.Exec(`
		INSERT INTO documents (collection, path, title, hash, created_at, modified_at, active)
		VALUES (?, ?, ?, ?, ?, ?, 1)
		ON CONFLICT(collection, path) DO UPDATE SET
			title = excluded.title,
			hash = excluded.hash,
			modified_at = excluded.modified_at,
			active = 1`,
		collection, path, title, hash, now, now,
	)
	if err != nil {
		return err
	}

	// Get document ID for FTS
	var docID int64
	err = tx.QueryRow(
		`SELECT id FROM documents WHERE collection = ? AND path = ?`,
		collection, path,
	).Scan(&docID)
	if err != nil {
		return err
	}

	// Update FTS index
	filepath := collection + "/" + path
	_, err = tx.Exec(`DELETE FROM documents_fts WHERE rowid = ?`, docID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		`INSERT INTO documents_fts (rowid, filepath, title, body) VALUES (?, ?, ?, ?)`,
		docID, filepath, title, content,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Search result types

type SearchResult struct {
	Collection string
	Path       string
	Title      string
	Snippet    string
	Score      float64
}

// Search performs FTS5 full-text search
func (s *Store) Search(query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}

	rows, err := s.db.Query(`
		SELECT d.collection, d.path, d.title,
			snippet(documents_fts, 2, '<mark>', '</mark>', '...', 32) as snippet,
			bm25(documents_fts) as score
		FROM documents_fts f
		JOIN documents d ON d.id = f.rowid
		WHERE documents_fts MATCH ? AND d.active = 1
		ORDER BY score
		LIMIT ?`,
		query, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	for rows.Next() {
		var r SearchResult
		if err := rows.Scan(&r.Collection, &r.Path, &r.Title, &r.Snippet, &r.Score); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

// Get retrieves a document by collection and path
func (s *Store) Get(collection, path string) (*Document, string, error) {
	row := s.db.QueryRow(`
		SELECT d.id, d.collection, d.path, d.title, d.hash, d.created_at, d.modified_at, d.active, c.doc
		FROM documents d
		JOIN content c ON c.hash = d.hash
		WHERE d.collection = ? AND d.path = ? AND d.active = 1`,
		collection, path,
	)

	var doc Document
	var content string
	var active int
	err := row.Scan(&doc.ID, &doc.Collection, &doc.Path, &doc.Title, &doc.Hash,
		&doc.CreatedAt, &doc.ModifiedAt, &active, &content)
	if err != nil {
		return nil, "", err
	}
	doc.Active = active == 1
	return &doc, content, nil
}

// MultiGetResult holds document with content
type MultiGetResult struct {
	Document *Document
	Content  string
}

// MultiGet retrieves multiple documents by paths
func (s *Store) MultiGet(paths []string, maxBytes int) ([]MultiGetResult, error) {
	if maxBytes <= 0 {
		maxBytes = 10 * 1024 // 10KB default
	}

	var results []MultiGetResult
	totalBytes := 0

	for _, p := range paths {
		// Parse collection/path format
		collection, docPath := splitPath(p)
		if collection == "" {
			continue
		}

		doc, content, err := s.Get(collection, docPath)
		if err != nil {
			continue // Skip missing documents
		}

		contentLen := len(content)
		if totalBytes+contentLen > maxBytes {
			// Truncate content to fit
			remaining := maxBytes - totalBytes
			if remaining > 0 {
				content = content[:remaining] + "\n... (truncated)"
			} else {
				break
			}
		}

		results = append(results, MultiGetResult{
			Document: doc,
			Content:  content,
		})
		totalBytes += contentLen

		if totalBytes >= maxBytes {
			break
		}
	}

	return results, nil
}

func splitPath(p string) (collection, path string) {
	idx := strings.Index(p, "/")
	if idx < 0 {
		return "", ""
	}
	return p[:idx], p[idx+1:]
}
