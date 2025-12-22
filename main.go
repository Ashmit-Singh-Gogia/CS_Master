package main

import (
	dbpkg "CS_Master/internal/db"
	routepkg "CS_Master/internal/routes"
	"fmt"
	"log"
	"net/http"
	"time" // <--- 1. Don't forget to import "time"!

	_ "github.com/lib/pq"
)

func main() {
	// 1. Database Setup
	db := dbpkg.Connect()
	fmt.Println("Connected to CS_MasterDB successfully!")

	dbpkg.RunMigrations(db)
	fmt.Println("Questions table and Users table created")

	// 2. Router Setup
	// 'r' is your "Handler". It knows how to route requests.
	r := routepkg.SetupRouter(db)

	// 3. Define the Custom Server
	// instead of passing everything into ListenAndServe, we pack it into a struct
	s := &http.Server{
		Addr:           ":8000",          // The port to listen on
		Handler:        r,                // Pass your router 'r' here!
		ReadTimeout:    10 * time.Second, // How long to wait for the client to send the request
		WriteTimeout:   10 * time.Second, // How long to wait for us to write the response
		MaxHeaderBytes: 1 << 20,          // 1 MB limit for headers (prevents attacks)
	}

	fmt.Println("Server running on http://localhost:8000")

	// 4. Start the Server
	// We call ListenAndServe on our custom struct 's'
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
