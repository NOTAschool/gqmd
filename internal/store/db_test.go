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
