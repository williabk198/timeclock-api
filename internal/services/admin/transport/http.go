package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/williabk198/timeclock/internal/models"
	"github.com/williabk198/timeclock/internal/services/admin/endpoints"
	"github.com/williabk198/timeclock/pkg/httputil"
)

func NewHttpHandler(adminEndpoints endpoints.Endpoints) http.Handler {
	rootRouter := mux.NewRouter()

	buildPersonRoutes(rootRouter.PathPrefix("/persons").Subrouter(), adminEndpoints.Person())
	buildPersonContactRoutes(rootRouter.PathPrefix("/persons/{personID}/contacts").Subrouter(), adminEndpoints.Contact())
	buildEmployeeEndpoints(rootRouter.PathPrefix("/employees").Subrouter(), adminEndpoints.Employee())

	return rootRouter
}

func buildPersonRoutes(personRouter *mux.Router, personEndpoints endpoints.PersonEndpoints) {
	routeHandleBuilder := httputil.RouteHandleBuilder{
		ErrorEncoder: errorEncoder,
		ErrorHandler: errorHandler,
	}

	personRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		personEndpoints.GetAll,
		decodeFetchAllRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		personEndpoints.Add,
		decodeCreateItemRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	personRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		personEndpoints.GetSpecific,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	personRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		personEndpoints.Update,
		decodeUpdateItemRequestData[endpoints.PersonData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)
}

func buildPersonContactRoutes(contactRouter *mux.Router, contactEndpoints endpoints.ContactEndpoints) {
	routeHandleBuilder := httputil.RouteHandleBuilder{
		ErrorEncoder: errorEncoder,
		ErrorHandler: errorHandler,
	}

	contactRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.GetPersonContacts,
		decodeFetchItemRequestData("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	contactRouter.Handle("/addresses", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.GetPersonContactAddresses,
		decodeFetchItemRequestData("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	contactRouter.Handle("/addresses", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.AddContactAddressForPerson,
		decodeAddSubItemRequestData[endpoints.PersonAddressData]("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	contactRouter.Handle("/addresses/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.UpdatePersonContactAddress,
		decodeUpdateContactRequestData[endpoints.PersonAddressData](),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	contactRouter.Handle("/addresses/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.DeleteContactAddressForPerson,
		decodeDeleteContactRequestData(),
		encodeResponseBodyJSON,
	)).Methods(http.MethodDelete)

	contactRouter.Handle("/emails", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.GetPersonContactEmails,
		decodeFetchItemRequestData("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	contactRouter.Handle("/emails", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.AddContactEmailForPerson,
		decodeAddSubItemRequestData[endpoints.PersonEmailData]("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	contactRouter.Handle("/emails/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.UpdatePersonContactEmail,
		decodeUpdateContactRequestData[endpoints.PersonEmailData](),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	contactRouter.Handle("/emails/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.DeleteContactEmailForPerson,
		decodeDeleteContactRequestData(),
		encodeResponseBodyJSON,
	)).Methods(http.MethodDelete)

	contactRouter.Handle("/phones", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.AddContactPhoneForPerson,
		decodeAddSubItemRequestData[endpoints.PersonPhoneData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPost)

	contactRouter.Handle("/phones", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.GetPersonContactPhones,
		decodeFetchItemRequestData("personID"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	contactRouter.Handle("/phones/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.UpdatePersonContactPhone,
		decodeUpdateContactRequestData[endpoints.PersonPhoneData](),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	contactRouter.Handle("/phones/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		contactEndpoints.DeleteContactPhoneForPerson,
		decodeDeleteContactRequestData(),
		encodeResponseBodyJSON,
	)).Methods(http.MethodDelete)
}

func buildEmployeeEndpoints(employeeRouter *mux.Router, employeeEndpoints endpoints.EmployeeEndpoints) {
	routeHandleBuilder := httputil.RouteHandleBuilder{
		ErrorEncoder: errorEncoder,
		ErrorHandler: errorHandler,
	}

	employeeRouter.Handle("", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.GetAll,
		decodeFetchAllRequestData,
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	employeeRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.GetSpecific,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodGet)

	employeeRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.Update,
		decodeUpdateItemRequestData[endpoints.EmployeeData]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	employeeRouter.Handle("/{id}", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.Delete,
		decodeFetchItemRequestData("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodDelete)

	employeeRouter.Handle("/{id}/exempt", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.UpdateExemptStatus,
		decodeUpdateItemRequestData[bool]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	employeeRouter.Handle("/{id}/pay", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.UpdatePay,
		decodeUpdateItemRequestData[models.EmployeePay]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	employeeRouter.Handle("/{id}/sickTime", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.UpdateSickTimeHours,
		decodeUpdateItemRequestData[float64]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	employeeRouter.Handle("/{id}/status", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.UpdateStatus,
		decodeUpdateItemRequestData[int]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)

	employeeRouter.Handle("/{id}/timeOff", httputil.BuildRouteHandler(
		routeHandleBuilder,
		employeeEndpoints.UpdateTimeOffHours,
		decodeUpdateItemRequestData[float64]("id"),
		encodeResponseBodyJSON,
	)).Methods(http.MethodPut)
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

func decodeDeleteContactRequestData() httputil.RequestDecoderFunc[endpoints.DeleteContactRequestData] {
	return func(ctx context.Context, r *http.Request) (endpoints.DeleteContactRequestData, error) {
		personID := mux.Vars(r)["personID"]
		contactID := mux.Vars(r)["id"]
		return endpoints.DeleteContactRequestData{
			PerosnID:  personID,
			ContactID: contactID,
		}, nil
	}
}

func decodeUpdateContactRequestData[T endpoints.ContactConstraint]() httputil.RequestDecoderFunc[endpoints.UpdateContactRequestData[T]] {
	return func(ctx context.Context, r *http.Request) (endpoints.UpdateContactRequestData[T], error) {
		personID := mux.Vars(r)["personID"]
		contactID := mux.Vars(r)["id"]
		data, err := decodeCreateItemRequestData[T](ctx, r)
		if err != nil {
			return endpoints.UpdateContactRequestData[T]{}, err
		}

		return endpoints.UpdateContactRequestData[T]{
			PersonID:  personID,
			ContactID: contactID,
			Data:      data,
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
