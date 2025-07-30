package transport

import (
	"context"
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
	panic("unimplemented")
}

func errorHandler(ctx context.Context, err error) {
	panic("unimplemented")
}

func decodeRequest(ctx context.Context, r *http.Request) (endpoints.PersonData, error) {
	panic("unimplemented")
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, data endpoints.PersonData) error {
	panic("unimplemented")
}
