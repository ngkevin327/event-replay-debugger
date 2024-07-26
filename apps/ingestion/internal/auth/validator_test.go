package auth_test

import (
	"testing"

	sharedauth "github.com/replay/platform/packages/shared-go/auth"
)

func TestValidatorKeyPrefixParsing(t *testing.T) {
	plain := "rk_live_" + "abcdef0123456789abcdef0123456789"
	prefix, err := sharedauth.ParseKeyPrefix(plain)
	if err != nil {
		t.Fatal(err)
	}
	if len(prefix) != 16 {
		t.Fatalf("prefix len %d", len(prefix))
	}
	hash, err := sharedauth.HashPrefix(plain)
	if err != nil || hash == "" {
		t.Fatalf("hash %v err %v", hash, err)
	}
}
