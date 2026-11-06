package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/config"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/outbox"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/redis"
)

func main() {
	cfg := config.Load()

	startupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(startupCtx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisClient, err := redis.Open(startupCtx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer redisClient.Close()

	repository := outbox.NewRepository(db)
	worker := outbox.NewWorker(repository, redisClient, 100)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	log.Println("outbox worker is running")

	worker.Run(ctx)

	log.Println("outbox worker stopped")
}
