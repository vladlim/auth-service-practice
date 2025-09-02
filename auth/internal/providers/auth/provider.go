package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, user models.RegisterUserData) (string, error)
	FindUserByUsername(ctx context.Context, username string) (string, string, error)
	FindUserByEmail(ctx context.Context, email string) (string, string, error)
	FindUserByID(ctx context.Context, userID string) (bool, error)

	CreateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error
	CreateTeacher(ctx context.Context, userID, universityID, degree string) error
	AddUserRole(ctx context.Context, userID, role string) error
	CheckUserRole(ctx context.Context, userID, role string) (bool, error)

	GetUserByID(ctx context.Context, userID string) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetStudentByID(ctx context.Context, userID string) (models.Student, error)
	GetStudentsByGroup(ctx context.Context, groupID string) ([]models.Student, error)
	GetTeacherByID(ctx context.Context, userID string) (models.Teacher, error)
	GetTeachersByUni(ctx context.Context, uniID string) ([]models.Teacher, error)
	GetUserRoles(ctx context.Context, userID string) ([]string, error)
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

// Activate keys...

func (p AuthProvider) ActivateStudent(ctx context.Context, userID string, claims jwt.MapClaims) error {
	groupID, ok1 := claims["group_id"].(string)
	universityID, ok2 := claims["university_id"].(string)
	enrollmentYear, ok3 := claims["enrollment_year"].(float64)

	if !ok1 || !ok2 || !ok3 {
		return errors.New("invalid student key parameters")
	}

	if activated, err := p.repository.CheckUserRole(ctx, userID, "student"); err != nil {
		return fmt.Errorf("failed to check role activation: %w", err)
	} else if activated {
		return errors.New("student role already activated")
	}

	if err := p.repository.CreateStudent(ctx, userID, groupID, universityID, int(enrollmentYear)); err != nil {
		return fmt.Errorf("failed to create student: %w", err)
	}

	if err := p.repository.AddUserRole(ctx, userID, "student"); err != nil {
		return fmt.Errorf("failed to add student role: %w", err)
	}

	return nil
}

func (p AuthProvider) ActivateTeacher(ctx context.Context, userID string, claims jwt.MapClaims) error {
	universityID, ok1 := claims["university_id"].(string)
	degree, ok2 := claims["degree"].(string)

	if !ok1 || !ok2 {
		return errors.New("invalid teacher key parameters")
	}

	if activated, err := p.repository.CheckUserRole(ctx, userID, "teacher"); err != nil {
		return fmt.Errorf("failed to check role activation: %w", err)
	} else if activated {
		return errors.New("teacher role already activated")
	}

	if err := p.repository.CreateTeacher(ctx, userID, universityID, degree); err != nil {
		return fmt.Errorf("failed to create teacher: %w", err)
	}

	if err := p.repository.AddUserRole(ctx, userID, "teacher"); err != nil {
		return fmt.Errorf("failed to add teacher role: %w", err)
	}

	return nil
}

// User Info...
func (p AuthProvider) GetUserByID(ctx context.Context, userID string) (User, error) {
	exists, err := p.repository.FindUserByID(ctx, userID)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	if !exists {
		return User{}, ErrUserNotFound
	}

	user, err := p.repository.GetUserByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	return DBUser2Provider(user), err
}

func (p AuthProvider) GetUserByEmail(ctx context.Context, email string) (User, error) {
	userID, _, err := p.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	exists, err := p.repository.FindUserByID(ctx, userID)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	if !exists {
		return User{}, ErrUserNotFound
	}

	user, err := p.repository.GetUserByEmail(ctx, email)
	if err != nil {
		return User{}, err
	}

	return DBUser2Provider(user), err
}

func (p AuthProvider) GetStudentByID(ctx context.Context, userID string) (Student, error) {
	exists, err := p.repository.FindUserByID(ctx, userID)
	if err != nil {
		return Student{}, ErrUserNotFound
	}
	if !exists {
		return Student{}, ErrUserNotFound
	}

	student, err := p.repository.GetStudentByID(ctx, userID)
	if err != nil {
		return Student{}, err
	}

	return DBStudent2Provider(student), err
}

func (p AuthProvider) GetStudents(ctx context.Context, groupIDs string) ([]models.Student, error) {
	return p.repository.GetStudentsByGroup(ctx, groupIDs)
}

func (p AuthProvider) GetTeacherByID(ctx context.Context, userID string) (Teacher, error) {
	exists, err := p.repository.FindUserByID(ctx, userID)
	if err != nil {
		return Teacher{}, ErrUserNotFound
	}
	if !exists {
		return Teacher{}, ErrUserNotFound
	}

	teacher, err := p.repository.GetTeacherByID(ctx, userID)
	if err != nil {
		return Teacher{}, err
	}

	return DBTeacher2Provider(teacher), err
}

func (p AuthProvider) GetTeachers(ctx context.Context, uniIDs string) ([]models.Teacher, error) {
	return p.repository.GetTeachersByUni(ctx, uniIDs)
}

func (p AuthProvider) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return p.repository.GetUserRoles(ctx, userID)
}
