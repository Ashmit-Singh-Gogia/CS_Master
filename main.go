package main

import (
	dbpkg "CS_Master/internal/db"
	routepkg "CS_Master/internal/routes"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {

	db := dbpkg.Connect()
	fmt.Println("Connected to CS_MasterDB successfully!")
	dbpkg.RunMigrations(db)
	fmt.Println("Questions table and Users table created")
	r := routepkg.SetupRouter(db)
	fmt.Println("Server running on http://localhost:8000")
	http.ListenAndServe(":8000", r)
}
