package main

import (
	"github.com/BentleyOph/htmx-go/pkg/routes"
	"github.com/gorilla/mux"

	// _ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
)

func main() {
	// It initializes the router, registers the book store routes, and starts the server.
	r := mux.NewRouter()
	routes.RegisterContactRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("localhost:9000", nil))
}
