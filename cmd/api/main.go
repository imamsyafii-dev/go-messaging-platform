package main

import (
	"log"
	"net/http"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/router"
)

func main() {

	addr := ":8080"
	log.Printf("API server is listening on %s", addr)

	if err := http.ListenAndServe(addr, router.New()); err != nil {
		log.Fatal(err)
	}
}
