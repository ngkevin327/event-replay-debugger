package collector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCountEventsInWindow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("42\n"))
	}))
	defer srv.Close()

	c := NewCounter(srv.URL)
	n, err := c.CountEventsInWindow(context.Background(), "proj-1", time.Now().Add(-time.Hour), time.Now(), []string{"payments"})
	if err != nil {
		t.Fatal(err)
	}
	if n != 42 {
		t.Fatalf("got %d", n)
	}
}
