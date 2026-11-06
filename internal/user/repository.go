package user

import (
	"context"

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
	email string,
	passwordHash string,
	username string,
	displayName string,
) (sqlc.User, error) {
	return r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
		Username:     username,
		DisplayName:  displayName,
	})
}

func (r *Repository) GetByEmail(
	ctx context.Context,
	email string,
) (sqlc.User, error) {
	return r.queries.GetUserByEmail(ctx, email)
}
