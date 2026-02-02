package store

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
)

// ScanResult holds scan statistics
type ScanResult struct {
	Added   int
	Updated int
	Removed int
	Errors  int
}

// ScanCollection scans a collection directory and indexes documents
func (s *Store) ScanCollection(name string) (*ScanResult, error) {
	col, err := s.GetCollection(name)
	if err != nil {
		return nil, err
	}

	result := &ScanResult{}

	// Walk directory and index files
	err = filepath.Walk(col.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors++
			return nil
		}

		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(col.Path, path)
		if err != nil {
			result.Errors++
			return nil
		}

		// Check if matches pattern (simple glob matching)
		if !matchGlob(col.Pattern, relPath) {
			return nil
		}

		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			result.Errors++
			return nil
		}

		// Calculate hash
		hash := hashContent(content)

		// Extract title from first line
		title := extractTitle(string(content), relPath)

		// Index document
		if err := s.IndexDocument(name, relPath, title, string(content), hash); err != nil {
			result.Errors++
			return nil
		}

		result.Added++
		return nil
	})

	return result, err
}

func hashContent(content []byte) string {
	h := sha256.Sum256(content)
	return hex.EncodeToString(h[:])
}

func extractTitle(content, fallback string) string {
	lines := strings.SplitN(content, "\n", 3)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	// Use filename without extension as fallback
	return strings.TrimSuffix(filepath.Base(fallback), filepath.Ext(fallback))
}

// matchGlob matches a path against a glob pattern
// Supports **/*.md style patterns
func matchGlob(pattern, path string) bool {
	// Handle **/*.ext pattern
	if strings.HasPrefix(pattern, "**/") {
		ext := strings.TrimPrefix(pattern, "**/")
		if strings.HasPrefix(ext, "*.") {
			suffix := strings.TrimPrefix(ext, "*")
			return strings.HasSuffix(path, suffix)
		}
	}

	// Handle *.ext pattern
	if strings.HasPrefix(pattern, "*.") {
		suffix := strings.TrimPrefix(pattern, "*")
		return strings.HasSuffix(path, suffix)
	}

	// Fallback to filepath.Match
	matched, _ := filepath.Match(pattern, filepath.Base(path))
	return matched
}
