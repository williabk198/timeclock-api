package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	errTestDecoder  = errors.New("test request decoder error")
	errTestEndpoint = errors.New("test endpoint logic error")
	errTestEncoder  = errors.New("test response encoder error")
)

func TestBuildRouteHandler(t *testing.T) {
	// To see if the HTTP handler function returned by BuildRouteHandler is correct,
	// we must also execute the returned HTTP handler and examine what was returned
	// if we expected a successful result. For test cases where we expect an error,
	// just simply checking whether the error exists is enough.

	type args[T, U any] struct {
		builder         RouteHandleBuilder
		endpointHandler EndpointLogicFunc[T, U]
		reqDecoder      RequestDecoderFunc[T]
		respEncoder     ResponseEncoderFunc[U]
	}
	type handlerArgs struct {
		w http.ResponseWriter
		r *http.Request
	}
	type wantResp struct {
		statusCode int
		body       []byte
	}

	tests := []struct {
		name        string
		args        args[testReqData, testRespData]
		handlerArgs handlerArgs
		wantResp    wantResp
		expectPanic bool
	}{
		{
			name: "Success",
			args: args[testReqData, testRespData]{
				builder: RouteHandleBuilder{
					ErrorEncoder: testErrorEncoder,
					ErrorHandler: testErrorHandler,
				},
				endpointHandler: testEndpointHandler,
				reqDecoder:      testReqDecoder,
				respEncoder:     testRespEncoder,
			},
			handlerArgs: handlerArgs{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("{\"val\":\"test\"}\n")),
			},
			wantResp: wantResp{
				statusCode: http.StatusOK,
				body:       []byte("{\"respVal\":\"test\"}\n"),
			},
		},
		{
			// Expect that the function will only return the first error that it encouters.
			// In this case, it expects that the ErrorEncoder is checked first...
			name:        "Error, Missing ErrorEncoder",
			expectPanic: true,
		},
		{
			// ... and then the ErrorHandler...
			name: "Error, Missing ErrorHandler",
			args: args[testReqData, testRespData]{
				builder: RouteHandleBuilder{
					ErrorEncoder: testErrorEncoder,
				},
			},
			expectPanic: true,
		},
		{
			// ... and then the EndpointHandlerFunc...
			name: "Error, Missing EndpointHandler",
			args: args[testReqData, testRespData]{
				builder: RouteHandleBuilder{
					ErrorEncoder: testErrorEncoder,
					ErrorHandler: testErrorHandler,
				},
			},
			expectPanic: true,
		},
		{
			// ... and then the RequestDecoder...
			name: "Error, Missing RequestDecoder",
			args: args[testReqData, testRespData]{
				builder: RouteHandleBuilder{
					ErrorEncoder: testErrorEncoder,
					ErrorHandler: testErrorHandler,
				},
				endpointHandler: testEndpointHandler,
			},
			expectPanic: true,
		},
		{
			// ... and then the ResponseEncoder...
			name: "Error, Missing ResponseEncoder",
			args: args[testReqData, testRespData]{
				builder: RouteHandleBuilder{
					ErrorEncoder: testErrorEncoder,
					ErrorHandler: testErrorHandler,
				},
				endpointHandler: testEndpointHandler,
				reqDecoder:      testReqDecoder,
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() {
					BuildRouteHandler(tt.args.builder, tt.args.endpointHandler, tt.args.reqDecoder, tt.args.respEncoder)
				})
				return
			}

			got := BuildRouteHandler(tt.args.builder, tt.args.endpointHandler, tt.args.reqDecoder, tt.args.respEncoder)
			got.ServeHTTP(tt.handlerArgs.w, tt.handlerArgs.r)
			resp := tt.handlerArgs.w.(*httptest.ResponseRecorder)

			gotBody, _ := io.ReadAll(resp.Body)

			assert.Equal(t, tt.wantResp.statusCode, resp.Code)
			assert.Equal(t, tt.wantResp.body, gotBody)
		})
	}
}

func TestRouteHandler_ServeHTTP(t *testing.T) {
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	type wantResp struct {
		statusCode int
		body       []byte
	}

	tests := []struct {
		name         string
		args         args
		routeHandler routeHandler[testReqData, testRespData]
		wantResp     wantResp
	}{
		{
			name: "Success",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("{\"val\":\"testing\"}\n")),
			},
			routeHandler: routeHandler[testReqData, testRespData]{
				endpoint: testEndpointHandler,
				decoder:  testReqDecoder,
				encoder:  testRespEncoder,
			},
			wantResp: wantResp{
				statusCode: http.StatusOK,
				body:       []byte("{\"respVal\":\"testing\"}\n"),
			},
		},
		{
			name: "Decode Error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/test", nil),
			},
			routeHandler: routeHandler[testReqData, testRespData]{
				decoder: func(ctx context.Context, r *http.Request) (testReqData, error) {
					return testReqData{}, errTestDecoder
				},
				errorEncoder: testErrorEncoder,
				errorHandler: func(ctx context.Context, err error) {},
			},
			wantResp: wantResp{
				statusCode: http.StatusBadRequest,
				body:       []byte(errTestDecoder.Error()),
			},
		},
		{
			name: "Endpoint Logic Error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("{}\n")),
			},
			routeHandler: routeHandler[testReqData, testRespData]{
				decoder: testReqDecoder,
				endpoint: func(ctx context.Context, reqData testReqData) (respData testRespData, err error) {
					err = errTestEndpoint
					return
				},
				errorEncoder: testErrorEncoder,
				errorHandler: func(ctx context.Context, err error) {},
			},
			wantResp: wantResp{
				statusCode: http.StatusInternalServerError,
				body:       []byte(errTestEndpoint.Error()),
			},
		},
		{
			name: "Encoder Error",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest(http.MethodGet, "/test", strings.NewReader("{\"val\":\"testing\"}\n")),
			},
			routeHandler: routeHandler[testReqData, testRespData]{
				decoder:  testReqDecoder,
				endpoint: testEndpointHandler,
				encoder: func(ctx context.Context, w http.ResponseWriter, trd testRespData) error {
					return errTestEncoder
				},
				errorEncoder: testErrorEncoder,
				errorHandler: func(ctx context.Context, err error) {},
			},
			wantResp: wantResp{
				statusCode: http.StatusNoContent,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.routeHandler.ServeHTTP(tt.args.w, tt.args.r)
			resp := tt.args.w.(*httptest.ResponseRecorder)

			assert.Equal(t, tt.wantResp.statusCode, resp.Code)
			assert.Equal(t, tt.wantResp.body, resp.Body.Bytes())
		})
	}
}

type testReqData struct {
	Value string `json:"val"`
}

type testRespData struct {
	RespValue string `json:"respVal"`
}

func testReqDecoder(ctx context.Context, r *http.Request) (testReqData, error) {
	var result testReqData

	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}

func testRespEncoder(ctx context.Context, w http.ResponseWriter, data testRespData) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(data)
}

func testErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	switch err {
	case errTestDecoder:
		w.WriteHeader(http.StatusBadRequest)
	case errTestEndpoint:
		w.WriteHeader(http.StatusInternalServerError)
	case errTestEncoder:
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}
	w.Write([]byte(err.Error()))
}

func testErrorHandler(ctx context.Context, err error) {
	slog.Default().ErrorContext(ctx, "test error", "error", err.Error())
}

func testEndpointHandler(ctx context.Context, reqData testReqData) (testRespData, error) {
	return testRespData{RespValue: reqData.Value}, nil
}
