package main

import (
	"log"
	"net/http"
	"os"
	"strings"
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

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
