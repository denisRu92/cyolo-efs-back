package api

import (
	logger "cyolo-efs/logging"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

// RequestLogger is used for standard logging
func (api *API) RequestLogger(next httprouter.Handle, name string) httprouter.Handle {
	return func(response http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		start := time.Now()

		next(response, request, ps)

		logger.Log.With("method", request.Method, "uri", request.RequestURI, "name", name, "duration", time.Since(start)).Info()
	}
}
