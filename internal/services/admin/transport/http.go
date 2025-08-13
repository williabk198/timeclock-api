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
		decodeCreateItemRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	rootRouter.Handle("/person/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().GetSpecific,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	rootRouter.Handle("/person/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().Update,
		decodeUpdateItemRequestData[endpoints.PersonData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

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

func decodeCreateItemRequestData[T any](ctx context.Context, r *http.Request) (reqData T, err error) {
	if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
		return reqData, fmt.Errorf("failed to parse data from request body: %w", err)
	}
	return reqData, nil
}

func decodeFetchItemRequestData(key string) func(context.Context, *http.Request) (string, error) {
	return func(ctx context.Context, r *http.Request) (string, error) {
		return mux.Vars(r)[key], nil
	}
}

func decodeUpdateItemRequestData[T any](key string) httputil.RequestDecoderFunc[endpoints.UpdateRequestData[T]] {
	return func(ctx context.Context, r *http.Request) (endpoints.UpdateRequestData[T], error) {
		id, _ := decodeFetchItemRequestData(key)(ctx, r)
		data, err := decodeCreateItemRequestData[T](ctx, r)
		if err != nil {
			return endpoints.UpdateRequestData[T]{}, err
		}

		return endpoints.UpdateRequestData[T]{
			ID:   id,
			Data: data,
		}, nil
	}
}

func encodeResponseBodyJSON[T any](ctx context.Context, w http.ResponseWriter, data T) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
