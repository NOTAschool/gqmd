package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenPath(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.sqlite")

	s, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("OpenPath failed: %v", err)
	}
	defer s.Close()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Error("Database file was not created")
	}
}

func TestGetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.sqlite")

	s, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("OpenPath failed: %v", err)
	}
	defer s.Close()

	status, err := s.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status.DBPath != dbPath {
		t.Errorf("DBPath = %q, want %q", status.DBPath, dbPath)
	}
	if status.Collections != 0 {
		t.Errorf("Collections = %d, want 0", status.Collections)
	}
}

func TestCollectionCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.sqlite")

	s, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("OpenPath failed: %v", err)
	}
	defer s.Close()

	// Add collection
	err = s.AddCollection("test", tmpDir, "**/*.md")
	if err != nil {
		t.Fatalf("AddCollection failed: %v", err)
	}

	// List collections
	cols, err := s.ListCollections()
	if err != nil {
		t.Fatalf("ListCollections failed: %v", err)
	}
	if len(cols) != 1 {
		t.Errorf("ListCollections = %d, want 1", len(cols))
	}

	// Get collection
	col, err := s.GetCollection("test")
	if err != nil {
		t.Fatalf("GetCollection failed: %v", err)
	}
	if col.Name != "test" {
		t.Errorf("Name = %q, want %q", col.Name, "test")
	}

	// Remove collection
	err = s.RemoveCollection("test")
	if err != nil {
		t.Fatalf("RemoveCollection failed: %v", err)
	}

	cols, _ = s.ListCollections()
	if len(cols) != 0 {
		t.Errorf("After remove, ListCollections = %d, want 0", len(cols))
	}
}

func TestIndexAndSearch(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.sqlite")

	s, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("OpenPath failed: %v", err)
	}
	defer s.Close()

	// Index a document
	err = s.IndexDocument("docs", "test.md", "Test Document", "Hello world content", "abc123")
	if err != nil {
		t.Fatalf("IndexDocument failed: %v", err)
	}

	// Search
	results, err := s.Search("hello", 10)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("Search results = %d, want 1", len(results))
	}

	// Get document
	doc, content, err := s.Get("docs", "test.md")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if doc.Title != "Test Document" {
		t.Errorf("Title = %q, want %q", doc.Title, "Test Document")
	}
	if content != "Hello world content" {
		t.Errorf("Content mismatch")
	}
}
