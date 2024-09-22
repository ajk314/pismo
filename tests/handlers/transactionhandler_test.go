package handlers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"pismo/handlers"
	"pismo/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateTransaction(t *testing.T) {
	mockService := new(mocks.MockTransactionService)
	handler := handlers.NewTransactionHandler(mockService)

	tests := []struct {
		name           string
		requestBody    string
		mockResponse   int64
		mockError      error
		mockCalls      func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Invalid request payload",
			requestBody:    `{`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "unexpected EOF\n",
			mockCalls:      func() {},
		},
		{
			name:           "No account ID provided",
			requestBody:    `{"account_id": 0}`,
			mockCalls:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No Account ID provided\n",
		},
		{
			name:           "No operation type ID provided",
			requestBody:    `{"account_id": 1, "operation_type_id": 0}`,
			mockCalls:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No Operation Type ID provided\n",
		},
		{
			name:           "No amount provided",
			requestBody:    `{"account_id": 1, "operation_type_id": 2, "amount": 0}`,
			mockCalls:      func() {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "No Amount provided\n",
		},
		{
			name:        "Database error during transaction creation",
			requestBody: `{"account_id": 1, "operation_type_id": 2, "amount": 12.34}`,
			mockCalls: func() {
				mockService.On("CreateTransaction", mock.Anything).Return(int64(0), errors.New("some db error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "some db error\n",
		},
		{
			name:        "Happy path: Successfully create a transaction",
			requestBody: `{"account_id": 1, "operation_type_id": 2, "amount": 12.34}`,
			mockCalls: func() {
				mockService.On("CreateTransaction", mock.Anything).Return(int64(1), nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "\"successfully created new transaction with ID 1\"\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService.ExpectedCalls = nil // clear any previous expectations
			tt.mockCalls()

			req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			handler.HandleCreateTransaction(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedBody, rr.Body.String())

			mockService.AssertExpectations(t)
		})
	}
}
