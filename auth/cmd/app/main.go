package main

import (
	"log"
	"os"

	"github.com/vladlim/auth-service-practice/auth/internal/config"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/auth"
	"github.com/vladlim/auth-service-practice/auth/internal/providers/tokens"
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

	if err := tokens.InitJWT(conf.AccessSecret, conf.RefreshSecret); err != nil {
		log.Default().Printf("[ERR] Init jwt parse error: %s\n", err.Error())
		panic(err)
	}

	storage, err := storage.New(conf.DB.GetDBURL(), conf.DB.MigrationsPath)
	if err != nil {
		panic(err)
	}

	facade := facade.New(storage)

	authProvider := auth.New(facade)
	tokensProvider := tokens.New(facade)

	s := server.New(conf, authProvider, tokensProvider)
	panic(s.Start())
}
