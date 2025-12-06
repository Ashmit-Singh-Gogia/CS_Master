package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5433 user=postgres password=postgres dbname=CS_MasterDB sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error opening DB: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Can't connect to DB: ", err)
	}
	fmt.Println("Connected to CS_MasterDB successfully!")
}
