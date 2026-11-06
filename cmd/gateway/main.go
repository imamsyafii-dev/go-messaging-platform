package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coder/websocket"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/auth"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/gateway"
)

func main() {
	tokenService := auth.NewTokenService(
		getEnv("JWT_SECRET", "development-only-secret"),
		24*time.Hour,
	)

	hub := gateway.NewHub()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		parts := strings.SplitN(header, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}

		userID, err := tokenService.Parse(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			log.Printf("accept websocket: %v", err)
			return
		}
		defer conn.Close(websocket.StatusNormalClosure, "connection closed")

		hub.Add(userID, conn)
		defer hub.Remove(userID, conn)

		log.Printf("websocket connected user_id=%s", userID)

		for {
			messageType, data, err := conn.Read(r.Context())
			if err != nil {
				return
			}

			if err := conn.Write(r.Context(), messageType, data); err != nil {
				return
			}
		}
	})

	addr := ":" + getEnv("GATEWAY_PORT", "8081")
	log.Printf("WebSocket gateway is listening on %s", addr)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	server := &http.Server{
		Addr:    addr,
		Handler: nil,
	}

	go func() {
		log.Printf("WebSocket gateway is listening on %s", addr)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
		log.Println("shutting down WebSocket gateway")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("gateway shutdown error: %v", err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
