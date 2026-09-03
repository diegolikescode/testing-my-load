package main

import (
	"os"

	"github.com/phuslu/log"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/diegolikescode/testing-my-load/internal/server"
	"github.com/diegolikescode/testing-my-load/pkg"
)

func main() {
	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok && lvl != "" {
		log.DefaultLogger.SetLevel(log.ParseLevel(lvl))
	}

	go pkg.StartPprofServer()

	log.Info().Msg("start application")
	repo := repository.NewRepository()
	server := server.NewServer()

	log.Fatal().Err(server.Start(repo))
}
