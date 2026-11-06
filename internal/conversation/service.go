package conversation

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

var (
	ErrInvalidConversationType = errors.New("invalid conversation type")
	ErrInvalidConversationName = errors.New("group conversation name is required")
)

type ConversationRepository interface {
	CreateWithOwner(
		ctx context.Context,
		conversationType string,
		name pgtype.Text,
		createdBy pgtype.UUID,
	) (sqlc.Conversation, error)
}

type Service struct {
	repository ConversationRepository
}

func NewService(repository ConversationRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	conversationType string,
	name string,
	createdBy pgtype.UUID,
) (sqlc.Conversation, error) {
	if conversationType != "direct" && conversationType != "group" {
		return sqlc.Conversation{}, ErrInvalidConversationType
	}

	if conversationType == "group" && name == "" {
		return sqlc.Conversation{}, ErrInvalidConversationName
	}

	conversationName := pgtype.Text{
		String: name,
		Valid:  name != "",
	}
	conversation, err := s.repository.CreateWithOwner(
		ctx,
		conversationType,
		conversationName,
		createdBy,
	)
	if err != nil {
		return sqlc.Conversation{}, fmt.Errorf("create conversation with owner: %w", err)
	}

	return conversation, nil
}
