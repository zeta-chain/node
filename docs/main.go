/*
This file is intended solely for testing the OpenAPI Spec Browser, allowing
developers to explore the OpenAPI specification without the need to build and
start the chain.

To start the server, run `go run docs/main.go` and navigate to
http://localhost:8080 in your browser.
*/

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/zeta-chain/zetacore/docs/openapi"
)

func main() {
	router := mux.NewRouter()
	openapi.RegisterOpenAPIService(router)

	http.Handle("/", router)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Println("Starting server on :8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
