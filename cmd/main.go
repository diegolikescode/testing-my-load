package main

import (
	"github.com/phuslu/log"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/diegolikescode/testing-my-load/internal/server"
)

func main() {
	log.Info().Msg("Start application")
	repo := repository.NewRepository()
	server := server.NewServer()

	log.Fatal().Err(server.Start(repo))
}
