package services

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"pismo/models"
	"pismo/services"
	"pismo/store"
)

func TestCreateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := store.NewRepository(db)
	service := services.NewTransactionService(repo)

	tests := []struct {
		name           string
		transaction    models.Transaction
		mockSetup      func(sqlmock.Sqlmock)
		expectedResult int64
		expectedError  string
	}{
		{
			name: "Successful deposit transaction",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4, // Deposit
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 4, 100.0, 100.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}))
				mock.ExpectExec(`UPDATE Transactions SET balance = \? WHERE transaction_id = \?`).
					WithArgs(100.0, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			expectedResult: 1,
			expectedError:  "",
		},
		{
			name: "Successful non-deposit transaction",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 1, // Purchase
				Amount:          -50.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 1, -50.0, -50.0).
					WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit()
			},
			expectedResult: 2,
			expectedError:  "",
		},
		{
			name: "Invalid operation direction",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,      // Deposit
				Amount:          -100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {},
			expectedResult: 0,
			expectedError:  "invalid transaction amount -100.00 for the given operation type ID 4: expected Credit direction",
		},
		{
			name: "Begin transaction error",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errors.New("db error"))
			},
			expectedResult: 0,
			expectedError:  "failed to begin transaction: db error",
		},
		{
			name: "Create transaction error",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 4, 100.0, 100.0).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			expectedResult: 0,
			expectedError:  "db error",
		},
		{
			name: "Process discharge error",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 4, 100.0, 100.0).
					WillReturnResult(sqlmock.NewResult(3, 1))
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}))
				mock.ExpectExec(`UPDATE Transactions SET balance = \? WHERE transaction_id = \?`).
					WithArgs(100.0, 3).
					WillReturnError(errors.New("discharge error"))
				mock.ExpectRollback()
			},
			expectedResult: 0,
			expectedError:  "failed to update deposit transaction: discharge error",
		},
		{
			name: "Commit error",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 4, 100.0, 100.0).
					WillReturnResult(sqlmock.NewResult(4, 1))
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}))
				mock.ExpectExec(`UPDATE Transactions SET balance = \? WHERE transaction_id = \?`).
					WithArgs(100.0, 4).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			expectedResult: 0,
			expectedError:  "failed to commit transaction: commit error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			result, err := service.CreateTransaction(tt.transaction)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}