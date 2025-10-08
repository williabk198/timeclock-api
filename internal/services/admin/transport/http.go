package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

	personRouter := rootRouter.PathPrefix("/persons").Subrouter()

	personRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().GetAll,
		decodeFetchAllRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().Add,
		decodeCreateItemRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	personRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().GetSpecific,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Person().Update,
		decodeUpdateItemRequestData[endpoints.PersonData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	personRouter.Handle("/{id}/contacts", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().GetPersonContacts,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("/{id}/contacts/addresses", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().GetPersonContactAddresses,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("/{id}/contacts/addresses", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().AddContactAddressForPerson,
		decodeAddSubItemRequestData[endpoints.PersonAddressData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	personRouter.Handle("/{id}/contacts/emails", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().GetPersonContactEmails,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("/{id}/contacts/emails", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().AddContactEmailForPerson,
		decodeAddSubItemRequestData[endpoints.PersonEmailData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	personRouter.Handle("/{id}/contacts/phones", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().AddContactPhoneForPerson,
		decodeAddSubItemRequestData[endpoints.PersonPhoneData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	personRouter.Handle("/{id}/contacts/phones", httputil.BuildRouteHandler(
		routeHandleBuilder,
		adminEndpoints.Contact().GetPersonContactPhones,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

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

func decodeFetchAllRequestData(ctx context.Context, r *http.Request) (endpoints.GetPaginatedRequestData, error) {
	var err error
	queryParams := r.URL.Query()

	var offset int
	if offsetStr := queryParams.Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return endpoints.GetPaginatedRequestData{}, err
		}
	}

	if offset < 0 {
		return endpoints.GetPaginatedRequestData{}, fmt.Errorf("invalid offset value; must be greater than 0")
	}

	limit := 500 // If there isn't a limit provided, use 500 as the default value
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return endpoints.GetPaginatedRequestData{}, err
		}
	}

	if limit < 0 {
		return endpoints.GetPaginatedRequestData{}, fmt.Errorf("invalid limit value; must be greater than 0")
	}

	return endpoints.GetPaginatedRequestData{
		Offset: uint(offset),
		Limit:  uint(limit),
	}, nil
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

func decodeAddSubItemRequestData[T any](idKey string) httputil.RequestDecoderFunc[endpoints.AddSubItemRequestData[T]] {
	return func(ctx context.Context, r *http.Request) (endpoints.AddSubItemRequestData[T], error) {
		parentId := mux.Vars(r)[idKey]
		var reqData T
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			return endpoints.AddSubItemRequestData[T]{}, fmt.Errorf("failed to parse data from request body: %w", err)
		}
		return endpoints.AddSubItemRequestData[T]{
			ParentID: parentId,
			Data:     reqData,
		}, nil
	}
}

func encodeResponseBodyJSON[T any](ctx context.Context, w http.ResponseWriter, data T) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}
