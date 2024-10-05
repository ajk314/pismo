package services

import (
	"pismo/helpers"
	"pismo/models"
	"pismo/store"
)

type TransactionServicer interface {
	CreateTransaction(transaction models.Transaction) (int64, error)
}

type TransactionService struct {
	db store.Repositoryer
}

func NewTransactionService(db store.Repositoryer) TransactionServicer {
	return &TransactionService{db: db}
}

func (s *TransactionService) CreateTransaction(transaction models.Transaction) (int64, error) {
	err := helpers.ValidateOperationDirection(transaction.OperationTypeID, transaction.Amount)
	if err != nil {
		return 0, err
	}

	if transaction.OperationTypeID < 4 {
		transaction.Balance = -transaction.Amount
	} else {
		transaction.Balance = transaction.Amount
	}
	
	transactionID, err := s.db.CreateTransaction(transaction)
	if err != nil {
		return 0, err
	}
	transaction.ID = transactionID
	
	err = s.db.DischargeTransaction(transaction)
	if err != nil {
		return 0.0, err
	}

	return transactionID, nil
}
