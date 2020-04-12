package main

import(
	"net/http"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {

	// Initialize router
	router := mux.NewRouter().StrictSlash(true)

	// Assign routes to handler functions
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
	return router
}