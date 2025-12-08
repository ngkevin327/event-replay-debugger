package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gwmw "github.com/replay/platform/apps/api-gateway/internal/middleware"
)

func TestShareTokenExpiredReturnsGone(t *testing.T) {
	token := "expired-token"
	shareMu.Lock()
	shareByTok[token] = shareRecord{
		IncidentID: "inc-1",
		OrgID:      "org-1",
		Scope:      gwmw.ShareScopeRead,
		ExpiresAt:  time.Now().Add(-time.Hour),
	}
	shareMu.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/v1/shared/incidents/inc-1?token="+token, nil)
	rec := httptest.NewRecorder()
	handler := gwmw.ShareTokenAuth(LookupShareToken)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rec, req)
	if rec.Code != http.StatusGone {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusGone)
	}
}

func TestShareTokenReadOnlyRejectsWrite(t *testing.T) {
	token := "valid-token"
	shareMu.Lock()
	shareByTok[token] = shareRecord{
		IncidentID: "inc-1",
		OrgID:      "org-1",
		Scope:      gwmw.ShareScopeRead,
		ExpiresAt:  time.Now().Add(time.Hour),
	}
	shareMu.Unlock()

	req := httptest.NewRequest(http.MethodPost, "/v1/shared/incidents/inc-1?token="+token, nil)
	rec := httptest.NewRecorder()
	chain := gwmw.ShareTokenAuth(LookupShareToken)(gwmw.RequireReadOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))
	chain.ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d want %d", rec.Code, http.StatusForbidden)
	}
}
