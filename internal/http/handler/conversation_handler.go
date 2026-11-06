package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/conversation"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/http/middleware"
)

type ConversationService interface {
	Create(
		ctx context.Context,
		conversationType string,
		name string,
		createdBy pgtype.UUID,
	) (sqlc.Conversation, error)
}

type ConversationHandler struct {
	service ConversationService
}

func NewConversationHandler(service ConversationService) *ConversationHandler {
	return &ConversationHandler{
		service: service,
	}
}

type createConversationRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type conversationResponse struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name,omitempty"`
	CreatedBy string `json:"created_by"`
}

func (h *ConversationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request createConversationRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	userID, ok := middleware.UserID(r.Context())
	if !ok {
		http.Error(w, "user identity not found", http.StatusUnauthorized)
		return
	}

	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		http.Error(w, "invalid user identity", http.StatusUnauthorized)
		return
	}

	createdConversation, err := h.service.Create(
		r.Context(),
		request.Type,
		request.Name,
		pgtype.UUID{
			Bytes: parsedUserID,
			Valid: true,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, conversation.ErrInvalidConversationType),
			errors.Is(err, conversation.ErrInvalidConversationName):
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "failed to create conversation", http.StatusInternalServerError)
		}

		return
	}

	response := conversationResponse{
		ID:        createdConversation.ID.String(),
		Type:      createdConversation.Type,
		Name:      createdConversation.Name.String,
		CreatedBy: createdConversation.CreatedBy.String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
