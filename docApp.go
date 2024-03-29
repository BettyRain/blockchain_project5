package main

import (
	"./p5"
	"log"
	"net/http"
	"os"
)

func main() {
	routerUser := p5.NewRouterDoc()
	if len(os.Args) > 1 {
		log.Fatal(http.ListenAndServe(":"+os.Args[1], routerUser))
	} else {
		log.Fatal(http.ListenAndServe(":8813", routerUser))
	}

}
