package endpoints

type RegisterUserData struct {
	PersonID string `json:"personID"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterUserResponse struct {
	ID string `json:"id"`
}

type LoginRequestData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenData struct {
	Token     string `json:"token"`
	Scope     string `json:"scope"`
	Kind      string `json:"kind"`
	IssuedAt  int64  `json:"issued"`  // UNIX Timestamp
	ExpiresAt int64  `json:"expires"` // UNIX Timestamp
}

type LoginResponseData struct {
	AccessToken  TokenData `json:"accessToken"`
	RefreshToken TokenData `json:"refreshToken"`
}

type RefreshRequestData struct {
	RefreshToken string `json:"refreshToken"`
}

type IntrospectRequestData struct {
	Token string `json:"token"`
}

// IntrospectResponseData mirrors the response shape defined by RFC 7662 (OAuth 2.0 Token
// Introspection), so other services can validate a token the same way they might
// already expect to from a standards-based introspection endpoint.
type IntrospectResponseData struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Subject   string `json:"sub,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

type RevokeRequestData struct {
	Token string `json:"token"`
}

type RevokeResponseData struct {
	Revoked bool `json:"revoked"`
}

type GrantScopeRequestData struct {
	UserID string
	Scope  string
}

type RevokeScopeRequestData struct {
	UserID string
	Scope  string
}

type ScopeResponseData struct {
	UserID string `json:"userID"`
	Scope  string `json:"scope"`
}
