package mocks

import (
	"pismo/models"

	"github.com/stretchr/testify/mock"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(transaction models.Transaction) (int64, error) {
	args := m.Called(transaction)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTransactionService) CreateTransactionsConcurrently(req models.Transaction, numTransactions int) ([]int64, error) {
	args := m.Called(req)
	return args.Get(0).([]int64), args.Error(1)
}
