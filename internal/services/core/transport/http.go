package transport

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/williabk198/timeclock/internal/services/core/endpoints"
)

func NewHttpHandler(coreEndpoints endpoints.Endpoints) http.Handler {
	rootRouter := mux.NewRouter()

	return rootRouter
}
