package authclient

import (
	"context"
	"net/http"
	"slices"
	"strings"
)

type contextKey string

const (
	subjectContextKey contextKey = "authclient.subject"
	scopeContextKey   contextKey = "authclient.scope"
)

// Introspector is satisfied by Client; it's defined here so tests can supply a fake.
type Introspector interface {
	Introspect(ctx context.Context, token string) (Introspection, error)
}

// RequireToken returns middleware that validates the bearer token on incoming requests by
// calling the given Introspector. Requests with a missing, malformed, or inactive token
// are rejected with 401 Unauthorized. On success, the token's subject (user ID) and scope
// are attached to the request context and can be retrieved with Subject and Scope.
func RequireToken(introspector Introspector) func(http.Handler) http.Handler {
	return requireToken(introspector, "")
}

// RequireScope returns middleware like RequireToken, but additionally requires the token to
// carry the given scope among its space-separated list of scopes. A request with an active
// token that doesn't carry the required scope is rejected with 403 Forbidden.
func RequireScope(introspector Introspector, scope string) func(http.Handler) http.Handler {
	return requireToken(introspector, scope)
}

func requireToken(introspector Introspector, requiredScope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}

			result, err := introspector.Introspect(r.Context(), token)
			if err != nil || !result.Active {
				http.Error(w, "invalid or expired token", http.StatusUnauthorized)
				return
			}

			if requiredScope != "" && !hasScope(result.Scope, requiredScope) {
				http.Error(w, "token does not have the required scope", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), subjectContextKey, result.Subject)
			ctx = context.WithValue(ctx, scopeContextKey, result.Scope)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Subject returns the user ID attached to the request context by RequireToken or
// RequireScope, if any.
func Subject(ctx context.Context) (string, bool) {
	subject, ok := ctx.Value(subjectContextKey).(string)
	return subject, ok
}

// Scope returns the token's scope attached to the request context by RequireToken or
// RequireScope, if any.
func Scope(ctx context.Context) (string, bool) {
	scope, ok := ctx.Value(scopeContextKey).(string)
	return scope, ok
}

// hasScope reports whether required is present among the space-separated scopes in tokenScope.
func hasScope(tokenScope, required string) bool {
	return slices.Contains(strings.Fields(tokenScope), required)
}

func bearerToken(r *http.Request) string {
	const prefix = "Bearer "

	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, prefix) {
		return ""
	}

	return strings.TrimPrefix(header, prefix)
}
