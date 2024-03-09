package routes

import (
	"github.com/BentleyOph/htmx-go/pkg/controllers" // Importing the controllers package
	"github.com/gorilla/mux"                        // mux is a package that implements a request router and dispatcher for matching incoming requests to their respective handler.
)

func RegisterContactRoutes(router *mux.Router) {
    router.HandleFunc("/contact/getform/{id}", controllers.UpdateForm).Methods("GET")
    router.HandleFunc("/contacts/new", controllers.CreateContact).Methods("POST", "GET")
    router.HandleFunc("/contacts/edit/{id}", controllers.UpdateContactByID).Methods("PUT")
    router.HandleFunc("/contacts/delete/{id}", controllers.DeleteContact).Methods("DELETE")
    router.HandleFunc("/contacts/{id}", controllers.GetContactById).Methods("GET")
    router.HandleFunc("/contacts", controllers.GetContacts).Methods("GET")
}