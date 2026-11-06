package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/config"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/handler"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/router"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/user"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	addr := ":" + cfg.HTTPPort
	log.Printf("API server is listening on %s", addr)

	if err := http.ListenAndServe(
		addr,
		router.New(userHandler.Create),
	); err != nil {
		log.Fatal(err)
	}
}
