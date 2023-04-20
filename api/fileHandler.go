package api

import (
	logger "cyolo-efs/logging"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

func (api *API) uploadFile(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodPut {
		ttl, err := strconv.Atoi(req.Header.Get("File-TTL"))
		if err != nil {
			logger.Log.With("error", err.Error()).Error("error reading ttl header")
			// Set to default ttl time
			ttl = 1
		}

		file, header, err := req.FormFile("file")

		if err != nil {
			logger.Log.With("error", err.Error()).Error("error reading request body")
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			logger.Log.With("error", err.Error()).Error("error reading request body")
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		url := api.service.UploadFile(req.Context(), header.Filename, data, time.Minute*time.Duration(ttl))

		JSON(writer, map[string]interface{}{"url": url})
	} else {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (api *API) downloadFile(writer http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	if req.Method == http.MethodGet {
		path := filepath.Base(req.URL.Path)

		data, err := api.service.DownloadFile(req.Context(), path)

		if err != nil {
			logger.Log.With("error", err.Error()).Error("error file missing")
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		// set the Content-Disposition header to force a download
		writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", path))

		// write the file data to the response
		writer.Write(data)
	} else {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
