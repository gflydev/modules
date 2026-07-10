package jwt

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key-for-jwt-unit-tests"

func setupSecret(t *testing.T) {
	t.Helper()
	if err := os.Setenv("JWT_SECRET_KEY", testSecret); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	// A non-zero TTL so the generated access token is not immediately expired.
	if err := os.Setenv("JWT_TTL_MINUTES", "15"); err != nil {
		t.Fatalf("setenv: %v", err)
	}
}

func TestExtractTokenMetadata_RoundTrip(t *testing.T) {
	setupSecret(t)

	tokens, err := GenerateTokens("42", nil)
	if err != nil {
		t.Fatalf("GenerateTokens: %v", err)
	}

	meta, err := ExtractTokenMetadata(tokens.Access)
	if err != nil {
		t.Fatalf("ExtractTokenMetadata: %v", err)
	}
	if meta == nil {
		t.Fatal("expected non-nil metadata")
	}
	if meta.UserID != 42 {
		t.Errorf("UserID = %d, want 42", meta.UserID)
	}
	if meta.Expires == 0 {
		t.Error("expected non-zero Expires")
	}
}

func TestExtractTokenMetadata_WrongSecretRejected(t *testing.T) {
	setupSecret(t)

	tokens, err := GenerateTokens("42", nil)
	if err != nil {
		t.Fatalf("GenerateTokens: %v", err)
	}

	// Verify with a different secret must fail.
	if err = os.Setenv("JWT_SECRET_KEY", "another-secret"); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	if meta, err := ExtractTokenMetadata(tokens.Access); err == nil {
		t.Fatalf("expected error for token signed with a different secret, got meta=%v", meta)
	}
}

func TestExtractTokenMetadata_NoneAlgRejected(t *testing.T) {
	setupSecret(t)

	// Craft a token with alg=none — must be rejected by the pinned method list.
	tok := jwt.NewWithClaims(jwt.SigningMethodNone, jwt.MapClaims{
		"id":      "42",
		"expires": float64(1 << 40),
	})
	signed, err := tok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none: %v", err)
	}

	if meta, err := ExtractTokenMetadata(signed); err == nil {
		t.Fatalf("expected alg=none token to be rejected, got meta=%v", meta)
	}
}

func TestExtractTokenMetadata_MalformedClaimsNoPanic(t *testing.T) {
	setupSecret(t)

	// `id` as a number (not the expected string) must return an error, not panic.
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":      42,
		"expires": float64(1 << 40),
	})
	signed, err := tok.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	if _, err := ExtractTokenMetadata(signed); err == nil {
		t.Fatal("expected error for numeric id claim")
	}
}
