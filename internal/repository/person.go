package repository

type Person struct {
	ID         string   `json:"id"`
	Nome       string   `json:"nome" validate:"required"`
	Apelido    string   `json:"apelido" validate:"required"`
	Nascimento string   `json:"nascimento" validate:"required,datetime=2006/01/02"`
	Stack      []string `json:"stack" validate:"omitzero,dive,max=32"`
}

type PersonRepository interface {
	Create(uuid, name, apelido, nascimento string, stack []string) error
	FindByID(id string) *Person
	FindByTerm(t string) []*Person
	Count() int
}
