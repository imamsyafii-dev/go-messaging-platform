package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/middleware"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/message"
)

type MessageService interface {
	Create(
		ctx context.Context,
		conversationID pgtype.UUID,
		senderID pgtype.UUID,
		clientMessageID pgtype.UUID,
		messageType string,
		content string,
	) (sqlc.Message, error)

	List(
		ctx context.Context,
		conversationID pgtype.UUID,
		userID pgtype.UUID,
		limit int32,
	) ([]sqlc.Message, error)

	ListBefore(
		ctx context.Context,
		conversationID pgtype.UUID,
		userID pgtype.UUID,
		cursorCreatedAt pgtype.Timestamptz,
		cursorID pgtype.UUID,
		limit int32,
	) ([]sqlc.Message, error)
}

type MessageHandler struct {
	service MessageService
}

type messageHistoryResponse struct {
	Messages []messageResponse `json:"messages"`
}

func NewMessageHandler(service MessageService) *MessageHandler {
	return &MessageHandler{
		service: service,
	}
}

type createMessageRequest struct {
	ClientMessageID string `json:"client_message_id"`
	Type            string `json:"type"`
	Content         string `json:"content"`
}

type messageResponse struct {
	ID              string `json:"id"`
	ConversationID  string `json:"conversation_id"`
	SenderID        string `json:"sender_id"`
	ClientMessageID string `json:"client_message_id"`
	Type            string `json:"type"`
	Content         string `json:"content,omitempty"`
}

func (h *MessageHandler) Create(w http.ResponseWriter, r *http.Request) {
	conversationID, err := uuid.Parse(chi.URLParam(r, "conversationID"))
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	senderIDString, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "user identity not found", http.StatusUnauthorized)
		return
	}

	senderID, err := uuid.Parse(senderIDString)
	if err != nil {
		http.Error(w, "invalid user identity", http.StatusUnauthorized)
		return
	}

	var request createMessageRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	clientMessageID, err := uuid.Parse(request.ClientMessageID)
	if err != nil {
		http.Error(w, "invalid client_message_id", http.StatusBadRequest)
		return
	}

	createdMessage, err := h.service.Create(
		r.Context(),
		pgtype.UUID{Bytes: conversationID, Valid: true},
		pgtype.UUID{Bytes: senderID, Valid: true},
		pgtype.UUID{Bytes: clientMessageID, Valid: true},
		request.Type,
		request.Content,
	)
	if err != nil {
		switch {
		case errors.Is(err, message.ErrInvalidMessageType),
			errors.Is(err, message.ErrEmptyMessage):
			http.Error(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, message.ErrNotConversationMember):
			http.Error(w, err.Error(), http.StatusForbidden)
		default:
			http.Error(w, "failed to create message", http.StatusInternalServerError)
		}

		return
	}

	response := messageResponse{
		ID:              createdMessage.ID.String(),
		ConversationID:  createdMessage.ConversationID.String(),
		SenderID:        createdMessage.SenderID.String(),
		ClientMessageID: createdMessage.ClientMessageID.String(),
		Type:            createdMessage.Type,
		Content:         createdMessage.Content.String,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

func (h *MessageHandler) List(w http.ResponseWriter, r *http.Request) {
	conversationID, err := uuid.Parse(chi.URLParam(r, "conversationID"))
	if err != nil {
		http.Error(w, "invalid conversation ID", http.StatusBadRequest)
		return
	}

	userIDString, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "user identity not found", http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		http.Error(w, "invalid user identity", http.StatusUnauthorized)
		return
	}

	limit := int32(50)
	if value := r.URL.Query().Get("limit"); value != "" {
		parsedLimit, parseErr := strconv.ParseInt(value, 10, 32)
		if parseErr != nil {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = int32(parsedLimit)
	}

	var messages []sqlc.Message

	cursorCreatedAt := r.URL.Query().Get("cursor_created_at")
	cursorID := r.URL.Query().Get("cursor_id")

	if cursorCreatedAt == "" && cursorID == "" {
		messages, err = h.service.List(
			r.Context(),
			pgtype.UUID{Bytes: conversationID, Valid: true},
			pgtype.UUID{Bytes: userID, Valid: true},
			limit,
		)
	} else {
		if cursorCreatedAt == "" || cursorID == "" {
			http.Error(w, "both cursor fields are required", http.StatusBadRequest)
			return
		}

		parsedTime, parseErr := time.Parse(time.RFC3339Nano, cursorCreatedAt)
		if parseErr != nil {
			http.Error(w, "invalid cursor_created_at", http.StatusBadRequest)
			return
		}

		parsedCursorID, parseErr := uuid.Parse(cursorID)
		if parseErr != nil {
			http.Error(w, "invalid cursor_id", http.StatusBadRequest)
			return
		}

		messages, err = h.service.ListBefore(
			r.Context(),
			pgtype.UUID{Bytes: conversationID, Valid: true},
			pgtype.UUID{Bytes: userID, Valid: true},
			pgtype.Timestamptz{Time: parsedTime, Valid: true},
			pgtype.UUID{Bytes: parsedCursorID, Valid: true},
			limit,
		)
	}

	if err != nil {
		if errors.Is(err, message.ErrNotConversationMember) {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		http.Error(w, "failed to list messages", http.StatusInternalServerError)
		return
	}

	response := messageHistoryResponse{
		Messages: make([]messageResponse, 0, len(messages)),
	}

	for _, currentMessage := range messages {
		response.Messages = append(response.Messages, messageResponse{
			ID:              currentMessage.ID.String(),
			ConversationID:  currentMessage.ConversationID.String(),
			SenderID:        currentMessage.SenderID.String(),
			ClientMessageID: currentMessage.ClientMessageID.String(),
			Type:            currentMessage.Type,
			Content:         currentMessage.Content.String,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
