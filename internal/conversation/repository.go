package conversation

import (
	"context"

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
	conversationType string,
	name pgtype.Text,
	createdBy pgtype.UUID,
) (sqlc.Conversation, error) {
	return r.queries.CreateConversation(ctx, sqlc.CreateConversationParams{
		Type:      conversationType,
		Name:      name,
		CreatedBy: createdBy,
	})
}

func (r *Repository) AddMember(
	ctx context.Context,
	conversationID pgtype.UUID,
	userID pgtype.UUID,
	role string,
) (sqlc.ConversationMember, error) {
	return r.queries.AddConversationMember(ctx, sqlc.AddConversationMemberParams{
		ConversationID: conversationID,
		UserID:         userID,
		Role:           role,
	})
}

func (r *Repository) GetByID(
	ctx context.Context,
	id pgtype.UUID,
) (sqlc.Conversation, error) {
	return r.queries.GetConversationByID(ctx, id)
}

func (r *Repository) ListMembers(
	ctx context.Context,
	conversationID pgtype.UUID,
) ([]sqlc.ConversationMember, error) {
	return r.queries.ListConversationMembers(ctx, conversationID)
}

func (r *Repository) CreateWithOwner(
	ctx context.Context,
	conversationType string,
	name pgtype.Text,
	createdBy pgtype.UUID,
) (sqlc.Conversation, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return sqlc.Conversation{}, err
	}
	defer tx.Rollback(ctx)

	queries := r.queries.WithTx(tx)

	conversation, err := queries.CreateConversation(ctx, sqlc.CreateConversationParams{
		Type:      conversationType,
		Name:      name,
		CreatedBy: createdBy,
	})
	if err != nil {
		return sqlc.Conversation{}, err
	}

	_, err = queries.AddConversationMember(ctx, sqlc.AddConversationMemberParams{
		ConversationID: conversation.ID,
		UserID:         createdBy,
		Role:           "admin",
	})
	if err != nil {
		return sqlc.Conversation{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return sqlc.Conversation{}, err
	}

	return conversation, nil
}
