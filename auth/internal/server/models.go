package server

import "github.com/vladlim/auth-service-practice/auth/internal/providers/auth"

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

type LoginUserData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
