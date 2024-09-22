package store

import (
	"pismo/models"
)

func (repo *Repository) CreateTransaction(t models.Transaction) (int64, error) {
    query := "INSERT INTO Transactions (account_id, operation_type_id, amount) VALUES (?, ?, ?)"
    row, err := repo.DB.Exec(query, t.AccountID, t.OperationTypeID, t.Amount)
    if err != nil {
        return 0, err
    }
    return row.LastInsertId()
}
