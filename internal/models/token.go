package models

import (
	"time"

	"github.com/google/uuid"
)

type TokenKind int

const (
	TokenKindAccess TokenKind = iota + 1
	TokenKindRefresh

	tokenKindStrAccess  string = "access"
	tokenKindStrRefresh string = "refresh"
)

func (tk TokenKind) String() string {
	switch tk {
	case TokenKindAccess:
		return tokenKindStrAccess
	case TokenKindRefresh:
		return tokenKindStrRefresh
	}

	return ""
}

// Token is the persisted record of an opaque token issued by the auth service.
// Only a SHA-256 hash of the raw token value is ever stored; the raw value exists
// only in memory at the moment it's issued, and can never be recovered from this record.
type Token struct {
	ID        uuid.UUID `jagsqlb:"id;omit"`
	UserID    uuid.UUID `jagsqlb:"user_id;omit-update"`
	Hash      string    `jagsqlb:"token_hash;omit-update"`
	Scope     string    `jagsqlb:"scope;omit-update"`
	Kind      TokenKind `jagsqlb:"kind;omit-update"`
	IssuedAt  time.Time `jagsqlb:"issued_at;omit-update"`
	ExpiresAt time.Time `jagsqlb:"expires_at;omit-update"`
	Revoked   bool      `jagsqlb:"revoked"`
}

// Active reports whether the token has neither expired nor been revoked as of `now`.
func (t Token) Active(now time.Time) bool {
	return !t.Revoked && now.Before(t.ExpiresAt)
}

// IssuedToken is the one-time representation of a newly issued opaque token returned to
// a client. Value holds the raw, unhashed token and must be treated as a secret; it is
// never persisted and can't be retrieved again once issued.
type IssuedToken struct {
	Value     string    `json:"token"`
	Scope     string    `json:"scope"`
	Kind      TokenKind `json:"kind"`
	IssuedAt  time.Time `json:"issued"`
	ExpiresAt time.Time `json:"expires"`
}
