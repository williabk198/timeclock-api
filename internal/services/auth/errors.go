package auth

import "errors"

var (
	// ErrInvalidCredentials indicates the supplied email/password combination did not
	// match any enabled user.
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrAccountLocked indicates the user's account has been locked, typically due to
	// repeated failed login attempts.
	ErrAccountLocked = errors.New("account is locked")

	// ErrAccountDisabled indicates the user's account has been disabled and can no
	// longer authenticate.
	ErrAccountDisabled = errors.New("account is disabled")

	// ErrEmailInUse indicates an attempt to register a user with an email address that
	// is already associated with another user.
	ErrEmailInUse = errors.New("email address is already in use")

	// ErrInvalidToken indicates a provided token is malformed, unknown, expired,
	// revoked, or otherwise not usable.
	ErrInvalidToken = errors.New("token is invalid or expired")

	// ErrInvalidRequest indicates the caller supplied malformed or otherwise invalid
	// request data.
	ErrInvalidRequest = errors.New("invalid request data")
)
