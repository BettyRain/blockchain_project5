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
		"ListOfPatients",
		"GET",
		"/patients",
		Patients,
	},
	Route{
		"AddData",
		"GET",
		"/add",
		AddData,
	},
	Route{
		"PatientData",
		"GET",
		"/patient",
		Patient,
	},
	Route{
		"AddData",
		"POST",
		"/add",
		AddData,
	},
	Route{
		"SendToBlc",
		"GET",
		"/send",
		SendToMiners,
	},
}
