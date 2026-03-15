package main

import (
	"log"

	"github.com/diegolikescode/testing-my-load/internal/server"
)

func main() {
	server := server.NewServer()

	log.Fatal(server.Start())
}
