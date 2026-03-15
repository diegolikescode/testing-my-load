package repository

import (
	"database/sql"
)

// type PersonRepository interface {
// 	Create(name, apelido, nascimento string, stack []string) error
// 	FindById(id string) *Person
// 	FindByTerm(t string) *Person
// 	Count() int
// }

type Repository struct {
	db *sql.DB
}

func NewPostgresConnection() {
	// dsn := "postgres://pqgo:password@localhost/pqgo?sslmode=verify-full"
	// db, err := sql.Open("postgres", dsn)
}
