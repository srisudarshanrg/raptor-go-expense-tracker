package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-expense-tracker/server/database"
	"github.com/srisudarshanrg/go-expense-tracker/server/setup"
)

const portNumber = ":8500"

func main() {
	// session
	session := scs.New()
	session.Cookie.Persist = true
	session.Lifetime = 1 * time.Hour
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	// database setup
	db, err := database.CreateDatabaseConn()
	if err != nil {
		log.Fatal(err)
	}

	setup.DBAccess(db)

	// handler calls
	http.HandleFunc("/login", setup.Login)
	http.HandleFunc("/register", setup.Register)
	http.HandleFunc("/expenses", setup.Expenses)
	http.HandleFunc("/tracker", setup.Tracker)
	http.HandleFunc("/budget", setup.Budget)
	http.HandleFunc("/profile", setup.Budget)
	http.HandleFunc("/logout", setup.Logout)

	log.Println("server running on port number", portNumber)
	http.ListenAndServe(portNumber, nil)
}
