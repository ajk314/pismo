package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pismo/handlers"
	"pismo/mocks"
	"pismo/models"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandleGetAccount(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	handler := handlers.NewAccountHandler(mockService)

	validResponse := models.Account{ID: 1, DocumentNumber: "123456789"}

	tests := []struct {
		name           string
		accountID      string
		mockResponse   models.Account
		mockCalls      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid non stringified int account ID will return error",
			accountID:      "abc",
			mockResponse:   models.Account{},
			mockCalls:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid account ID: abc\n",
		},
		{
			name:         "Account does not exist",
			accountID:    "2",
			mockResponse: models.Account{},
			mockCalls: func() {
				mockService.On("GetAccountByID", 2).Return(models.Account{}, sql.ErrNoRows)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Account not found\n",
		},
		{
			name:         "Db error",
			accountID:    "2",
			mockResponse: models.Account{},
			mockCalls: func() {
				mockService.On("GetAccountByID", 2).Return(models.Account{}, errors.New("some db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "some db error\n",
		},
		{
			name:         "Happy path: Successfully fetch an account",
			accountID:    "1",
			mockResponse: validResponse,
			mockCalls: func() {
				mockService.On("GetAccountByID", 1).Return(validResponse, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"document_number":"123456789"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil // clear any previous expectations
			tt.mockCalls()

			req := httptest.NewRequest(http.MethodGet, "/accounts/"+tt.accountID, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tt.accountID})

			rr := httptest.NewRecorder()
			handler.HandleGetAccount(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var accountResponse models.Account
				err := json.NewDecoder(rr.Body).Decode(&accountResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockResponse, accountResponse)
			} else {
				assert.Equal(t, tt.expectedBody, rr.Body.String())
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestHandleCreateAccount(t *testing.T) {
	mockService := new(mocks.MockAccountService)
	handler := handlers.NewAccountHandler(mockService)

	tests := []struct {
		name           string
		requestBody    string
		mockResponse   int64
		mockError      error
		expectedStatus int
		expectedBody   string
		mockCalls      func()
	}{
		{
			name:           "Invalid request payload",
			requestBody:    `{`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected EOF\n",
			mockCalls:      func() {},
		},
		{
			name:           "No document number provided",
			requestBody:    `{"document_number": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No document number provided\n",
			mockCalls:      func() {},
		},
		{
			name:           "Non stringified int for document number provided",
			requestBody:    `{"document_number": "abc123"}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid document number in payload, must be of type int/long: abc123\n",
			mockCalls:      func() {},
		},
		{
			name:           "Db error",
			requestBody:    `{"document_number": "123456789"}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "some db error\n",
			mockCalls: func() {
				mockService.On("CreateAccount", "123456789").Return(int64(0), errors.New("some db error"))
			},
		},
		{
			name:           "Happy path: Account successfully created",
			requestBody:    `{"document_number": "123456789"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   "\"successfully created new account with ID 1\"\n",
			mockCalls: func() {
				mockService.On("CreateAccount", "123456789").Return(int64(1), nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil // clear any previous expectations
			tt.mockCalls()

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			handler.HandleCreateAccount(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
