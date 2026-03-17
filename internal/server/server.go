package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/diegolikescode/testing-my-load/internal/handlers"
	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/diegolikescode/testing-my-load/pkg"
	"github.com/labstack/echo/v5"
	"github.com/phuslu/log"
)

type Server struct {
	echoServer *echo.Echo
	Host       string
	Port       int
}

/*
TODO: write logs in a file, map the file in docker-compose to a external one so I can have insights into whats happening
*/

func (s *Server) setupRoutes(repo repository.PersonRepository) {
	s.echoServer.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Message string }{Message: "all good"})
	})

	createHandler := handlers.NewCreatePersonHandler(repo)
	findHandler := handlers.NewFindPersonHandler(repo)

	s.echoServer.POST("/pessoas", createHandler.HandleCreatePerson)
	s.echoServer.GET("pessoas/:id", findHandler.HandleFindPersonById)
}

func (s *Server) Start(repo repository.PersonRepository) error {
	s.setupRoutes(repo)

	log.Info().Msg(fmt.Sprintf("Starting server on host %s and port %d", s.Host, s.Port))
	return s.echoServer.Start(s.Host + ":" + strconv.Itoa(s.Port))
}

func NewServer() *Server {
	e := echo.New()

	portStr := pkg.GetEnvOrDieTrying("SERVER_PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(fmt.Sprintf("the port variable was not valid, expected number got=%s", portStr))
	}

	return &Server{
		Host:       pkg.GetEnvOrDieTrying("SERVER_HOST"),
		Port:       port,
		echoServer: e,
	}
}
