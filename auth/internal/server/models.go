package server

import (
	"github.com/vladlim/auth-service-practice/auth/internal/providers/auth"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/tokens"
)

// Register Info...
type RegisterUserData struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func ServerRegisterReq2Provider(req RegisterUserData) auth.RegisterUserData {
	return auth.RegisterUserData{
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}
}

// Tokens...
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ServerTokens2Provider(req Tokens) tokens.Tokens {
	return tokens.Tokens{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}
}

func ProviderTokens2Server(req tokens.Tokens) Tokens {
	return Tokens{
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	}
}

// Login info...
type LoginUserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
