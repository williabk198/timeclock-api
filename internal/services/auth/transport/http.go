package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/williabk198/timeclock/internal/services/auth"
	"github.com/williabk198/timeclock/internal/services/auth/endpoints"
	"github.com/williabk198/timeclock/pkg/authclient"
	"github.com/williabk198/timeclock/pkg/httputil"
)

// requiredAdminScope is the token scope required to grant or revoke scopes on a user.
const requiredAdminScope = "user:admin"

func NewHttpHandler(authEndpoints endpoints.Endpoints) http.Handler {
	rootRouter := mux.NewRouter()
	buildAuthRoutes(rootRouter, authEndpoints)
	return rootRouter
}

// endpointIntrospector adapts endpoints.Endpoints.Introspect to satisfy
// authclient.Introspector, so the auth service's own admin-only routes can be
// protected the same way other services protect theirs (via
// authclient.RequireScope), without making a network round-trip to itself.
type endpointIntrospector struct {
	authEndpoints endpoints.Endpoints
}

func (ei endpointIntrospector) Introspect(ctx context.Context, token string) (authclient.Introspection, error) {
	resp, err := ei.authEndpoints.Introspect(ctx, endpoints.IntrospectRequestData{Token: token})
	if err != nil {
		return authclient.Introspection{}, err
	}

	return authclient.Introspection{
		Active:    resp.Active,
		Scope:     resp.Scope,
		TokenType: resp.TokenType,
		Subject:   resp.Subject,
		IssuedAt:  resp.IssuedAt,
		ExpiresAt: resp.ExpiresAt,
	}, nil
}

func buildAuthRoutes(router *mux.Router, authEndpoints endpoints.Endpoints) {
	routeHandleBuilder := httputil.RouteHandleBuilder{
		ErrorEncoder: errorEncoder,
		ErrorHandler: errorHandler,
	}

	router.Handle("/users", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.Register,
		decodeJSONBody[endpoints.RegisterUserData],
		encodeJSONResponse[endpoints.RegisterUserResponse](http.StatusCreated),
	)).Methods(http.MethodPost)

	router.Handle("/login", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.Login,
		decodeJSONBody[endpoints.LoginRequestData],
		encodeJSONResponse[endpoints.LoginResponseData](http.StatusOK),
	)).Methods(http.MethodPost)

	router.Handle("/introspect", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.Introspect,
		decodeJSONBody[endpoints.IntrospectRequestData],
		encodeJSONResponse[endpoints.IntrospectResponseData](http.StatusOK),
	)).Methods(http.MethodPost)

	router.Handle("/refresh", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.Refresh,
		decodeJSONBody[endpoints.RefreshRequestData],
		encodeJSONResponse[endpoints.LoginResponseData](http.StatusOK),
	)).Methods(http.MethodPost)

	router.Handle("/revoke", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.Revoke,
		decodeJSONBody[endpoints.RevokeRequestData],
		encodeJSONResponse[endpoints.RevokeResponseData](http.StatusOK),
	)).Methods(http.MethodPost)

	scopeRouter := router.PathPrefix("/users/{id}/scopes").Subrouter()
	scopeRouter.Use(authclient.RequireScope(endpointIntrospector{authEndpoints: authEndpoints}, requiredAdminScope))

	scopeRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.GrantScope,
		decodeGrantScopeRequestData,
		encodeJSONResponse[endpoints.ScopeResponseData](http.StatusCreated),
	)).Methods(http.MethodPost)

	scopeRouter.Handle("/{scope}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		authEndpoints.RevokeScope,
		decodeRevokeScopeRequestData,
		encodeJSONResponse[endpoints.ScopeResponseData](http.StatusOK),
	)).Methods(http.MethodDelete)
}

func decodeGrantScopeRequestData(ctx context.Context, r *http.Request) (endpoints.GrantScopeRequestData, error) {
	var body struct {
		Scope string `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return endpoints.GrantScopeRequestData{}, fmt.Errorf("%w: failed to parse data from request body: %v", auth.ErrInvalidRequest, err)
	}

	return endpoints.GrantScopeRequestData{
		UserID: mux.Vars(r)["id"],
		Scope:  body.Scope,
	}, nil
}

func decodeRevokeScopeRequestData(ctx context.Context, r *http.Request) (endpoints.RevokeScopeRequestData, error) {
	vars := mux.Vars(r)
	return endpoints.RevokeScopeRequestData{
		UserID: vars["id"],
		Scope:  vars["scope"],
	}, nil
}

func decodeJSONBody[T any](ctx context.Context, r *http.Request) (reqData T, err error) {
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return reqData, fmt.Errorf("%w: failed to parse data from request body: %v", auth.ErrInvalidRequest, err)
	}
	return reqData, nil
}

func encodeJSONResponse[T any](statusCode int) httputil.ResponseEncoderFunc[T] {
	return func(ctx context.Context, w http.ResponseWriter, data T) error {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		return json.NewEncoder(w).Encode(data)
	}
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	statusCode := http.StatusInternalServerError
	switch {
	case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, auth.ErrInvalidToken):
		statusCode = http.StatusUnauthorized
	case errors.Is(err, auth.ErrAccountLocked), errors.Is(err, auth.ErrAccountDisabled):
		statusCode = http.StatusForbidden
	case errors.Is(err, auth.ErrEmailInUse):
		statusCode = http.StatusConflict
	case errors.Is(err, auth.ErrInvalidRequest):
		statusCode = http.StatusBadRequest
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{
		err.Error(),
	})
}

func errorHandler(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
}
