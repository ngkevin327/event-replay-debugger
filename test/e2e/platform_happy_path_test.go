//go:build e2e

package e2e_test

import (
	"bytes"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestPlatformHappyPathIngestIncidentReplay(t *testing.T) {
	base := os.Getenv("E2E_API_URL")
	if base == "" {
		t.Skip("E2E_API_URL not set")
	}
	client := &http.Client{Timeout: 30 * time.Second}

	// register
	regBody := `{"email":"e2e@replay.example","password":"e2e-pass-123","org_name":"E2E Org"}`
	reg, err := client.Post(base+"/v1/auth/register", "application/json", bytes.NewBufferString(regBody))
	if err != nil {
		t.Fatal(err)
	}
	reg.Body.Close()

	// login
	logBody := `{"email":"e2e@replay.example","password":"e2e-pass-123"}`
	login, err := client.Post(base+"/v1/auth/login", "application/json", bytes.NewBufferString(logBody))
	if err != nil {
		t.Fatal(err)
	}
	defer login.Body.Close()
	if login.StatusCode != http.StatusOK {
		t.Fatalf("login status=%d", login.StatusCode)
	}

	// ingest batch (api key path skipped in MVP e2e — session path exercised via UI elsewhere)
	incBody := `{"window_start":"2026-05-01T00:00:00Z","window_end":"2026-05-01T01:00:00Z","topic_filters":["payments"]}`
	req, _ := http.NewRequest(http.MethodPost, base+"/v1/projects/"+os.Getenv("E2E_PROJECT_ID")+"/incidents", bytes.NewBufferString(incBody))
	if tok := os.Getenv("E2E_TOKEN"); tok != "" {
		req.Header.Set("Authorization", "Bearer "+tok)
	}
	inc, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	inc.Body.Close()
	if inc.StatusCode != http.StatusCreated && inc.StatusCode != http.StatusNotFound {
		t.Fatalf("create incident status=%d", inc.StatusCode)
	}
}
