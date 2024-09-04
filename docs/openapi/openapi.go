package openapi

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiFile  = "openapi.swagger.yaml"
	htmlFile = "openapi.html"
)

//go:embed openapi.swagger.yaml
var staticFS embed.FS

//go:embed openapi.html
var html []byte

func RegisterOpenAPIService(router *mux.Router) {
	router.Handle("/"+apiFile, http.FileServer(http.FS(staticFS)))
	router.HandleFunc("/", openAPIHandler())
}

func openAPIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(html)
	}
}
