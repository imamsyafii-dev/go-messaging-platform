package outbox

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: sqlc.New(pool),
	}
}

func (r *Repository) Create(
	ctx context.Context,
	aggregateType string,
	aggregateID pgtype.UUID,
	eventType string,
	payload []byte,
) (sqlc.OutboxEvent, error) {
	return r.queries.CreateOutboxEvent(ctx, sqlc.CreateOutboxEventParams{
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		EventType:     eventType,
		Payload:       payload,
	})
}

func (r *Repository) ListPending(
	ctx context.Context,
	limit int32,
) ([]sqlc.OutboxEvent, error) {
	return r.queries.ListPendingOutboxEvents(ctx, limit)
}

func (r *Repository) MarkPublished(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.MarkOutboxEventPublished(ctx, id)
}

func (r *Repository) IncrementAttempts(
	ctx context.Context,
	id pgtype.UUID,
) error {
	return r.queries.IncrementOutboxEventAttempts(ctx, id)
}
