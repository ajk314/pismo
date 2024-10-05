package services

import (
	"pismo/helpers"
	"pismo/models"
	"pismo/store"

	"fmt"
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
	
	// this is the db transaction that will be used to commit to db, and rollback everything
	// in case of any failures
	tx, err := s.db.BeginTransaction()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// roll back everything if there are any errors at the end
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Insert the deposit transaction into the db using the same db transaction context
	// (not the monetary transaction)
	transaction.Balance = transaction.Amount
	transactionID, err := s.db.CreateTransactionWithTx(tx, transaction)
	if err != nil {
		return 0, err
	}
	transaction.ID = transactionID
	
	// discharge the transaction only if its a deposit
	if transaction.OperationTypeID == 4 {
		// use the same db transaction context (not monetary transaction) so we can 
		// commit or rollback all at once
		err = s.db.ProcessDischargeTransactionWithTx(tx, transaction)
		if err != nil {
			return 0, err
		}
	}

	// only commit Inserting the deposit transaction and updating the debt transactions
	// if there are no errors with updating any of those in the db
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transactionID, nil
}