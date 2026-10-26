package authclient

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeIntrospector struct {
	result Introspection
	err    error
}

func (fi fakeIntrospector) Introspect(ctx context.Context, token string) (Introspection, error) {
	return fi.result, fi.err
}

func TestRequireToken(t *testing.T) {
	type args struct {
		introspector Introspector
		authHeader   string
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantCalled     bool
	}{
		{
			name: "Success",
			args: args{
				introspector: fakeIntrospector{result: Introspection{Active: true, Subject: "user-1"}},
				authHeader:   "Bearer valid-token",
			},
			wantStatusCode: http.StatusOK,
			wantCalled:     true,
		},
		{
			name: "MissingHeader",
			args: args{
				introspector: fakeIntrospector{result: Introspection{Active: true}},
				authHeader:   "",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "MalformedHeader",
			args: args{
				introspector: fakeIntrospector{result: Introspection{Active: true}},
				authHeader:   "NotBearer sometoken",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "InactiveToken",
			args: args{
				introspector: fakeIntrospector{result: Introspection{Active: false}},
				authHeader:   "Bearer expired-token",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name: "IntrospectionError",
			args: args{
				introspector: fakeIntrospector{err: assert.AnError},
				authHeader:   "Bearer valid-token",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			var gotSubject string
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				gotSubject, _ = Subject(r.Context())
				w.WriteHeader(http.StatusOK)
			})

			handler := RequireToken(tt.args.introspector)(next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.args.authHeader != "" {
				req.Header.Set("Authorization", tt.args.authHeader)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantCalled, called)
			if tt.wantCalled {
				assert.Equal(t, tt.args.introspector.(fakeIntrospector).result.Subject, gotSubject)
			}
		})
	}
}

func TestRequireScope(t *testing.T) {
	type args struct {
		introspector  Introspector
		requiredScope string
		authHeader    string
	}

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantCalled     bool
	}{
		{
			name: "HasRequiredScope",
			args: args{
				introspector:  fakeIntrospector{result: Introspection{Active: true, Scope: "user:admin"}},
				requiredScope: "user:admin",
				authHeader:    "Bearer valid-token",
			},
			wantStatusCode: http.StatusOK,
			wantCalled:     true,
		},
		{
			name: "HasRequiredScopeAmongMultiple",
			args: args{
				introspector:  fakeIntrospector{result: Introspection{Active: true, Scope: "default user:admin"}},
				requiredScope: "user:admin",
				authHeader:    "Bearer valid-token",
			},
			wantStatusCode: http.StatusOK,
			wantCalled:     true,
		},
		{
			name: "MissingRequiredScope",
			args: args{
				introspector:  fakeIntrospector{result: Introspection{Active: true, Scope: "default"}},
				requiredScope: "user:admin",
				authHeader:    "Bearer valid-token",
			},
			wantStatusCode: http.StatusForbidden,
		},
		{
			name: "InactiveToken",
			args: args{
				introspector:  fakeIntrospector{result: Introspection{Active: false, Scope: "user:admin"}},
				requiredScope: "user:admin",
				authHeader:    "Bearer expired-token",
			},
			wantStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				w.WriteHeader(http.StatusOK)
			})

			handler := RequireScope(tt.args.introspector, tt.args.requiredScope)(next)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", tt.args.authHeader)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantStatusCode, rec.Code)
			assert.Equal(t, tt.wantCalled, called)
		})
	}
}
