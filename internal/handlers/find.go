package handlers

import (
	"net/http"
	"strconv"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
	"github.com/phuslu/log"
)

type FindPersonHandler struct {
	validate   *validator.Validate
	repository repository.PersonRepository
}

func (h *FindPersonHandler) HandleFindPersonByID(c *echo.Context) error {
	id := c.Param("id")
	if id == "" {
		log.Warn().Msg("[FIND_BY_ID] the id param is empty")
		return nil
	}

	person := h.repository.FindByID(id)
	if person == nil {
		c.NoContent(http.StatusNotFound)
		return nil
	}
	c.JSON(http.StatusOK, person)

	// log stuff
	return nil
}

func (h *FindPersonHandler) HandleFindPersonByTerm(c *echo.Context) error {
	t := c.QueryParam("t")
	if t == "" {
		log.Warn().Msg("[FIND_BY_TERM] the term QueryParam is empty")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	people := h.repository.FindByTerm(t)

	c.JSON(http.StatusOK, people)
	return nil
}

func (h *FindPersonHandler) HandleCountPeople(c *echo.Context) error {
	count := h.repository.Count()
	if count == 0 {
		log.Error().Msg("something went wrong while trying to count people")
	}

	c.String(http.StatusOK, strconv.Itoa(count))
	return nil
}

func NewFindPersonHandler(repo repository.PersonRepository) *FindPersonHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &FindPersonHandler{validate: validate, repository: repo}
}
