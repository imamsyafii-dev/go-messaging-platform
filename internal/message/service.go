package message

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

var (
	ErrInvalidMessageType    = errors.New("invalid message type")
	ErrEmptyMessage          = errors.New("message content is required")
	ErrNotConversationMember = errors.New("user is not a conversation member")
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

	ListBefore(
		ctx context.Context,
		conversationID pgtype.UUID,
		cursorCreatedAt pgtype.Timestamptz,
		cursorID pgtype.UUID,
		limit int32,
	) ([]sqlc.Message, error)

	List(
		ctx context.Context,
		conversationID pgtype.UUID,
		limit int32,
	) ([]sqlc.Message, error)
}

type MembershipChecker interface {
	IsMember(
		ctx context.Context,
		conversationID pgtype.UUID,
		userID pgtype.UUID,
	) (bool, error)
}

type Service struct {
	repository MessageRepository
	membership MembershipChecker
}

func NewService(
	repository MessageRepository,
	membership MembershipChecker,
) *Service {
	return &Service{
		repository: repository,
		membership: membership,
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
	isMember, err := s.membership.IsMember(
		ctx,
		conversationID,
		senderID,
	)
	if err != nil {
		return sqlc.Message{}, fmt.Errorf("check conversation membership: %w", err)
	}

	if !isMember {
		return sqlc.Message{}, ErrNotConversationMember
	}

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

func (s *Service) ListBefore(
	ctx context.Context,
	conversationID pgtype.UUID,
	userID pgtype.UUID,
	cursorCreatedAt pgtype.Timestamptz,
	cursorID pgtype.UUID,
	limit int32,
) ([]sqlc.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	isMember, err := s.membership.IsMember(
		ctx,
		conversationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("check conversation membership: %w", err)
	}

	if !isMember {
		return nil, ErrNotConversationMember
	}

	return s.repository.ListBefore(
		ctx,
		conversationID,
		cursorCreatedAt,
		cursorID,
		limit,
	)
}

func (s *Service) List(
	ctx context.Context,
	conversationID pgtype.UUID,
	userID pgtype.UUID,
	limit int32,
) ([]sqlc.Message, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	isMember, err := s.membership.IsMember(
		ctx,
		conversationID,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("check conversation membership: %w", err)
	}

	if !isMember {
		return nil, ErrNotConversationMember
	}

	return s.repository.List(ctx, conversationID, limit)
}
