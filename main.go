package main

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/go-expense-tracker/server/setup"
)

const portNumber = ":8500"

func main() {
	http.HandleFunc("/login", setup.Login)
	http.HandleFunc("/register", setup.Register)

	log.Println("server running on port number", portNumber)
	http.ListenAndServe(portNumber, nil)
}
