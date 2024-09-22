package store

import (
	"database/sql"
	"pismo/models"
)

func (repo *Repository) GetAccountByID(id int) (models.Account, error) {
	query := "SELECT * FROM Accounts WHERE account_id = ?"
	row := repo.DB.QueryRow(query, id)

	var account models.Account
	if err := row.Scan(&account.ID, &account.DocumentNumber); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, err
		}
		return models.Account{}, err
	}

	return account, nil
}

func (repo *Repository) GetAccountByDocumentNumber(id string) (models.Account, error) {
	query := "SELECT account_id, document_number FROM Accounts WHERE document_number = ?"
	row := repo.DB.QueryRow(query, id)

	var account models.Account
	if err := row.Scan(&account.ID, &account.DocumentNumber); err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, err
		}
		return models.Account{}, err
	}

	return account, nil
}


func (repo *Repository) CreateAccount(documentNumber string) (int64, error) {
	query := "INSERT INTO Accounts (document_number) VALUES (?)"
	row, err := repo.DB.Exec(query, documentNumber)
	if err != nil {
		return 0, err
	}
	return row.LastInsertId()
}
