package message

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *Repository) Create(
	ctx context.Context,
	conversationID pgtype.UUID,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
	messageType string,
	content pgtype.Text,
) (sqlc.Message, error) {
	return r.queries.CreateMessage(ctx, sqlc.CreateMessageParams{
		ConversationID:  conversationID,
		SenderID:        senderID,
		ClientMessageID: clientMessageID,
		Type:            messageType,
		Content:         content,
	})
}

func (r *Repository) GetByID(
	ctx context.Context,
	id pgtype.UUID,
) (sqlc.Message, error) {
	return r.queries.GetMessageByID(ctx, id)
}

func (r *Repository) GetByClientID(
	ctx context.Context,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
) (sqlc.Message, error) {
	return r.queries.GetMessageByClientID(ctx, sqlc.GetMessageByClientIDParams{
		SenderID:        senderID,
		ClientMessageID: clientMessageID,
	})
}

func (r *Repository) List(
	ctx context.Context,
	conversationID pgtype.UUID,
	limit int32,
) ([]sqlc.Message, error) {
	return r.queries.ListMessages(ctx, sqlc.ListMessagesParams{
		ConversationID: conversationID,
		Limit:          limit,
	})
}

func (r *Repository) CreateWithOutbox(
	ctx context.Context,
	conversationID pgtype.UUID,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
	messageType string,
	content pgtype.Text,
) (sqlc.Message, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return sqlc.Message{}, err
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	message, err := queries.CreateMessage(ctx, sqlc.CreateMessageParams{
		ConversationID:  conversationID,
		SenderID:        senderID,
		ClientMessageID: clientMessageID,
		Type:            messageType,
		Content:         content,
	})
	if err != nil {
		return sqlc.Message{}, err
	}

	payload, err := json.Marshal(map[string]any{
		"id":                message.ID,
		"conversation_id":   message.ConversationID,
		"sender_id":         message.SenderID,
		"client_message_id": message.ClientMessageID,
		"type":              message.Type,
		"content":           message.Content,
		"created_at":        message.CreatedAt,
	})
	if err != nil {
		return sqlc.Message{}, err
	}

	_, err = queries.CreateOutboxEvent(ctx, sqlc.CreateOutboxEventParams{
		AggregateType: "message",
		AggregateID:   message.ID,
		EventType:     "message.created",
		Payload:       payload,
	})
	if err != nil {
		return sqlc.Message{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Message{}, err
	}

	return message, nil
}

func (r *Repository) ListBefore(
	ctx context.Context,
	conversationID pgtype.UUID,
	cursorCreatedAt pgtype.Timestamptz,
	cursorID pgtype.UUID,
	limit int32,
) ([]sqlc.Message, error) {
	return r.queries.ListMessagesBefore(ctx, sqlc.ListMessagesBeforeParams{
		ConversationID: conversationID,
		CreatedAt:      cursorCreatedAt,
		ID:             cursorID,
		Limit:          limit,
	})
}
