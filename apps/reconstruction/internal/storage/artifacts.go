package storage

import (
	"context"
	"fmt"
)

// ArtifactStore writes reconstruction artifacts.
type ArtifactStore interface {
	PutArtifact(ctx context.Context, key string, body []byte) (string, error)
}

// MemoryStore is an in-memory artifact store for tests.
type MemoryStore struct {
	objects map[string][]byte
}

// NewMemoryStore creates a test store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{objects: make(map[string][]byte)}
}

// PutArtifact stores bytes and returns an s3-style URI.
func (m *MemoryStore) PutArtifact(ctx context.Context, key string, body []byte) (string, error) {
	m.objects[key] = append([]byte(nil), body...)
	return fmt.Sprintf("s3://replay-artifacts/%s", key), nil
}
