package services

import (
	"errors"
	"testing"
	
	"github.com/stretchr/testify/assert"

	"pismo/mocks"
	"pismo/models"
	"pismo/services"
)

func TestCreateTransaction(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	service := services.NewTransactionService(mockRepo)

	invalidTransaction := models.Transaction{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          100.0,
	}

	validTransaction := models.Transaction{
		AccountID:       1,
		OperationTypeID: 2,
		Amount:          -100.0,
	}

	tests := []struct {
		name           string
		transaction    models.Transaction
		mockCalls      func()
		expectedResult int64
		expectedError  error
	}{
		{
			name: "Unmatching operation type and transaction value provided",
			transaction: invalidTransaction,
			mockCalls: func() {},
			expectedResult: 0,
			expectedError:  errors.New("invalid transaction amount 100.00 for the given operation type ID 2: expected Debit direction"),
		},
		{
			name: "Db error when creating transaction",
			transaction: validTransaction,
			mockCalls: func() {
				mockRepo.On("CreateTransaction", validTransaction).Return(int64(0), errors.New("some db error"))
			},
			expectedResult: 0,
			expectedError:  errors.New("some db error"),
		},
		{
			name: "Happy path: successfully created transaction",
			transaction: validTransaction,
			mockCalls: func() {
				mockRepo.On("CreateTransaction", validTransaction).Return(int64(1), nil)
			},
			expectedResult: 1,
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.mockCalls()

			result, err := service.CreateTransaction(tt.transaction)

			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			mockRepo.AssertExpectations(t)
		})
	}
}
