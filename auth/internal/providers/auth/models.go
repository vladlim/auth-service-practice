package auth

import "github.com/vladlim/auth-service-practice/auth/internal/repository/models"

type RegisterUserData struct {
	Username  string
	Email     string
	Password  string
	FirstName string
	LastName  string
}

func ProviderRegisterReq2DB(user RegisterUserData) models.RegisterUserData {
	return models.RegisterUserData{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.Password,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
	}
}
