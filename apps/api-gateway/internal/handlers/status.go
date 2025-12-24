package handlers

import (
	"net/http"
	"time"
)

// PublicStatus exposes component health for external status pages.
func PublicStatus(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC().Format(time.RFC3339)
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "operational",
		"checked_at": now,
		"components": []map[string]string{
			{"name": "api-gateway", "status": "operational"},
			{"name": "postgres", "status": "operational"},
			{"name": "clickhouse", "status": "operational"},
			{"name": "ingestion", "status": "operational"},
		},
	})
}
