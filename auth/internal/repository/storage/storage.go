package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	storage "github.com/vladlim/auth-service-practice/auth/internal/repository/storage/queries"
	"github.com/vladlim/utils/db/psql"
)

type Storage struct {
	db *sql.DB
}

func New(dbURL string, migrationsPath string) (Storage, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return Storage{}, err
	}
	err = db.Ping()
	if err != nil {
		return Storage{}, err
	}
	psql.MigrateDB(db, migrationsPath, psql.PGDriver)
	return Storage{
		db: db,
	}, nil
}

func (s Storage) CreateUser(ctx context.Context, user models.RegisterUserData) (string, error) {
	var userID string

	err := s.db.QueryRowContext(ctx, storage.CreateUserQuery,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.FirstName,
		user.LastName,
	).Scan(
		&userID,
	)

	return userID, err
}

func (s Storage) FindUserByUsername(ctx context.Context, username string) (string, string, error) {
	var userID, userPassword string
	err := s.db.QueryRowContext(ctx, storage.FindUserByUsernameQuery, username).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s Storage) FindUserByEmail(ctx context.Context, email string) (string, string, error) {
	var userID, userPassword string
	err := s.db.QueryRowContext(ctx, storage.FindUserByEmailQuery, email).Scan(&userID, &userPassword)
	return userID, userPassword, err
}

func (s Storage) GetPeople(ctx context.Context) ([]models.Person, error) {
	q := `
		SELECT
			id,
			username
		FROM people
	`
	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	people := make([]models.Person, 0)
	err = sqlx.StructScan(rows, &people)
	if err != nil {
		return nil, err
	}
	return people, nil
}
