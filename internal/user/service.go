package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/imamsyafii-dev/go-messaging-platform/internal/database/sqlc"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		email string,
		passwordHash string,
		username string,
		displayName string,
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
		return sqlc.User{}, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}
