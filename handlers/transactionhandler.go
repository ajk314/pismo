package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"pismo/models"
	"pismo/services"
)

const (
	numConcurrentTransactions = 10 // will try to create 10 concurrent transactions
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

func (h *TransactionHandler) HandleCreateTransactionRaceCondition(w http.ResponseWriter, r *http.Request) {
    var req models.Transaction
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

	// redundant checks I know... be better to put some validator function to keep code DRY
	// but just copy pasting this in here for now
    if req.AccountID == 0 {
        http.Error(w, "No Account ID provided", http.StatusBadRequest)
        return
    }
    if req.OperationTypeID == 0 {
        http.Error(w, "No Operation Type ID provided", http.StatusBadRequest)
        return
    }
    if req.Amount == 0.0 {
        http.Error(w, "No Amount provided", http.StatusBadRequest)
        return
    }

	transactionIDs, err := h.transactionService.CreateTransactionsConcurrently(req, numConcurrentTransactions)
	if err != nil {
		http.Error(w, "Error creating transaction.", http.StatusInternalServerError)
        return
	}

    if len(transactionIDs) == 0 {
        http.Error(w, "Failed to create any transactions due to race condition.", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    resp := fmt.Sprintf("Successfully created %d new transactions with IDs: %v", len(transactionIDs), transactionIDs)
    if err := json.NewEncoder(w).Encode(resp); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
