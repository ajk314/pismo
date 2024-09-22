package mocks

import (
	"github.com/stretchr/testify/mock"
	"pismo/models"
)

type MockAccountService struct {
	mock.Mock
}

func (m *MockAccountService) GetAccountByID(id int) (models.Account, error) {
	args := m.Called(id)
	return args.Get(0).(models.Account), args.Error(1)
}

func (m *MockAccountService) CreateAccount(documentNumber string) (int64, error) {
	args := m.Called(documentNumber)
	return args.Get(0).(int64), args.Error(1)
}
