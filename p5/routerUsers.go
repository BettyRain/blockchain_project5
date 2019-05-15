package p5

import (
	"github.com/gorilla/mux"
	"net/http"
)

//Doctor's application
func NewRouterDoc() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routesDoc {
		var handler http.Handler
		handler = route.HandlerUser
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}

//Patient's application
func NewRouterPat() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routesPat {
		var handler http.Handler
		handler = route.HandlerUser
		handler = Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	return router
}
