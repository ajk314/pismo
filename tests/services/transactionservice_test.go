package services

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
	"pismo/mocks"
	"pismo/models"
	"pismo/services"
)

func TestCreateTransaction(t *testing.T) {
	mockRepo := new(mocks.MockRepository)
	// mockTx := new(mocks.MockTx) // Create the mock transaction object
	mockTx := &sql.Tx{}
	service := services.NewTransactionService(mockRepo)

	// Define test transactions
	// invalidTransaction := models.Transaction{
	// 	AccountID:       1,
	// 	OperationTypeID: 2,
	// 	Amount:          100.0,
	// }

	validTransaction := models.Transaction{
		AccountID:       1,
		OperationTypeID: 4, // Deposit
		Amount:          100.0,
	}

	tests := []struct {
		name           string
		transaction    models.Transaction
		mockCalls      func()
		expectedResult int64
		expectedError  error
	}{
		// {
		// 	name: "Unmatching operation type and transaction value provided",
		// 	transaction: invalidTransaction,
		// 	mockCalls: func() {
		// 		// No need for any DB calls here, this should fail before DB interaction
		// 	},
		// 	expectedResult: 0,
		// 	expectedError:  errors.New("invalid transaction amount 100.00 for the given operation type ID 2: expected Debit direction"),
		// },
		{
			name: "Db error when creating transaction",
			transaction: validTransaction,
			mockCalls: func() {
				mockRepo.On("BeginTransaction").Return(mockTx, nil)
				mockRepo.On("CreateTransactionWithTx", &mocks.MockTx{}, validTransaction).Return(int64(0), errors.New("db error"))
				// mockTx.On("Rollback").Return(nil)
			},
			expectedResult: 0,
			expectedError:  errors.New("db error"),
		},
		// {
		// 	name: "Happy path: successfully created transaction without discharge",
		// 	transaction: models.Transaction{
		// 		AccountID:       1,
		// 		OperationTypeID: 1, // Not a deposit
		// 		Amount:          -50.0,
		// 	},
		// 	mockCalls: func() {
		// 		mockRepo.On("BeginTransaction").Return(mockTx, nil)
		// 		mockRepo.On("CreateTransactionWithTx", mockTx, mock.Anything).Return(int64(1), nil)
		// 		mockTx.On("Commit").Return(nil)
		// 	},
		// 	expectedResult: 1,
		// 	expectedError:  nil,
		// },
		// {
		// 	name: "Transaction created and discharge processed successfully",
		// 	transaction: validTransaction,
		// 	mockCalls: func() {
		// 		mockRepo.On("BeginTransaction").Return(mockTx, nil)
		// 		mockRepo.On("CreateTransactionWithTx", mockTx, mock.Anything).Return(int64(1), nil)
		// 		mockRepo.On("ProcessDischargeTransactionWithTx", mockTx, validTransaction).Return(nil)
		// 		mockTx.On("Commit").Return(nil)
		// 	},
		// 	expectedResult: 1,
		// 	expectedError:  nil,
		// },
		// {
		// 	name: "Discharge transaction failed",
		// 	transaction: validTransaction,
		// 	mockCalls: func() {
		// 		mockRepo.On("BeginTransaction").Return(mockTx, nil)
		// 		mockRepo.On("CreateTransactionWithTx", mockTx, mock.Anything).Return(int64(1), nil)
		// 		mockRepo.On("ProcessDischargeTransactionWithTx", mockTx, validTransaction).Return(errors.New("discharge failed"))
		// 		mockTx.On("Rollback").Return(nil)
		// 	},
		// 	expectedResult: 0,
		// 	expectedError:  errors.New("discharge failed"),
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the mock expectations between tests
			mockRepo.ExpectedCalls = nil
			// mockTx.ExpectedCalls = nil

			// Set up the mock calls defined in the test case
			tt.mockCalls()

			// Call the service function
			result, err := service.CreateTransaction(tt.transaction)

			// Assert that the result and errors are as expected
			assert.Equal(t, tt.expectedResult, result)
			assert.Equal(t, tt.expectedError, err)

			// Verify that all mocked expectations were met
			mockRepo.AssertExpectations(t)
			// mockTx.AssertExpectations(t)
		})
	}
}
