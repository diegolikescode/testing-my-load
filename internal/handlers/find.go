package handlers

import (
	"net/http"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type FindPersonHandler struct {
	validate   *validator.Validate
	repository repository.PersonRepository
}

func (h *FindPersonHandler) HandleFindPersonById(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		// log stuff
		return nil
	}

	person := h.repository.FindByID(id)
	if person == nil {
		// not found or problem
	}
	c.JSON(http.StatusOK, person)

	// log stuff
	return nil
}

func NewFindPersonHandler(repo repository.PersonRepository) *FindPersonHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &FindPersonHandler{validate: validate, repository: repo}
}
