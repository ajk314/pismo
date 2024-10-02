package store

import (
	"pismo/models"

    "fmt"
)

func (repo *Repository) CreateTransaction(t models.Transaction) (int64, error) {
    query := "INSERT INTO Transactions (account_id, operation_type_id, amount, balance) VALUES (?, ?, ?, ?)"
    row, err := repo.DB.Exec(query, t.AccountID, t.OperationTypeID, t.Amount, t.Balance)
    if err != nil {
        return 0, err
    }
    return row.LastInsertId()
}

func (repo *Repository) DischargeTransaction(depositTransaction models.Transaction) error {
    query := `SELECT balance 
        from Transactions
        WHERE account_id = ?
        AND operation_type_id < 4
        AND balance > 0
    `
    rows, err := repo.DB.Query(query, depositTransaction.AccountID)
    if err != nil {
        return err
    }

    transactions := []models.Transaction{}
    for rows.Next() {
        trans := models.Transaction{}
        err = rows.Scan(&trans.ID, &trans.AccountID, &trans.OperationTypeID, &trans.Amount, &trans.Balance, &trans.EventDate)
        if err != nil {
            return err
        }
        transactions = append(transactions, trans)
    }

    remainingDeposit := depositTransaction.Balance
    for _, t := range transactions {
        currentBalance := t.Balance
        if remainingDeposit > currentBalance {
            remainingDeposit -= t.Balance
            t.Balance = 0
        } else {
            t.Balance -= remainingDeposit
            remainingDeposit = 0
            break
        }
    }

    for _, t := range transactions {
        query := "UPDATE Transactions SET balance = ? WHERE transaction_id = ?"
        row, err := repo.DB.Exec(query, t.Balance, t.ID)
        if err != nil {
            return err
        }
        fmt.Println(row.LastInsertId())
    }

    depositQuery := "INSERT INTO Transactions (account_id, operation_type_id, amount, balance) VALUES (?, ?, ?, ?)"
    depositTransaction.Balance = remainingDeposit
    row, err := repo.DB.Exec(depositQuery, depositTransaction.AccountID, depositTransaction.OperationTypeID, depositTransaction.Amount, depositTransaction.Balance)
    if err != nil {
        return err
    }
    fmt.Println(row.LastInsertId())

    return nil
}
