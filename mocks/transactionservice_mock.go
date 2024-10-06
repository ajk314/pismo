package mocks

import (
	"database/sql"
	
	"github.com/stretchr/testify/mock"
	
	"pismo/models"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) BeginTransaction() (*sql.Tx, error) {
	args := m.Called()
	return args.Get(0).(*sql.Tx), args.Error(1)
}

func (m *MockTransactionService) CreateTransactionWithTx(tx *sql.Tx, transaction models.Transaction) (int64, error) {
	args := m.Called(tx, transaction)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockTransactionService) ProcessDischargeTransactionWithTx(tx *sql.Tx, transaction models.Transaction) error {
	args := m.Called(tx, transaction)
	return args.Error(0)
}