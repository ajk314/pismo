package services

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"sync"
	
	"github.com/go-sql-driver/mysql"

	"pismo/helpers"
	"pismo/models"
	"pismo/store"
)

const (
    ErrCodeDeadlock = 1213
)

type TransactionServicer interface {
	CreateTransaction(transaction models.Transaction) (int64, error)
	CreateTransactionsConcurrently(req models.Transaction, count int) ([]int64, error)
}

type TransactionService struct {
	db store.Repositoryer
}

func NewTransactionService(db store.Repositoryer) TransactionServicer {
	return &TransactionService{db: db}
}

func (s *TransactionService) CreateTransactionsConcurrently(req models.Transaction, numTransactions int) ([]int64, error) {
    var wg sync.WaitGroup
    var transactionIDs []int64
    var err error

    create := func() {
        defer wg.Done()
        time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

        transactionID, tempErr := s.CreateTransaction(req)
        if tempErr != nil {
            err = tempErr
            return
        }

        transactionIDs = append(transactionIDs, transactionID)
    }

    for i := 0; i < numTransactions; i++ {
        wg.Add(1)
        go create()
    }

    wg.Wait()

    if err != nil {
        return nil, err
    }
    return transactionIDs, nil
}

func (s *TransactionService) CreateTransaction(transaction models.Transaction) (int64, error) {
	var transactionID int64
	err := helpers.ValidateOperationDirection(transaction.OperationTypeID, transaction.Amount)
	if err != nil {
		return 0, err
	}

	// there are cases where theres no race conditions but a transaction fails to execute due to deadlocks
	// this gives three attempts to create a transaction. For 10 concurrent transaction, this timeout is
	// plenty to make sure all three transactions are met. In a prod scenario, it may not be so simple.
	// for example, if 1000 rows are selected, they would all be locked by the FOR UPDATE sql statement,
	// increasing deadlocks
	for i := 0; i < 3; i++ {
		transactionID, err = s.attemptTransactionCreationWithRollback(transaction)
		if err != nil {
			var mysqlErr *mysql.MySQLError
			// errors are wrapped with my custom error messages, so I specifically need to use errors.As
			if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeDeadlock {
				fmt.Printf("Deadlock detected, retrying transaction: attempt %d\n", i+1)
				time.Sleep(time.Duration(i+1) * time.Second) // exponential back-off
				continue
			}
			return 0, err
		}
		break // transaction was finally processed
	}
	return transactionID, err
}

func (s *TransactionService) attemptTransactionCreationWithRollback(transaction models.Transaction) (int64, error) {
	// this is the db transaction that will be used to commit to db, and rollback everything
	// in case of any failures
	tx, err := s.db.BeginTransaction()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
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
		err = s.db.ProcessDischargeTransactionWithTx(tx, transaction)
		if err != nil {
			return 0, err
		}
	}

	// only commit if there are no errors with updating any of those in the db
	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transactionID, nil
}
