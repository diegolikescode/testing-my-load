package server

import (
	"net/http"
	"strconv"

	"github.com/diegolikescode/testing-my-load/internal/handlers"
	"github.com/labstack/echo/v5"
)

type Server struct {
	echoServer *echo.Echo
	Host       string
	Port       int
}

func (s *Server) setupRoutes() {
	s.echoServer.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Message string }{Message: "all good"})
	})

	createHandler := handlers.NewCreatePersonHandler()

	s.echoServer.POST("/pessoas", createHandler.HandleCreatePerson)
}

func (s *Server) Start() error {
	s.setupRoutes()
	return s.echoServer.Start(s.Host + ":" + strconv.Itoa(s.Port))
	// return s.echoServer.Start("localhost:6969")
}

func NewServer(host string, port int) *Server {
	e := echo.New()

	return &Server{
		Host:       host,
		Port:       port,
		echoServer: e,
	}
}
