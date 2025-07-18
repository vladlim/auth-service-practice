package facade

import (
	"context"

	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
)

type Storage interface {
	CreateUser(ctx context.Context, user models.RegisterUserData) (string, error)
	FindUserByUsername(ctx context.Context, username string) (string, string, error)
	FindUserByEmail(ctx context.Context, email string) (string, string, error)
}

type Facade struct {
	storage Storage
}

func New(storage Storage) Facade {
	return Facade{
		storage: storage,
	}
}

func (f Facade) CreateUser(ctx context.Context, user models.RegisterUserData) (string, error) {
	return f.storage.CreateUser(ctx, user)
}

func (f Facade) FindUserByUsername(ctx context.Context, username string) (string, string, error) {
	return f.storage.FindUserByUsername(ctx, username)
}

func (f Facade) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	return f.storage.FindUserByEmail(ctx, email)
}
