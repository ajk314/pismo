package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"pismo/helpers"
)

func TestValidateOperationDirection(t *testing.T) {
	tests := []struct {
		name              string
		operationTypeID   int
		transactionAmount float32
		expectedError     string
	}{
		{
			name:              "Invalid operation type ID",
			operationTypeID:   99, // Invalid ID
			transactionAmount: 50.0,
			expectedError:     "invalid operation type ID: 99",
		},
		{
			name:              "Debit expected but received credit amount",
			operationTypeID:   1, // Normal Purchase
			transactionAmount: 50.0,
			expectedError:     "invalid transaction amount 50.00 for the given operation type ID 1: expected Debit direction",
		},
		{
			name:              "Credit expected but received debit amount",
			operationTypeID:   4, // Credit Voucher
			transactionAmount: -50.0,
			expectedError:     "invalid transaction amount -50.00 for the given operation type ID 4: expected Credit direction",
		},
		{
			name:              "Happy path: Valid debit transaction",
			operationTypeID:   1, // Normal Purchase
			transactionAmount: -100.0,
			expectedError:     "",
		},
		{
			name:              "Happy path: Valid credit transaction",
			operationTypeID:   4, // Credit Voucher
			transactionAmount: 200.0,
			expectedError:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helpers.ValidateOperationDirection(tt.operationTypeID, tt.transactionAmount)

			if tt.expectedError == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.expectedError)
			}
		})
	}
}
