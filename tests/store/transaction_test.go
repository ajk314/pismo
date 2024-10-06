package store

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"pismo/models"
	"pismo/store"
)

func TestCreateTransactionWithTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &store.Repository{DB: db}

	tests := []struct {
		name           string
		transaction    models.Transaction
		mockSetup      func(sqlmock.Sqlmock)
		expectedID     int64
		expectedError  string
	}{
		{
			name: "Successful transaction creation",
			transaction: models.Transaction{
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
				Balance:         100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(1, 4, 100.0, 100.0).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedID:    1,
			expectedError: "",
		},
		{
			name: "Database error",
			transaction: models.Transaction{
				AccountID:       2,
				OperationTypeID: 1,
				Amount:          -50.0,
				Balance:         -50.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(`INSERT INTO Transactions \(account_id, operation_type_id, amount, balance\) VALUES \(\?, \?, \?, \?\)`).
					WithArgs(2, 1, -50.0, -50.0).
					WillReturnError(errors.New("db error"))
			},
			expectedID:    0,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			tx, err := db.Begin()
			assert.NoError(t, err)

			id, err := repo.CreateTransactionWithTx(tx, tt.transaction)

			if tt.expectedError != "" {
				assert.EqualError(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.expectedID, id)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestProcessDischargeTransactionWithTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := &store.Repository{DB: db}

	tests := []struct {
		name               string
		depositTransaction models.Transaction
		mockSetup          func(sqlmock.Sqlmock)
		expectedError      string
	}{
		{
			name: "Successful discharge - single transaction",
			depositTransaction: models.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}).AddRow(2, -50.0))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(0.0, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(50.0, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: "",
		},
		{
			name: "Successful discharge - multiple transactions",
			depositTransaction: models.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}).
						AddRow(2, -30.0).
						AddRow(3, -50.0).
						AddRow(4, -40.0))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(0.0, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(0.0, 3).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(-20.0, 4).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(0.0, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: "",
		},
		{
			name: "Query error",
			depositTransaction: models.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnError(errors.New("query error"))
			},
			expectedError: "failed to query transactions: query error",
		},
		{
			name: "Update error",
			depositTransaction: models.Transaction{
				ID:              1,
				AccountID:       1,
				OperationTypeID: 4,
				Amount:          100.0,
			},
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`SELECT transaction_id, balance FROM Transactions WHERE account_id = \? AND operation_type_id < 4 AND balance < 0 ORDER BY event_date ASC`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"transaction_id", "balance"}).AddRow(2, -50.0))
				mock.ExpectExec("UPDATE Transactions SET balance = \\? WHERE transaction_id = \\?").
					WithArgs(0.0, 2).
					WillReturnError(errors.New("update error"))
			},
			expectedError: "failed to update balance for transaction 2: update error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)

			tx, err := db.Begin()
			assert.NoError(t, err)

			err = repo.ProcessDischargeTransactionWithTx(tx, tt.depositTransaction)

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