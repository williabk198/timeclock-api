package httputil

import (
	"context"
	"net/http"

	"github.com/williabk198/timeclock/pkg/errors"
)

// RouteHandleBuilder might be a bit of a misnomer since it doesn't have a Build function,
// but it holds data that could be used across multiple http.Handlers.
// An example, if your creating a micro-service or just a colletion of related endpoints,
// then you probably would want to have the same error handling and error response for
// those endpoints. That is where this data type comes into play.
type RouteHandleBuilder struct {
	ErrorEncoder ErrorEncoderFunc
	ErrorHandler ErrorHandlerFunc
}

// BuildRouteHandler does what you think, it creates an http.Handler using the provided builder,
// endpoint handler, encode and decoder. This function isn't attached to RouteHandleBuilder type because
// of its need for type parameters which are used to hold the types of the request and response data types.
// In which that would have negatively effected
func BuildRouteHandler[T, U any](
	builder RouteHandleBuilder,
	endpointHandler EndpointLogicFunc[T, U],
	reqDecoder RequestDecoderFunc[T],
	respEncoder ResponseEncoderFunc[U],
) (http.Handler, error) {
	switch {
	case builder.ErrorEncoder == nil:
		return nil, errors.NewNilValueError("builder.ErrorEncoder")
	case builder.ErrorHandler == nil:
		return nil, errors.NewNilValueError("builder.ErrorHandler")
	case endpointHandler == nil:
		return nil, errors.NewNilValueError("endpointHandler")
	case reqDecoder == nil:
		return nil, errors.NewNilValueError("reqDecoder")
	case respEncoder == nil:
		return nil, errors.NewNilValueError("respEncoder")
	}

	return routeHandler[T, U]{
		endpoint:     endpointHandler,
		decoder:      reqDecoder,
		encoder:      respEncoder,
		errorHandler: builder.ErrorHandler,
		errorEncoder: builder.ErrorEncoder,
	}, nil
}

// routeHandler is heavily inspired by the github.com/go-kit/kit package implementation of its `http.Server` type.
// The major difference being that it uses generics to represent incoming and outgoing data types instead of using `any`
type routeHandler[T, U any] struct {
	endpoint     EndpointLogicFunc[T, U]
	decoder      RequestDecoderFunc[T]
	encoder      ResponseEncoderFunc[U]
	errorEncoder ErrorEncoderFunc
	errorHandler ErrorHandlerFunc
}

// ServeHTTP statisfies the http.Handler interface. It executes the RequestDecoderFunc, EndpointLogicFunc, and ResponseEncoderFunc
// of routeHandler and if any errors an encoutered along the way, they will be handled by the ErrorEncoderFunc and ErrorHandlerFunc
// that was also provided to the routeHandler.
func (rh routeHandler[T, U]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	data, err := rh.decoder(ctx, r)
	if err != nil {
		rh.errorHandler(ctx, err)
		rh.errorEncoder(ctx, err, w)
		return
	}

	respData, err := rh.endpoint(ctx, data)
	if err != nil {
		rh.errorHandler(ctx, err)
		rh.errorEncoder(ctx, err, w)
		return
	}

	err = rh.encoder(ctx, w, respData)
	if err != nil {
		rh.errorHandler(ctx, err)
		rh.errorEncoder(ctx, err, w)
		return
	}
}

// EndpointLogicFunc represents the core logic to be execuded by an HTTP handler.
type EndpointLogicFunc[T, U any] func(ctx context.Context, reqData T) (respData U, err error)

// RequestDecoderFunc represents logic that will take in the HTTP request and
// transforms the data in the request into data that can be consumed by an
// implementation of a EndpointLogicFunc later.
type RequestDecoderFunc[T any] func(context.Context, *http.Request) (T, error)

// ResponseEncoderFunc represents the logic that will take in data,
// presumably from an implementation of EndpointLogicFunc,
// and transforms, and then sends the transfomed data back to the client.
type ResponseEncoderFunc[T any] func(context.Context, http.ResponseWriter, T) error

// ErrorEncoderFunc represents the logic that will handle errors and handle a response
// back to the client related to the given error
type ErrorEncoderFunc func(context.Context, error, http.ResponseWriter)

// ErrorHandlerFunc represent the logic that will handle errors that occur while handling an HTTP request.
type ErrorHandlerFunc func(context.Context, error)
