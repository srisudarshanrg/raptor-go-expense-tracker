package setup

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/go-expense-tracker/server/models"
)

// LoginPost is the handler for the post requests from the login page
func LoginPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	log.Println(username, password)

	err = RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{})
	if err != nil {
		log.Println(err)
	}
}
