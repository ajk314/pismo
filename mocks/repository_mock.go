package mocks

import (
	"pismo/models"

	"github.com/stretchr/testify/mock"
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

func (m *MockRepository) CreateTransaction(t models.Transaction) (int64, error) {
	args := m.Called(t)
	return args.Get(0).(int64), args.Error(1)
}
