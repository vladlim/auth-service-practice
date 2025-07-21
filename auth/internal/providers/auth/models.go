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

// User...

type User struct {
	ID        string
	Username  string
	Email     string
	FirstName string
	LastName  string
	CreatedAt string
}

func DBUser2Provider(user models.User) User {
	return User{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
	}
}

// Student...

type Student struct {
	ID             string
	GroupID        string
	UniversityID   string
	EnrollmentYear int
	User           User
}

func DBStudent2Provider(student models.Student) Student {
	return Student{
		ID:             student.UserID,
		GroupID:        student.GroupID,
		UniversityID:   student.UniversityID,
		EnrollmentYear: student.EnrollmentYear,
		User:           DBUser2Provider(student.User),
	}
}

// Teacher

type Teacher struct {
	ID           string
	UniversityID string
	Degree       string
	User         User
}

func DBTeacher2Provider(teacher models.Teacher) Teacher {
	return Teacher{
		ID:           teacher.UserID,
		UniversityID: teacher.UniversityID,
		Degree:       teacher.Degree,
		User:         DBUser2Provider(teacher.User),
	}
}
