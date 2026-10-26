package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

// opaqueTokenByteLength is the number of random bytes used to generate an opaque token.
// 32 bytes (256 bits) of entropy is well beyond what is practical to brute force.
const opaqueTokenByteLength = 32

// GenerateOpaqueToken returns a cryptographically random, URL-safe opaque token string.
// The returned value carries no embedded information (unlike a JWT); it must be looked
// up against a datastore via HashToken to be validated.
func GenerateOpaqueToken() (string, error) {
	buf := make([]byte, opaqueTokenByteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate random token data: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashToken returns the hex-encoded SHA-256 hash of an opaque token value.
// Only the hash of a token is ever persisted to a datastore, so a compromised
// datastore can't be used to forge or replay tokens.
func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
