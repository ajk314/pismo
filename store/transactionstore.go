package store

import (
	"fmt"
	"math"

	"database/sql"

	"pismo/models"
)

func (repo *Repository) BeginTransaction() (*sql.Tx, error) {
	return repo.DB.Begin()
}

func (repo *Repository) CreateTransactionWithTx(tx *sql.Tx, t models.Transaction) (int64, error) {
	query := "INSERT INTO Transactions (account_id, operation_type_id, amount, balance) VALUES (?, ?, ?, ?)"
	row, err := tx.Exec(query, t.AccountID, t.OperationTypeID, t.Amount, t.Balance)
	if err != nil {
		return 0, err
	}
	return row.LastInsertId()
}

func (repo *Repository) ProcessDischargeTransactionWithTx(tx *sql.Tx, depositTransaction models.Transaction) error {
	query := `SELECT transaction_id, balance FROM Transactions
        WHERE account_id = ?
        AND operation_type_id < 4
        AND balance < 0
        ORDER BY event_date ASC
    `

	rows, err := tx.Query(query, depositTransaction.AccountID)
	if err != nil {
		return fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	remainingDeposit := depositTransaction.Amount
	var updatedTransactions []models.Transaction

	for rows.Next() {
		var trans models.Transaction
		if err := rows.Scan(&trans.ID, &trans.Balance); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}

		absCurrentBalance := math.Abs(trans.Balance)
		if remainingDeposit > absCurrentBalance {
			remainingDeposit -= absCurrentBalance
			trans.Balance = 0
		} else {
			trans.Balance += remainingDeposit
			remainingDeposit = 0
		}

		updatedTransactions = append(updatedTransactions, trans)

		if remainingDeposit <= 0 {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error during row iteration: %w", err)
	}

	// not redundant close, must close here before running INSERTS or UPDATES
	rows.Close()

	// UPDATE balances for the transactions
	updateBalanceQuery := "UPDATE Transactions SET balance = ? WHERE transaction_id = ?"
	for _, update := range updatedTransactions {
		if _, err := tx.Exec(updateBalanceQuery, update.Balance, update.ID); err != nil {
			return fmt.Errorf("failed to update balance for transaction %d: %w", update.ID, err)
		}
	}

	// UPDATE the remaining balance for the deposit transaction
	if _, err := tx.Exec(updateBalanceQuery, remainingDeposit, depositTransaction.ID); err != nil {
		return fmt.Errorf("failed to update deposit transaction: %w", err)
	}

    // Test: uncomment this error to determine if the rollback is working properly after committing to db
    // return fmt.Errorf("failure")

	return nil
}
