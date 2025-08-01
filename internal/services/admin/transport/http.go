package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/williabk198/timeclock/internal/services/admin/endpoints"
	"github.com/williabk198/timeclock/pkg/httputil"
)

func NewHttpHandler(adminEndpoints endpoints.Endpoints) http.Handler {
	rootRouter := mux.NewRouter()
	routeHandleBuilder := httputil.RouteHandleBuilder{
		ErrorEncoder: errorEncoder,
		ErrorHandler: errorHandler,
	}

	rootRouter.Handle("/person", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().Add,
		decodeRequest,
		encodeResponse,
	)).Methods(http.MethodPost)

	return rootRouter
}

func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError) // TODO: Base this and the returned message on `err`
	json.NewEncoder(w).Encode(struct {
		Error string
	}{
		err.Error(),
	})
}

func errorHandler(ctx context.Context, err error) {
	slog.ErrorContext(ctx, err.Error())
}

func decodeRequest(ctx context.Context, r *http.Request) (endpoints.PersonData, error) {
	var personData endpoints.PersonData
	if err := json.NewDecoder(r.Body).Decode(&personData); err != nil {
		return personData, fmt.Errorf("failed to parse data from request body: %w", err)
	}
	return personData, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, data endpoints.PersonData) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
