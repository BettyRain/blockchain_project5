package p5

import (
	"net/http"
)

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
}

func Patient(w http.ResponseWriter, r *http.Request) {
	//View data by personal code
}

func Patients(w http.ResponseWriter, r *http.Request) {
	//View data from blocks (by doctor)
}

func AddData(w http.ResponseWriter, r *http.Request) {
	//Add patient data
}
