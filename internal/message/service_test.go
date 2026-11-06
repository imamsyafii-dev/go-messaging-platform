package message_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/imamsyafii-dev/go-messaging-platform/internal/message"
)

type mockMessageRepository struct {
	called bool
}
type mockMembershipChecker struct {
	isMember bool
}

func (m *mockMembershipChecker) IsMember(
	ctx context.Context,
	conversationID pgtype.UUID,
	userID pgtype.UUID,
) (bool, error) {
	return m.isMember, nil
}

func (m *mockMessageRepository) Create(
	ctx context.Context,
	conversationID pgtype.UUID,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
	messageType string,
	content pgtype.Text,
) (sqlc.Message, error) {
	return sqlc.Message{}, nil
}

func (m *mockMessageRepository) CreateWithOutbox(
	ctx context.Context,
	conversationID pgtype.UUID,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
	messageType string,
	content pgtype.Text,
) (sqlc.Message, error) {
	m.called = true

	return sqlc.Message{
		Type:    messageType,
		Content: content,
	}, nil
}

func (m *mockMessageRepository) GetByClientID(
	ctx context.Context,
	senderID pgtype.UUID,
	clientMessageID pgtype.UUID,
) (sqlc.Message, error) {
	return sqlc.Message{}, nil
}

func TestCreateRejectsInvalidType(t *testing.T) {
	repository := &mockMessageRepository{}
	service := message.NewService(
		repository,
		&mockMembershipChecker{isMember: true},
	)

	_, err := service.Create(
		context.Background(),
		pgtype.UUID{},
		pgtype.UUID{},
		pgtype.UUID{},
		"invalid",
		"hello",
	)

	if err != message.ErrInvalidMessageType {
		t.Fatalf("expected invalid message type error, got %v", err)
	}

	if repository.called {
		t.Fatal("repository should not be called")
	}
}

func TestCreateTextMessage(t *testing.T) {
	repository := &mockMessageRepository{}

	service := message.NewService(
		repository,
		&mockMembershipChecker{isMember: true},
	)

	_, err := service.Create(
		context.Background(),
		pgtype.UUID{},
		pgtype.UUID{},
		pgtype.UUID{},
		"text",
		"hello",
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !repository.called {
		t.Fatal("repository should be called")
	}
}

func TestCreateRejectsNonMember(t *testing.T) {
	repository := &mockMessageRepository{}
	service := message.NewService(
		repository,
		&mockMembershipChecker{isMember: false},
	)

	_, err := service.Create(
		context.Background(),
		pgtype.UUID{},
		pgtype.UUID{},
		pgtype.UUID{},
		"text",
		"hello",
	)

	if err != message.ErrNotConversationMember {
		t.Fatalf("expected membership error, got %v", err)
	}

	if repository.called {
		t.Fatal("repository should not be called for non-member")
	}
}
