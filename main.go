package main

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/go-expense-tracker/server/database"
	"github.com/srisudarshanrg/go-expense-tracker/server/setup"
)

const portNumber = ":8500"

func main() {
	// database setup
	db, err := database.CreateDatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	setup.DBAccess(db)

	// handler calls
	http.HandleFunc("/login", setup.Login)
	http.HandleFunc("/register", setup.Register)

	log.Println("server running on port number", portNumber)
	http.ListenAndServe(portNumber, nil)
}
