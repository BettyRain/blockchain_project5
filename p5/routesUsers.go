package p5

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerUser http.HandlerFunc
}

type Routes []Route

var routesDoc = Routes{
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
		"AddData",
		"POST",
		"/add",
		AddData,
	},
	Route{
		"RegisterPatient",
		"GET",
		"/start",
		StartDoc,
	},
	Route{
		"HeartBeatReceive",
		"POST",
		"/heartbeat/receive",
		HeartBeatReceive,
	},
	Route{
		"Show",
		"GET",
		"/show",
		Show,
	},
}

var routesPat = Routes{
	Route{
		"PatientData",
		"GET",
		"/patient",
		Patient,
	},
	Route{
		"PatientData",
		"POST",
		"/patient",
		Patient,
	},
	Route{
		"RegisterPatient",
		"GET",
		"/start",
		StartPat,
	},
}
