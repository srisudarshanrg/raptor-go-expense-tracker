package setup

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-expense-tracker/server/models"
)

var db *sql.DB
var session *scs.SessionManager

// DBAccess provides the handlers with access to the database
func DBAccessHandlers(dbAccess *sql.DB) {
	db = dbAccess
}

// SessionAccess provides the handlers with access to the sessions
func SessionAccessHandlers(sessionAccess *scs.SessionManager) {
	session = sessionAccess
}

// Login is the handler for the login page
func Login(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Register is the handler for the register page
func Register(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "register.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Expense is the handler for the register page
func Expenses(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "expenses.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Tracker is the handler for the register page
func Tracker(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Budget is the handler for the register page
func Budget(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "budget.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Profile is the handler for the register page
func Profile(w http.ResponseWriter, r *http.Request) {
	err := RenderTemplate(w, r, "profile.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}

// Logout is the handler to logout of the web app
func Logout(w http.ResponseWriter, r *http.Request) {

}
