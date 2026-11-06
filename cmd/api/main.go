package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/config"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/conversation"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/handler"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/router"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/message"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/redis"
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

	redisClient, err := redis.Open(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	conversationRepository := conversation.NewRepository(db)
	conversationService := conversation.NewService(conversationRepository)
	conversationHandler := handler.NewConversationHandler(conversationService)
	tokenService := auth.NewTokenService(
		cfg.JWTSecret,
		24*time.Hour,
	)

	authHandler := handler.NewAuthHandler(
		userService,
		tokenService,
	)
	userHandler := handler.NewUserHandler(userService)

	messageRepository := message.NewRepository(db)
	messageService := message.NewService(
		messageRepository,
		conversationRepository,
	)
	messageHandler := handler.NewMessageHandler(messageService)

	addr := ":" + cfg.HTTPPort
	log.Printf("API server is listening on %s", addr)

	if err := http.ListenAndServe(
		addr,
		router.New(
			userHandler.Create,
			authHandler.Login,
			userHandler.Me,
			conversationHandler.Create,
			messageHandler.Create,
			messageHandler.List,
			tokenService,
		),
	); err != nil {
		log.Fatal(err)
	}
}
