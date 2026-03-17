package handlers

import (
	"net/http"

	"github.com/diegolikescode/testing-my-load/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type CreatePersonHandler struct {
	validate   *validator.Validate
	repository repository.PersonRepository
}

func (h *CreatePersonHandler) validateDto(dto repository.Person) error {
	err := h.validate.Struct(dto)
	return err
}

func (h *CreatePersonHandler) HandleCreatePerson(c *echo.Context) error {
	var dto repository.Person
	err := c.Bind(&dto)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		// log stuff
		return err
	}
	err = h.validateDto(dto)
	if err != nil {
		c.String(http.StatusUnprocessableEntity, err.Error())
		// log stuff
		return err
	}

	exists := h.repository.CheckIfExists(dto.Apelido)
	if exists {
		c.Response().WriteHeader(http.StatusUnprocessableEntity)
		// log stuff
		return nil
	}

	uid := uuid.NewString()
	h.repository.Create(uid, dto.Apelido, dto.Nome, dto.Nascimento, dto.Stack)
	c.Response().Header().Set("Location", uid)

	// log stuff
	return nil
}

func NewCreatePersonHandler(repo repository.PersonRepository) *CreatePersonHandler {
	validate := validator.New(validator.WithRequiredStructEnabled())
	return &CreatePersonHandler{validate: validate, repository: repo}
}
