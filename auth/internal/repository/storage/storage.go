package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	storage "github.com/vladlim/auth-service-practice/auth/internal/repository/storage/queries"
	"github.com/vladlim/utils/db/psql"
)

type Storage interface {
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

	BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error)
}

// Tx interface for transactions
type Tx interface {
	CreateUser(ctx context.Context, user models.RegisterUserData) (string, error)
	FindUserByUsername(ctx context.Context, username string) (string, string, error)
	FindUserByEmail(ctx context.Context, email string) (string, string, error)
	CreateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error
	CreateTeacher(ctx context.Context, userID, universityID, degree string) error
	AddUserRole(ctx context.Context, userID, role string) error
	CheckUserRole(ctx context.Context, userID, role string) (bool, error)

	Commit() error
	Rollback() error
}

// DBStorage Storage implementation
type DBStorage struct {
	db *sql.DB
}

func New(dbURL string, migrationsPath string) (*DBStorage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	psql.MigrateDB(db, migrationsPath, psql.PGDriver)
	return &DBStorage{db: db}, nil
}

// storageTx Tx implementation
type storageTx struct {
	tx *sql.Tx
}

func (s *DBStorage) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	tx, err := s.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &storageTx{tx: tx}, nil
}

// DBStorage without transactions
func (s *DBStorage) CreateUser(ctx context.Context, user models.RegisterUserData) (string, error) {
	var userID string
	err := s.db.QueryRowContext(ctx, storage.CreateUserQuery,
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName).Scan(&userID)
	return userID, err
}

func (s *DBStorage) FindUserByUsername(ctx context.Context, username string) (string, string, error) {
	var userID, userPassword string
	err := s.db.QueryRowContext(ctx, storage.FindUserByUsernameQuery, username).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s *DBStorage) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	var userID, userPassword string
	err := s.db.QueryRowContext(ctx, storage.FindUserByEmailQuery, email).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s *DBStorage) FindUserByID(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, storage.FindUserByIDQuery, userID).Scan(&exists)
	return exists, err
}

func (s *DBStorage) CreateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error {
	_, err := s.db.ExecContext(ctx, storage.CreateStudentQuery, userID, groupID, universityID, enrollmentYear)
	return err
}

func (s *DBStorage) CreateTeacher(ctx context.Context, userID, universityID, degree string) error {
	_, err := s.db.ExecContext(ctx, storage.CreateTeacherQuery, userID, universityID, degree)
	return err
}

func (s *DBStorage) AddUserRole(ctx context.Context, userID, role string) error {
	_, err := s.db.ExecContext(ctx, storage.AddUserRoleQuery, userID, role)
	return err
}

func (s *DBStorage) CheckUserRole(ctx context.Context, userID, role string) (bool, error) {
	var exists bool
	err := s.db.QueryRowContext(ctx, storage.CheckUserRoleQuery, userID, role).Scan(&exists)
	return exists, err
}

func (s *DBStorage) GetUserByID(ctx context.Context, userID string) (models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx, storage.GetUserByIDQuery, userID).Scan(&user.ID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}

func (s *DBStorage) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx, storage.GetUserByEmailQuery, email).Scan(&user.ID, &user.Username, &user.Email,
		&user.FirstName, &user.LastName, &user.CreatedAt)
	return user, err
}

func (s *DBStorage) GetStudentByID(ctx context.Context, studentID string) (models.Student, error) {
	var student models.Student
	var user models.User

	err := s.db.QueryRowContext(ctx, storage.GetStudentByIDQuery, studentID).Scan(
		&student.UserID, &student.GroupID, &student.UniversityID, &student.EnrollmentYear,
		&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err != nil {
		return models.Student{}, err
	}

	student.User = user
	return student, nil
}

func (s *DBStorage) GetStudentsByGroup(ctx context.Context, groupIDs string) ([]models.Student, error) {
	rows, err := s.db.QueryContext(ctx, storage.GetStudentsByGroupQuery, groupIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []models.Student
	for rows.Next() {
		var student models.Student
		var user models.User

		err := rows.Scan(
			&student.UserID, &student.GroupID, &student.UniversityID, &student.EnrollmentYear,
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		student.User = user
		students = append(students, student)
	}

	return students, nil
}

func (s *DBStorage) GetTeacherByID(ctx context.Context, teacherID string) (models.Teacher, error) {
	var teacher models.Teacher
	var user models.User

	err := s.db.QueryRowContext(ctx, storage.GetTeacherByIDQuery, teacherID).Scan(
		&teacher.UserID, &teacher.UniversityID, &teacher.Degree,
		&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err != nil {
		return models.Teacher{}, err
	}

	teacher.User = user
	return teacher, nil
}

func (s *DBStorage) GetTeachersByUni(ctx context.Context, uniIDs string) ([]models.Teacher, error) {
	rows, err := s.db.QueryContext(ctx, storage.GetTeachersByUniversityQuery, uniIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var teachers []models.Teacher
	for rows.Next() {
		var teacher models.Teacher
		var user models.User

		err := rows.Scan(
			&teacher.UserID, &teacher.UniversityID, &teacher.Degree,
			&user.ID, &user.Username, &user.Email, &user.FirstName, &user.LastName, &user.CreatedAt)
		if err != nil {
			return nil, err
		}

		teacher.User = user
		teachers = append(teachers, teacher)
	}

	return teachers, nil
}

func (s *DBStorage) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, storage.GetUserRolesQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []string
	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return roles, nil
}

// storageTx (transactions)
func (s *storageTx) CreateUser(ctx context.Context, user models.RegisterUserData) (string, error) {
	var userID string
	err := s.tx.QueryRowContext(ctx, storage.CreateUserQuery,
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName).Scan(&userID)
	return userID, err
}

func (s *storageTx) FindUserByUsername(ctx context.Context, username string) (string, string, error) {
	var userID, userPassword string
	err := s.tx.QueryRowContext(ctx, storage.FindUserByUsernameQuery, username).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s *storageTx) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	var userID, userPassword string
	err := s.tx.QueryRowContext(ctx, storage.FindUserByEmailQuery, email).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s *storageTx) CreateStudent(ctx context.Context, userID, groupID, universityID string, enrollmentYear int) error {
	_, err := s.tx.ExecContext(ctx, storage.CreateStudentQuery, userID, groupID, universityID, enrollmentYear)
	return err
}

func (s *storageTx) CreateTeacher(ctx context.Context, userID, universityID, degree string) error {
	_, err := s.tx.ExecContext(ctx, storage.CreateTeacherQuery, userID, universityID, degree)
	return err
}

func (s *storageTx) AddUserRole(ctx context.Context, userID, role string) error {
	_, err := s.tx.ExecContext(ctx, storage.AddUserRoleQuery, userID, role)
	return err
}

func (s *storageTx) CheckUserRole(ctx context.Context, userID, role string) (bool, error) {
	var exists bool
	err := s.tx.QueryRowContext(ctx, storage.CheckUserRoleQuery, userID, role).Scan(&exists)
	return exists, err
}

func (s *storageTx) Commit() error {
	return s.tx.Commit()
}

func (s *storageTx) Rollback() error {
	return s.tx.Rollback()
}
