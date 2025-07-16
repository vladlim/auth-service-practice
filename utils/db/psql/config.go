package psql

import "fmt"

const dbConnectionStringTemplate = "%s://%s:%s@%s:%d/%s?sslmode=%s"

// DB ...
type DB struct {
	Schema   string `yaml:"schema"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	IP       string `yaml:"ip"`
	Port     uint16 `yaml:"port"`
	DataBase string `yaml:"database"`
	SSL      string `yaml:"ssl"`

	MigrationsPath string `yaml:"migrations_path"`
}

// GetDBURL ...
func (db DB) GetDBURL() string {
	return fmt.Sprintf(
		dbConnectionStringTemplate,
		db.Schema,
		db.User,
		db.Password,
		db.IP,
		db.Port,
		db.DataBase,
		db.SSL,
	)
}
