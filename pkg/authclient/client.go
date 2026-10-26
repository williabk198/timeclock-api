package authclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Introspection mirrors the response of the auth service's "/introspect" endpoint, which
// follows the shape defined by RFC 7662 (OAuth 2.0 Token Introspection).
type Introspection struct {
	Active    bool   `json:"active"`
	Scope     string `json:"scope,omitempty"`
	TokenType string `json:"token_type,omitempty"`
	Subject   string `json:"sub,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

// Client calls the auth service's token introspection endpoint to validate the opaque
// tokens it issues. Since those tokens are opaque, this is the only way for another
// service to determine whether a token is active and who it belongs to.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient returns a Client that calls the auth service running at baseURL (e.g.
// "http://auth:8080") to introspect tokens.
func NewClient(baseURL string) Client {
	return Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Introspect asks the auth service whether the given raw token value is currently active,
// returning its associated scope, kind, and subject (user ID) if so.
func (c Client) Introspect(ctx context.Context, token string) (Introspection, error) {
	reqBody, err := json.Marshal(struct {
		Token string `json:"token"`
	}{Token: token})
	if err != nil {
		return Introspection{}, fmt.Errorf("failed to build introspection request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/introspect", bytes.NewReader(reqBody))
	if err != nil {
		return Introspection{}, fmt.Errorf("failed to create introspection request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return Introspection{}, fmt.Errorf("failed to call introspection endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Introspection{}, fmt.Errorf("introspection endpoint returned unexpected status: %s", resp.Status)
	}

	var result Introspection
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return Introspection{}, fmt.Errorf("failed to decode introspection response: %w", err)
	}

	return result, nil
}
