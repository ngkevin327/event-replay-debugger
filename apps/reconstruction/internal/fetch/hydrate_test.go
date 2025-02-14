package fetch

import (
	"context"
	"testing"
)

type memStore struct {
	data map[string][]byte
}

func (m *memStore) GetPayload(_ context.Context, uri string) ([]byte, error) {
	return m.data[uri], nil
}

func TestHydratePool(t *testing.T) {
	rows := []EventRow{{S3URI: "s3://b/1"}, {S3URI: "s3://b/2"}}
	p := &HydratePool{
		Store: &memStore{data: map[string][]byte{
			"s3://b/1": []byte("a"),
			"s3://b/2": []byte("b"),
		}},
		Workers: 2,
	}
	if err := p.Hydrate(context.Background(), rows); err != nil {
		t.Fatal(err)
	}
	if string(rows[0].Payload) != "a" {
		t.Fatal()
	}
}
