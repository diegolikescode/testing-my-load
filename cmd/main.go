package main

import (
	"github.com/phuslu/log"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/diegolikescode/testing-my-load/internal/server"
	"github.com/diegolikescode/testing-my-load/pkg"
)

func main() {
	go pkg.StartPprofServer()

	log.Info().Msg("start application")
	repo := repository.NewRepository()
	server := server.NewServer()

	log.Fatal().Err(server.Start(repo))
}
