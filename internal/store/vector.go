package store

import (
	"encoding/binary"
	"math"
)

// Vector represents an embedding vector
type Vector []float32

// VectorResult holds vector search result
type VectorResult struct {
	Collection string
	Path       string
	Title      string
	Score      float64
	ChunkIdx   int
}

// vectorToBlob converts float32 slice to bytes
func vectorToBlob(v Vector) []byte {
	buf := make([]byte, len(v)*4)
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}

// blobToVector converts bytes to float32 slice
func blobToVector(b []byte) Vector {
	v := make(Vector, len(b)/4)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[i*4:]))
	}
	return v
}

// cosineSimilarity calculates cosine similarity between two vectors
func cosineSimilarity(a, b Vector) float64 {
	if len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// StoreEmbedding stores a vector embedding for a document
func (s *Store) StoreEmbedding(hash string, chunkIdx int, model string, vec Vector) error {
	blob := vectorToBlob(vec)
	now := nowISO()

	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO embeddings (hash, chunk_idx, model, dimensions, vector, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		hash, chunkIdx, model, len(vec), blob, now,
	)
	return err
}

// VectorSearch performs vector similarity search
func (s *Store) VectorSearch(queryVec Vector, limit int) ([]VectorResult, error) {
	if limit <= 0 {
		limit = 10
	}

	// Get all embeddings
	rows, err := s.db.Query(`
		SELECT e.hash, e.chunk_idx, e.vector, d.collection, d.path, d.title
		FROM embeddings e
		JOIN documents d ON d.hash = e.hash AND d.active = 1`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type scored struct {
		result VectorResult
		score  float64
	}
	var all []scored

	for rows.Next() {
		var hash string
		var chunkIdx int
		var blob []byte
		var col, path, title string

		if err := rows.Scan(&hash, &chunkIdx, &blob, &col, &path, &title); err != nil {
			continue
		}

		vec := blobToVector(blob)
		score := cosineSimilarity(queryVec, vec)

		all = append(all, scored{
			result: VectorResult{
				Collection: col,
				Path:       path,
				Title:      title,
				Score:      score,
				ChunkIdx:   chunkIdx,
			},
			score: score,
		})
	}

	// Sort by score descending
	for i := 0; i < len(all)-1; i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].score > all[i].score {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	// Return top results
	results := make([]VectorResult, 0, limit)
	for i := 0; i < len(all) && i < limit; i++ {
		results = append(results, all[i].result)
	}

	return results, nil
}
