package mocks

import (
	"github.com/stretchr/testify/mock"

	"pismo/models"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetAccountByID(id int) (models.Account, error) {
	args := m.Called(id)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockRepository) GetAccountByDocumentNumber(documentNumber string) (models.Account, error) {
	args := m.Called(documentNumber)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockRepository) CreateAccount(documentNumber string) (int64, error) {
	args := m.Called(documentNumber)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) BeginTransaction() (*MockTx, error) {
	args := m.Called()
	return args.Get(0).(*MockTx), args.Error(1)
}

func (m *MockRepository) CreateTransactionWithTx(tx *MockTx, transaction models.Transaction) (int64, error) {
	args := m.Called(transaction)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRepository) ProcessDischargeTransactionWithTx(*MockTx, models.Transaction) error {
	args := m.Called()
	return args.Error(1)
}
