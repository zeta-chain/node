package openapi

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiFile   = "openapi.swagger.yaml"
	templateFile = "template.tpl"
)

//go:embed openapi.swagger.yaml
var staticFS embed.FS

//go:embed template.tpl
var templateFS embed.FS

func RegisterOpenAPIService(router *mux.Router) {
	router.Handle("/"+apiFile, http.FileServer(http.FS(staticFS)))
	router.HandleFunc("/", openAPIHandler())
}

func openAPIHandler() http.HandlerFunc {
	tmpl, err := template.ParseFS(templateFS, templateFile)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		err := tmpl.Execute(w, struct {
			URL string
		}{
			"/" + apiFile,
		})
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
