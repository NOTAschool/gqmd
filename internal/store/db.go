package store

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

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
	schema := `
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY,
		name TEXT UNIQUE NOT NULL,
		path TEXT NOT NULL,
		pattern TEXT DEFAULT '**/*.md',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := s.db.Exec(schema)
	return err
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

	status.TotalDocs = 0
	status.HasVectorIndex = false

	return status, nil
}
