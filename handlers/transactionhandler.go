package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"pismo/models"
	"pismo/services"
)

type TransactionHandler struct {
	transactionService services.TransactionServicer
}

func NewTransactionHandler(transactionService services.TransactionServicer) *TransactionHandler {
    return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) HandleCreateTransaction(w http.ResponseWriter, r *http.Request) {
    var req models.Transaction
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    if req.AccountID == 0 {
        http.Error(w, "No Account ID provided", http.StatusBadRequest) // 400
        return
    }
	if req.OperationTypeID == 0 {
        http.Error(w, "No Operation Type ID provided", http.StatusBadRequest) // 400
        return
    }
	if req.Amount == 0.0 {
        http.Error(w, "No Amount provided", http.StatusBadRequest) // 400
        return
    }

	transactionID, err := h.transactionService.CreateTransaction(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	resp := fmt.Sprintf("successfully created new transaction with ID %d", transactionID)
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}