package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
    user := "root"
    password := "root"
	hostname := "localhost:3306" // using localhost for now
    dbname := "pismo_db"

    connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, hostname, dbname)
    conn, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}

	// health check
	err = conn.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database:", err)
	} else {
		fmt.Println("Successfully connected to the database!")
	}
	return conn, err
}
