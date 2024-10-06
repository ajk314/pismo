package services

import (
	"database/sql"
	"errors"
	"testing"

	"pismo/mocks"
	"pismo/models"
	"pismo/services"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountByID(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := services.NewAccountService(mockRepo)

	tests := []struct {
		name           string
		accountID      int
		mockResponse   models.Account
		mockError      error
		mockCalls      func()
		expectedResult models.Account
		expectedError  error
	}{
		{
			name:      "Account not found",
			accountID: 2,
			mockCalls: func() {
				mockRepo.On("GetAccountByID", 2).Return(models.Account{}, sql.ErrNoRows)
			},
			expectedResult: models.Account{},
			expectedError:  sql.ErrNoRows,
		},
		{
			name:      "Database error",
			accountID: 2,
			mockCalls: func() {
				mockRepo.On("GetAccountByID", 2).Return(models.Account{}, errors.New("some db error"))
			},
			expectedResult: models.Account{},
			expectedError:  errors.New("some db error"),
		},
		{
			name:      "Successfully fetched account",
			accountID: 1,
			mockCalls: func() {
				mockRepo.On("GetAccountByID", 1).Return(models.Account{ID: 1, DocumentNumber: "123456789"}, nil)
			},
			expectedResult: models.Account{ID: 1, DocumentNumber: "123456789"},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Clear previous expectations
			tt.mockCalls()

			result, err := service.GetAccountByID(tt.accountID)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreateAccount(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := services.NewAccountService(mockRepo)

	tests := []struct {
		name           string
		documentNumber string
		mockCalls      func()
		expectedResult int64
		expectedError  error
	}{
		{
			name:           "Db error when checking if account already exists",
			documentNumber: "123456789",
			expectedResult: 0,
			expectedError:  errors.New("some db error"),
			mockCalls: func() {
				mockRepo.On("GetAccountByDocumentNumber", "123456789").Return(models.Account{}, errors.New("some db error"))
			},
		},
		{
			name:           "Account with same document number already exists",
			documentNumber: "123456789",
			expectedResult: 0,
			expectedError:  errors.New("an account with that document number already exists"), // Make sure this matches exactly with your actual error message
			mockCalls: func() {
				mockRepo.On("GetAccountByDocumentNumber", "123456789").Return(models.Account{ID: 1, DocumentNumber: "123456789"}, nil)
			},
		},
		{
			name:           "Db error when creating account",
			documentNumber: "987654321",
			mockCalls: func() {
				mockRepo.On("GetAccountByDocumentNumber", "987654321").Return(models.Account{}, sql.ErrNoRows)
				mockRepo.On("CreateAccount", "987654321").Return(int64(0), errors.New("database error during creation"))
			},
			expectedResult: 0,
			expectedError:  errors.New("database error during creation"),
		},
		{
			name:           "Successfully created account",
			documentNumber: "987654321",
			mockCalls: func() {
				mockRepo.On("GetAccountByDocumentNumber", "987654321").Return(models.Account{}, sql.ErrNoRows)
				mockRepo.On("CreateAccount", "987654321").Return(int64(1), nil)
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // Clear previous expectations

			tt.mockCalls()
			result, err := service.CreateAccount(tt.documentNumber)

			assert.Equal(t, tt.expectedResult, result)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
