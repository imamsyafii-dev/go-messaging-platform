package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/redis/go-redis/v9"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

const EventsChannel = "messaging.events"

type EventRepository interface {
	ListPending(ctx context.Context, limit int32) ([]sqlc.OutboxEvent, error)
	MarkPublished(ctx context.Context, id pgtype.UUID) error
	IncrementAttempts(ctx context.Context, id pgtype.UUID) error
}

type Worker struct {
	repository EventRepository
	redis      *redis.Client
	batchSize  int32
}

func NewWorker(
	repository EventRepository,
	redisClient *redis.Client,
	batchSize int32,
) *Worker {
	if batchSize <= 0 {
		batchSize = 100
	}

	return &Worker{
		repository: repository,
		redis:      redisClient,
		batchSize:  batchSize,
	}
}

func (w *Worker) RunOnce(ctx context.Context) error {
	events, err := w.repository.ListPending(ctx, w.batchSize)
	if err != nil {
		return fmt.Errorf("list pending outbox events: %w", err)
	}

	for _, event := range events {
		if err := w.redis.Publish(ctx, EventsChannel, event.Payload).Err(); err != nil {
			if incrementErr := w.repository.IncrementAttempts(ctx, event.ID); incrementErr != nil {
				log.Printf("increment outbox attempts: %v", incrementErr)
			}

			return fmt.Errorf("publish outbox event %s: %w", event.ID, err)
		}

		if err := w.repository.MarkPublished(ctx, event.ID); err != nil {
			return fmt.Errorf("mark outbox event %s published: %w", event.ID, err)
		}
	}

	return nil
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		if err := w.RunOnce(ctx); err != nil {
			log.Printf("outbox worker: %v", err)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
