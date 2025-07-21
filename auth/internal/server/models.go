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

// User

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
}

func ProviderUser2Server(user auth.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

// Student

type Student struct {
	ID             string `json:"id"`
	GroupID        string `json:"group_id"`
	UniversityID   string `json:"university_id"`
	EnrollmentYear int    `json:"enrollment_year"`
	User           User   `json:"user"`
}

func ProviderStudent2Server(student auth.Student) Student {
	return Student{
		ID:             student.ID,
		GroupID:        student.GroupID,
		UniversityID:   student.UniversityID,
		EnrollmentYear: student.EnrollmentYear,
		User:           ProviderUser2Server(student.User),
	}
}

// Teacher

type Teacher struct {
	ID           string `json:"id"`
	UniversityID string `json:"unversity_id"`
	Degree       string `json:"degree"`
	User         User   `json:"user"`
}

func ProviderTeacher2Server(teacher auth.Teacher) Teacher {
	return Teacher{
		ID:           teacher.ID,
		UniversityID: teacher.UniversityID,
		Degree:       teacher.Degree,
		User:         ProviderUser2Server(teacher.User),
	}
}
