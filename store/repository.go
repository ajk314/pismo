package store

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"pismo/models"
)

type Tx interface {
    Commit() error
    Rollback() error
}

type Repositoryer interface {
	GetAccountByID(id int) (models.Account, error)
	GetAccountByDocumentNumber(documentNumber string) (models.Account, error)
	CreateAccount(documentNumber string) (int64, error)
	BeginTransaction() (*sql.Tx, error)
	CreateTransactionWithTx(*sql.Tx, models.Transaction) (int64, error)
	ProcessDischargeTransactionWithTx(*sql.Tx, models.Transaction) error
}

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	return &Repository{DB: db}
}
