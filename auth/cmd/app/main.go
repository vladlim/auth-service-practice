package main

import (
	"os"

	"github.com/vladlim/auth-service-practice/auth/internal/config"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/people"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/facade"
	"github.com/vladlim/auth-service-practice/auth/internal/repository/storage"
	"github.com/vladlim/auth-service-practice/auth/internal/server"
)

func main() {
	configPath := os.Args[1]
	conf, err := config.Parse(configPath)
	if err != nil {
		panic(err)
	}

	storage, err := storage.New(conf.DB.GetDBURL(), conf.DB.MigrationsPath)
	if err != nil {
		panic(err)
	}

	facade := facade.New(storage)

	provider := people.New(facade)

	s := server.New(conf, provider)
	panic(s.Start())
}
