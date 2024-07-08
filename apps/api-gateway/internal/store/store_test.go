package store_test

import (
	"testing"

	"github.com/replay/platform/apps/api-gateway/internal/store"
)

func TestNewStoreRequiresPool(t *testing.T) {
	s := store.NewStore(nil)
	if s == nil {
		t.Fatal("expected non-nil store wrapper")
	}
}
