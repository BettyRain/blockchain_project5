package p5

import (
	"fmt"
	"net/http"
)

func init() {
	// This function will be executed before everything else.
	// Do some initialization here.
}

// Register ID, download BlockChain, start HeartBeat
func Patient(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HEY IT IS WORKING")
}
