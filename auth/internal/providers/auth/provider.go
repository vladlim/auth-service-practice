package auth

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.RegisterUserData) (string, error)
	FindUserByUsername(ctx context.Context, username string) (string, string, error)
	FindUserByEmail(ctx context.Context, email string) (string, string, error)
}

type AuthProvider struct {
	repository Repository
}

func New(repository Repository) AuthProvider {
	return AuthProvider{
		repository: repository,
	}
}

func (p AuthProvider) RegisterUser(ctx context.Context, user RegisterUserData) (string, error) {
	userConv := ProviderRegisterReq2DB(user)
	if _, _, err := p.repository.FindUserByUsername(ctx, userConv.Username); err == nil {
		return "", ErrUsernameExists
	} else if err != sql.ErrNoRows {
		return "", err
	}

	if _, _, err := p.repository.FindUserByEmail(ctx, userConv.Email); err == nil {
		return "", ErrEmailExists
	} else if err != sql.ErrNoRows {
		return "", err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashingPassword
	}

	user.Password = string(hashedPassword)

	userID, err := p.repository.CreateUser(ctx, ProviderRegisterReq2DB(user))

	return userID, err
}

func (p AuthProvider) LoginUser(ctx context.Context, login, password string) (string, error) {
	var userID, userPassword string
	var err error
	if strings.Contains(login, "@") {
		userID, userPassword, err = p.repository.FindUserByEmail(ctx, login)
	} else {
		userID, userPassword, err = p.repository.FindUserByUsername(ctx, login)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password)); err != nil {
		return "", ErrIncorrectPassword
	}

	return userID, nil
}
