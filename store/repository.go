package store

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"pismo/models"
)

type Repositoryer interface {
	GetAccountByID(id int) (models.Account, error)
	GetAccountByDocumentNumber(documentNumber string) (models.Account, error)
	CreateAccount(documentNumber string) (int64, error)
	CreateTransaction(t models.Transaction) (int64, error)
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}