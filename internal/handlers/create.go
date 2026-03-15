package handlers

import (
	"net/http"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v5"
)

type CreatePersonHandler struct {
	validate   *validator.Validate
	repository *repository.PersonRepository
}

func (h *CreatePersonHandler) validateDto(dto repository.Person) error {
	err := h.validate.Struct(dto)
	return err
}

func (h *CreatePersonHandler) HandleCreatePerson(c *echo.Context) error {
	var dto repository.Person
	err := c.Bind(&dto)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}
	err = h.validateDto(dto)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return err
	}

	return nil
}

func NewCreatePersonHandler() *CreatePersonHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())

	return &CreatePersonHandler{validate: validate}
}
