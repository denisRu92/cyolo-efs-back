package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Health returns a JSON with the service status.
func (api *API) Health(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	JSON(w, map[string]interface{}{
		"status": "ok",
	})
}
