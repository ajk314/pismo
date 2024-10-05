package store

import (
	"pismo/models"

    "fmt"
    "slices"
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
    query := `SELECT * 
        from Transactions
        WHERE account_id = ?
        AND operation_type_id < 4
        AND balance < 0
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

    updatedTransactionIDs := []int64{}
    remainingDeposit := depositTransaction.Balance
    for i, t := range transactions {
        if remainingDeposit == 0 {
            break
        }

        absCurrentBalance := t.Balance * -1
        if remainingDeposit > absCurrentBalance { // using absolute value for comparison to simplify logic here
            remainingDeposit -= absCurrentBalance
            transactions[i].Balance = 0
        } else {
            // there is more debt than remaining deposit, so pay off what we can and break out
            transactions[i].Balance += remainingDeposit // balance is negative, so "lower" the balance by adding the remaining deposit
            remainingDeposit = 0
        }
        updatedTransactionIDs = append(updatedTransactionIDs, t.ID)
    }

    updateBalanceQuery := "UPDATE Transactions SET balance = ? WHERE transaction_id = ?"
    // update only transactions that had a chance in balance
    for _, t := range transactions {
        if slices.Contains(updatedTransactionIDs, t.ID) {
            row, err := repo.DB.Exec(updateBalanceQuery, t.Balance, t.ID)
            if err != nil {
                return err
            }
            fmt.Println(row.LastInsertId())
        }
    }

    // update the deposit transaction with remaining balance
    row, err := repo.DB.Exec(updateBalanceQuery, remainingDeposit, depositTransaction.ID)
    if err != nil {
        return err
    }
    fmt.Println(row.LastInsertId())

    return nil
}
