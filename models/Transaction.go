package models

import (
	"time"
)

type Transaction struct {
	ID              int64     `json:"id"`
	AccountID       int       `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	Balance         float64   `json:"balance"`
	EventDate       time.Time `json:"event_date"`
}
