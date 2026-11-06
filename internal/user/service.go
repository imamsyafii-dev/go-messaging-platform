package user

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	Create(
		ctx context.Context,
		email string,
		passwordHash string,
		username string,
		displayName string,
	) (sqlc.User, error)
	GetByEmail(
		ctx context.Context,
		email string,
	) (sqlc.User, error)
}

type Service struct {
	repository UserRepository
}

func NewService(repository UserRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUser(
	ctx context.Context,
	email string,
	password string,
	username string,
	displayName string,
) (sqlc.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return sqlc.User{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.repository.Create(
		ctx,
		email,
		string(passwordHash),
		username,
		displayName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sqlc.User{}, ErrUserAlreadyExists
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return sqlc.User{}, ErrUserAlreadyExists
		}

		return sqlc.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *Service) Authenticate(
	ctx context.Context,
	email string,
	password string,
) (sqlc.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return sqlc.User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return sqlc.User{}, ErrInvalidCredentials
	}

	return user, nil
}
