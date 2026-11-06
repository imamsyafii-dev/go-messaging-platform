package main

import (
	"log"
	"net/http"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/config"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/router"
)

func main() {
	cfg := config.Load()
	addr := ":" + cfg.HTTPPort

	log.Printf("API server is listening on %s", addr)

	if err := http.ListenAndServe(addr, router.New()); err != nil {
		log.Fatal(err)
	}
}
