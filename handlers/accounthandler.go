package handlers

import (
	"database/sql"
	"encoding/json"
    "fmt"
	"net/http"
	"strconv"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"pismo/models"
	"pismo/services"
)

type AccountHandler struct {
	accountService services.AccountServicer
}

func NewAccountHandler(accountService services.AccountServicer) *AccountHandler {
    return &AccountHandler{accountService: accountService}
}

func (h *AccountHandler) HandleGetAccount(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    idString := vars["id"]
    
    idInt, err := strconv.Atoi(idString)
    if err != nil {
        msg := fmt.Sprintf("Invalid account ID: %s", idString)
        http.Error(w, msg, http.StatusBadRequest) // 400
        return
    }

    account, err := h.accountService.GetAccountByID(idInt)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            http.Error(w, "Account not found", http.StatusNotFound) // 404
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError) // 500
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    err = json.NewEncoder(w).Encode(account)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func (h *AccountHandler) HandleCreateAccount(w http.ResponseWriter, r *http.Request) {
    var req models.Account
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest) // 400
        return
    }

    if req.DocumentNumber == "" {
        http.Error(w, "No document number provided", http.StatusBadRequest) // 400
        return
    }
    _, err := strconv.Atoi(req.DocumentNumber)
    if err != nil {
        msg := fmt.Sprintf("Invalid document number in payload, must be of type int/long: %s", req.DocumentNumber)
        http.Error(w, msg, http.StatusBadRequest) // 400
        return
    }

    accountID, err := h.accountService.CreateAccount(req.DocumentNumber)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError) // 500
        return
    }

    w.Header().Set("Content-Type", "application/json")
    resp := fmt.Sprintf("successfully created new account with ID %d", accountID)
    err = json.NewEncoder(w).Encode(resp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
