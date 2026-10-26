package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateOpaqueToken(t *testing.T) {
	token1, err := GenerateOpaqueToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token1)

	token2, err := GenerateOpaqueToken()
	require.NoError(t, err)
	assert.NotEmpty(t, token2)

	// Each generated token should be unique.
	assert.NotEqual(t, token1, token2)
}

func TestHashToken(t *testing.T) {
	hash1 := HashToken("some-raw-token-value")
	hash2 := HashToken("some-raw-token-value")
	hash3 := HashToken("a-different-raw-token-value")

	// Hashing the same input should always produce the same, deterministic output.
	assert.Equal(t, hash1, hash2)
	// Hashing different input should produce different output.
	assert.NotEqual(t, hash1, hash3)
	// A hex-encoded SHA-256 hash is always 64 characters long.
	assert.Len(t, hash1, 64)
}
