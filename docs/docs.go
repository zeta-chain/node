package docs

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	apiFile   = "/static/openapi.swagger.yaml"
	indexFile = "template/index.tpl"
)

//go:embed static
var staticFS embed.FS

//go:embed template
var templateFS embed.FS

func RegisterOpenAPIService(router *mux.Router) {
	router.Handle(apiFile, http.FileServer(http.FS(staticFS)))
	router.HandleFunc("/", openAPIHandler())
}

func openAPIHandler() http.HandlerFunc {
	tmpl, err := template.ParseFS(templateFS, indexFile)
	if err != nil {
		panic(err)
	}

	return func(w http.ResponseWriter, req *http.Request) {
		err := tmpl.Execute(w, struct {
			URL string
		}{
			apiFile,
		})
		if err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
	}
}
