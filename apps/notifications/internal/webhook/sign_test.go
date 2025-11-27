package webhook

import "testing"

func TestSignAndVerifyVectors(t *testing.T) {
	secret := "partner-secret"
	body := []byte(`{"event":"incident.ready","id":"inc-1"}`)
	sig := SignPayload(secret, body)
	if err := VerifySignature(secret, body, sig); err != nil {
		t.Fatal(err)
	}
	if err := VerifySignature(secret, body, "sha256=deadbeef"); err == nil {
		t.Fatal("expected invalid signature")
	}
}
