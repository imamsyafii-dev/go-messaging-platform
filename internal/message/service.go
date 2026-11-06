package message

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

var (
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrEmptyMessage       = errors.New("message content is required")
)

type MessageRepository interface {
	Create(
		ctx context.Context,
		conversationID pgtype.UUID,
		senderID pgtype.UUID,
		clientMessageID pgtype.UUID,
		messageType string,
		content pgtype.Text,
	) (sqlc.Message, error)

	GetByClientID(
		ctx context.Context,
		senderID pgtype.UUID,
		clientMessageID pgtype.UUID,
	) (sqlc.Message, error)

	CreateWithOutbox(
		ctx context.Context,
		conversationID pgtype.UUID,
		senderID pgtype.UUID,
		clientMessageID pgtype.UUID,
		messageType string,
		content pgtype.Text,
	) (sqlc.Message, error)
}

type Service struct {
	repository MessageRepository
}

func NewService(repository MessageRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	conversationID pgtype.UUID,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
	messageType string,
	content string,
) (sqlc.Message, error) {
	if messageType != "text" &&
		messageType != "image" &&
		messageType != "file" {
		return sqlc.Message{}, ErrInvalidMessageType
	}

	if messageType == "text" && content == "" {
		return sqlc.Message{}, ErrEmptyMessage
	}

	messageContent := pgtype.Text{
		String: content,
		Valid:  content != "",
	}

	message, err := s.repository.CreateWithOutbox(
		ctx,
		conversationID,
		senderID,
		clientMessageID,
		messageType,
		messageContent,
	)
	if err != nil {
		return sqlc.Message{}, fmt.Errorf("create message: %w", err)
	}

	return message, nil
}
