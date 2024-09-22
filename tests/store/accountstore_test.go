package store

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"pismo/models"
	"pismo/store"
)

func TestGetAccountByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &store.Repository{DB: db}

	tests := []struct {
		name           string
		accountID      int
		mockSetup      func()
		expectedResult models.Account
		expectedError  error
	}{
		{
			name:      "Account not found",
			accountID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM Accounts WHERE account_id = ?").
					WithArgs(2).
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: models.Account{},
			expectedError:  sql.ErrNoRows,
		},
		{
			name:      "Database error",
			accountID: 2,
			mockSetup: func() {
				mock.ExpectQuery("SELECT \\* FROM Accounts WHERE account_id = ?").
					WithArgs(2).
					WillReturnError(errors.New("some db error"))
			},
			expectedResult: models.Account{},
			expectedError:  errors.New("some db error"),
		},
		{
			name:      "Successfully fetched account",
			accountID: 1,
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"account_id", "document_number"}).
					AddRow(1, "123456789")
				mock.ExpectQuery("SELECT \\* FROM Accounts WHERE account_id = ?").
					WithArgs(1).
					WillReturnRows(rows)
			},
			expectedResult: models.Account{ID: 1, DocumentNumber: "123456789"},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.GetAccountByID(tt.accountID)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestGetAccountByDocumentNumber(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &store.Repository{DB: db}

	tests := []struct {
		name           string
		documentNumber string
		mockSetup      func()
		expectedResult models.Account
		expectedError  error
	}{
		{
			name:           "Account not found",
			documentNumber: "123456789",
			mockSetup: func() {
				mock.ExpectQuery("SELECT account_id, document_number FROM Accounts WHERE document_number = ?").
					WithArgs("123456789").
					WillReturnError(sql.ErrNoRows)
			},
			expectedResult: models.Account{},
			expectedError:  sql.ErrNoRows,
		},
		{
			name:           "Database error",
			documentNumber: "123456789",
			mockSetup: func() {
				mock.ExpectQuery("SELECT account_id, document_number FROM Accounts WHERE document_number = ?").
					WithArgs("123456789").
					WillReturnError(errors.New("some db error"))
			},
			expectedResult: models.Account{},
			expectedError:  errors.New("some db error"),
		},
		{
			name:           "Successfully fetched account",
			documentNumber: "123456789",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"account_id", "document_number"}).
					AddRow(1, "123456789")
				mock.ExpectQuery("SELECT account_id, document_number FROM Accounts WHERE document_number = ?").
					WithArgs("123456789").
					WillReturnRows(rows)
			},
			expectedResult: models.Account{ID: 1, DocumentNumber: "123456789"},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.GetAccountByDocumentNumber(tt.documentNumber)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreateAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &store.Repository{DB: db}

	tests := []struct {
		name           string
		documentNumber string
		mockSetup      func()
		expectedResult int64
		expectedError  error
	}{
		{
			name:           "Database error during creation",
			documentNumber: "123456789",
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO Accounts \\(document_number\\) VALUES \\(\\?\\)").
					WithArgs("123456789").
					WillReturnError(errors.New("some db error"))
			},
			expectedResult: 0,
			expectedError:  errors.New("some db error"),
		},
		{
			name:           "Successfully created account",
			documentNumber: "987654321",
			mockSetup: func() {
				mock.ExpectExec("INSERT INTO Accounts \\(document_number\\) VALUES \\(\\?\\)").
					WithArgs("987654321").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			result, err := repo.CreateAccount(tt.documentNumber)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
