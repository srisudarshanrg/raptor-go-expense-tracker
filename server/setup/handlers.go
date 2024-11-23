package setup

import (
	"database/sql"
	"log"
	"net/http"
)

var db *sql.DB

func DBAccess(dbAccess *sql.DB) {
	db = dbAccess
}

// Login is the handler for the login page
func Login(w http.ResponseWriter, r *http.Request) {
	template, err := RenderTemplate("./templates/login.page.tmpl")
	if err != nil {
		log.Println(err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}

// Register is the handler for the register page
func Register(w http.ResponseWriter, r *http.Request) {
	template, err := RenderTemplate("./templates/register.page.tmpl")
	if err != nil {
		log.Println(err)
	}
	err = template.Execute(w, nil)
	if err != nil {
		log.Println(err)
	}
}
