package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"pismo/database"
	"pismo/handlers"
	"pismo/services"
	"pismo/store"
)

func main() {
	conn, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	db := store.NewRepository(conn)

	r := mux.NewRouter()

	accountService := services.NewAccountService(db)
	accountHandler := handlers.NewAccountHandler(accountService)

	transactionService := services.NewTransactionService(db)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	r.HandleFunc("/accounts/{id}", accountHandler.HandleGetAccount).Methods("GET")
	r.HandleFunc("/accounts", accountHandler.HandleCreateAccount).Methods("POST")
	r.HandleFunc("/transactions", transactionHandler.HandleCreateTransaction).Methods("POST")
	r.HandleFunc("/transactions-race-condition", transactionHandler.HandleCreateTransactionRaceCondition).Methods("POST")

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
