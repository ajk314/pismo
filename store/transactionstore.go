package store

import (
	"fmt"

	"pismo/models"
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
	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()

	query := `SELECT transaction_id, balance FROM Transactions
        WHERE account_id = ?
        AND operation_type_id < 4
        AND balance < 0
        ORDER BY event_date ASC`

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

        absCurrentBalance := trans.Balance * -1
        if remainingDeposit > absCurrentBalance { // using absolute value for comparison to simplify logic here
            remainingDeposit -= absCurrentBalance
            trans.Balance = 0
        } else {
            // there is more debt than remaining deposit, so pay off what we can and break out
            trans.Balance += remainingDeposit // balance is negative, so "lower" the balance by adding the remaining deposit
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

    // not a redundant close since we do that in a defer
	// rows must be closed before performing updates
	rows.Close()

	// update balances for debt transactions
	updateBalanceQuery := "UPDATE Transactions SET balance = ? WHERE transaction_id = ?"
	for _, update := range updatedTransactions {
		if _, err := tx.Exec(updateBalanceQuery, update.Balance, update.ID); err != nil {
			return fmt.Errorf("failed to update balance for transaction %d: %w", update.ID, err)
		}
	}

	// update remaining balance for the deposit transaction
	if _, err := tx.Exec(updateBalanceQuery, remainingDeposit, depositTransaction.ID); err != nil {
		return fmt.Errorf("failed to update deposit transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
