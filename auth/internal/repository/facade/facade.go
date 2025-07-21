package facade

import (
	"context"

	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/storage"
)

type Storage interface {
	storage.Storage
}

type Facade struct {
	storage Storage
}

func New(storage Storage) Facade {
	return Facade{
		storage: storage,
	}
}

// Auth...

func (f Facade) CreateUser(ctx context.Context, user models.RegisterUserData) (string, error) {
	return f.storage.CreateUser(ctx, user)
}

func (f Facade) FindUserByUsername(ctx context.Context, username string) (string, string, error) {
	return f.storage.FindUserByUsername(ctx, username)
}

func (f Facade) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	return f.storage.FindUserByEmail(ctx, email)
}

func (f Facade) FindUserByID(ctx context.Context, userID string) (bool, error) {
	return f.storage.FindUserByID(ctx, userID)
}

// User Info...

func (f Facade) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	return f.storage.GetUserByID(ctx, userID)
}

func (f Facade) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return f.storage.GetUserByEmail(ctx, email)
}

func (f Facade) GetStudentByID(ctx context.Context, userID string) (models.Student, error) {
	return f.storage.GetStudentByID(ctx, userID)
}

func (f Facade) GetStudentsByGroup(ctx context.Context, groupID string) ([]models.Student, error) {
	return f.storage.GetStudentsByGroup(ctx, groupID)
}

func (f Facade) GetTeacherByID(ctx context.Context, userID string) (models.Teacher, error) {
	return f.storage.GetTeacherByID(ctx, userID)
}

func (f Facade) GetTeachersByUni(ctx context.Context, uniID string) ([]models.Teacher, error) {
	return f.storage.GetTeachersByUni(ctx, uniID)
}

func (f Facade) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	return f.storage.GetUserRoles(ctx, userID)
}

// Keys...

func (f Facade) CreateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error {
	return f.storage.CreateStudent(ctx, userID, groupID, universityID, enrollmentYear)
}

func (f Facade) CreateTeacher(ctx context.Context, userID, universityID, degree string) error {
	return f.storage.CreateTeacher(ctx, userID, universityID, degree)
}

func (f Facade) AddUserRole(ctx context.Context, userID, role string) error {
	return f.storage.AddUserRole(ctx, userID, role)
}

func (f Facade) CheckUserRole(ctx context.Context, userID, role string) (bool, error) {
	return f.storage.CheckUserRole(ctx, userID, role)
}

// Transactions (activate keys)...

func (f *Facade) ActivateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error {
	tx, err := f.storage.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.CreateStudent(ctx, userID, groupID, universityID, enrollmentYear); err != nil {
		return err
	}

	if err := tx.AddUserRole(ctx, userID, "student"); err != nil {
		return err
	}

	return tx.Commit()
}

func (f *Facade) ActivateTeacher(ctx context.Context, userID, universityID, degree string) error {
	tx, err := f.storage.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.CreateTeacher(ctx, userID, universityID, degree); err != nil {
		return err
	}

	if err := tx.AddUserRole(ctx, userID, "teacher"); err != nil {
		return err
	}

	return tx.Commit()
}
