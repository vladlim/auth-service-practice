package storage

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/models"
	"github.com/vladlim/auth-service-practice/utils/db/psql"
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
