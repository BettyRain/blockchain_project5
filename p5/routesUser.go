package p5

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerUser http.HandlerFunc
}

type Routes []Route

var routesUser = Routes{
	Route{
		"Patient",
		"GET",
		"/patient",
		Patient,
	},
}
