package services

import (
	"errors"
	"database/sql"

	"pismo/models"
	"pismo/store"
)

type AccountServicer interface {
	GetAccountByID(id int) (models.Account, error)
	CreateAccount(documentNumber string) (int64, error)
}

type AccountService struct {
	db store.Repositoryer
}

func NewAccountService(db store.Repositoryer) AccountServicer {
	return &AccountService{db: db}
}

func (s *AccountService) GetAccountByID(id int) (models.Account, error) {
	account, err := s.db.GetAccountByID(id)
	if err != nil {
		return models.Account{}, err
	}
	return account, nil
}

func (s *AccountService) CreateAccount(documentNumber string) (int64, error) {
	// are document numbers unique?
	// if they are, we first need to check that we arent creating another account with the same document number
	account, err := s.db.GetAccountByDocumentNumber(documentNumber)
	if err != nil && err != sql.ErrNoRows { // if there is a db error and not a no record found error
		return 0, err
	}
	if account != (models.Account{}) { // struct is not empty, so account with that document number already exists
		return 0, errors.New("an account with that document number already exists")
	}

	accountID, err := s.db.CreateAccount(documentNumber)
	if err != nil {
		return 0, err
	}
	return accountID, nil
}
