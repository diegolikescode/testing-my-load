package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/diegolikescode/testing-my-load/pkg"
	_ "github.com/lib/pq"
	"github.com/phuslu/log"
)

// queries
const (
	queryCreate       = "INSERT INTO pessoas(id, apelido, nome, nascimento, stack, busca_termos) VALUES($1, $2, $3, $4, $5, $6);"
	queryFindByID     = "SELECT apelido, nome, nascimento, stack FROM pessoas WHERE id = $1;"
	queryFindByTermo  = "SELECT id, apelido, nome, nascimento, stack FROM pessoas WHERE busca_termos LIKE $1;"
	queryCheckExists  = "SELECT EXISTS(SELECT 1 FROM pessoas WHERE apelido=$1);"
	queryCountPessoas = "SELECT COUNT(*) FROM pessoas;"
)

type Repository struct {
	db        *sql.DB
	prepStmts map[string]*sql.Stmt
}

func (r *Repository) setupStmts() {
	createStmt, err := r.db.Prepare(queryCreate)
	if err != nil {
		panic(fmt.Sprintf("couldnt prepare create stmt, got=%s", err))
	}

	findByIDStmt, err := r.db.Prepare(queryFindByID)
	if err != nil {
		panic(fmt.Sprintf("couldnt prepare find by ID stmt, got=%s", err))
	}

	findByTermStmt, err := r.db.Prepare(queryFindByTermo)
	if err != nil {
		panic(fmt.Sprintf("couldnt prepare find by Term stmt, got=%s", err))
	}

	checkExistsStmt, err := r.db.Prepare(queryCheckExists)
	if err != nil {
		panic(fmt.Sprintf("couldnt prepare Check Exists stmt, got=%s", err))
	}

	countPessoasStmt, err := r.db.Prepare(queryCountPessoas)
	if err != nil {
		panic(fmt.Sprintf("couldnt prepare Count Pessoas stmt, got=%s", err))
	}

	r.prepStmts["create"] = createStmt
	r.prepStmts["findById"] = findByIDStmt
	r.prepStmts["findByTerm"] = findByTermStmt
	r.prepStmts["checkExists"] = checkExistsStmt
	r.prepStmts["countPeople"] = countPessoasStmt
}

func (r *Repository) stackStr(s []string) string {
	return strings.Join(s, ";")
}

func (r *Repository) stackSlice(s string) []string {
	return strings.Split(s, ";")
}

func (r *Repository) pessoaSearchTerm(apelido, nome, stackStr string) string {
	return apelido + ";" + nome + ";" + stackStr
}

func (r *Repository) CheckIfExists(apelido string) bool {
	exists := r.prepStmts["checkExists"].QueryRow(apelido)
	var val bool
	err := exists.Scan(&val)
	if err != nil {
		log.Error().Err(err)
	}
	return val
}

func (r *Repository) Create(id, nome, apelido, nascimento string, stack []string) error {
	stackStr := r.stackStr(stack)
	searchTerm := r.pessoaSearchTerm(apelido, nome, stackStr)
	_, err := r.prepStmts["create"].Exec(id, nome, apelido, nascimento, stackStr, searchTerm)
	if err != nil {
		// log
	}

	return nil
}

func (r *Repository) FindByID(id string) *Person {
	result := r.prepStmts["findById"].QueryRow(id)
	if result.Err() != nil {
		// log
	}

	var pessoa Person
	var stackStr string
	err := result.Scan(&pessoa.Apelido, &pessoa.Nome, &pessoa.Nascimento, &stackStr)
	if err != nil {
		// log
		return nil
	}

	pessoa.Stack = r.stackSlice(stackStr)

	return &pessoa
}

func (r *Repository) FindByTerm(t string) []*Person {
	t = "%" + t + "%"
	result, err := r.prepStmts["findByTerm"].Query(t)
	if err != nil {
		log.Error().Msg("error finding by term=" + result.Err().Error())
		return []*Person{}
	}

	pessoas := []*Person{}
	for result.Next() {
		var p Person
		var stackStr string
		err := result.Scan(&p.ID, &p.Apelido, &p.Nome, &p.Nascimento, &stackStr)
		if err != nil {
			log.Error().Msg("error parsing the result into struct=" + err.Error())
			return nil
		}

		p.Stack = r.stackSlice(stackStr)
		pessoas = append(pessoas, &p)
	}

	return pessoas
}

func (r *Repository) Count() int {
	result := r.prepStmts["countPeople"].QueryRow()
	if result.Err() != nil {
		log.Error().Msg("error counting people=" + result.Err().Error())
	}

	var count int
	err := result.Scan(&count)
	if err != nil {
		log.Error().Msg("error counting people=" + err.Error())
		return 0
	}

	return count
}

func NewRepository() *Repository {
	dbHost := pkg.GetEnvOrDieTrying("POSTGRES_HOST")
	dbUser := pkg.GetEnvOrDieTrying("POSTGRES_USER")
	dbPasswd := pkg.GetEnvOrDieTrying("POSTGRES_PASSWORD")
	dbName := pkg.GetEnvOrDieTrying("POSTGRES_DB")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPasswd, dbHost, dbName)
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	repo := &Repository{
		db:        db,
		prepStmts: make(map[string]*sql.Stmt),
	}
	repo.setupStmts()

	return repo
}
