package psql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

// Driver ...
type Driver = string

const (
	// PGDriver ...
	PGDriver Driver = "postgres"
)

// MigrateDB ...
func MigrateDB(database *sql.DB, migrationsPath string, driver Driver) error {
	db, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %+v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		driver,
		db,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %+v", err)
	}

	err = m.Up()
	if err != nil {
		if err == migrate.ErrNoChange {
			log.Printf("Migration did not change DB")
			return nil
		}
		return fmt.Errorf("failed to migrate: %+v", err)
	}

	return nil
}
