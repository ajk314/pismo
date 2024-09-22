package helpers

import (
	"fmt"
)

// Direction represents the direction of a transaction (Debit or Credit)
type Direction int

const (
	Debit  Direction = -1 // Deduct money
	Credit Direction = 1  // Add money
)

// always keep this private
// this should never be mutated unless we are adding another transaction type
var operationDirectionMap = map[int]Direction{
	1: Debit,  // Normal Purchase
	2: Debit,  // Purchase with installments
	3: Debit,  // Withdrawal
	4: Credit, // Credit Voucher
}

// String returns the string representation of a Direction
func (d Direction) String() string {
	switch d {
	case Debit:
		return "Debit"
	case Credit:
		return "Credit"
	default:
		return "Unknown"
	}
}

// ValidateOperationDirection validates the direction of the transaction
func ValidateOperationDirection(operationTypeID int, transactionAmount float32) error {
	direction, ok := operationDirectionMap[operationTypeID]
	if !ok {
		return fmt.Errorf("invalid operation type ID: %d", operationTypeID)
	}

	transactionDirection := Credit
	if transactionAmount < 0.0 {
		transactionDirection = Debit
	}

	if transactionDirection != direction {
		return fmt.Errorf(
			"invalid transaction amount %.2f for the given operation type ID %d: expected %s direction",
			transactionAmount, operationTypeID, direction.String(),
		)
	}

	return nil
}
