package main

import (
	"./p3"
	"./p5"
	"log"
	"net/http"
	"os"
)

func main() {
	router := p3.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], router))
	} else {
		log.Fatal(http.ListenAndServe(":6686", router))
	}

	//TODO: how we can make two parallel routers running? Do a new main?
	routerUser := p5.NewRouter()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], routerUser))
	} else {
		log.Fatal(http.ListenAndServe(":6613", routerUser))
	}

}
