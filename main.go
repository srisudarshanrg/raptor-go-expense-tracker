package main

import (
	"net/http"

	"github.com/srisudarshanrg/go-expense-tracker/server"
)

const portNumber = ":8500"

func main() {
	http.HandleFunc("/login", server.Login)
	http.HandleFunc("/register", server.Register)

	http.ListenAndServe(portNumber, nil)
}
