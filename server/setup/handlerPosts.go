package setup

import (
	"log"
	"net/http"

	"github.com/srisudarshanrg/go-expense-tracker/server/functions"
	"github.com/srisudarshanrg/go-expense-tracker/server/models"
	"github.com/srisudarshanrg/go-expense-tracker/server/validations"
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

func RegisterPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	username := r.Form.Get("username")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	passwordConfirm := r.Form.Get("passwordConfirm")

	// form validations
	validations.MaxLength(username, 30)
	validations.MinLength(username, 2)
	validations.ValidEmail(email)
	validations.PasswordEqualConfirmPassword(password, passwordConfirm)
	validations.UsernameExists(username)
	validations.EmailExists(email)

	// get error list
	errorList := validations.GetErrorList()
	if len(errorList) > 0 {
		RenderTemplate(w, r, "register.page.tmpl", models.TemplateData{
			Data: errorList,
		})
		return
	}

	// create user
	user, err := functions.CreateNewUser(username, email, password)
	if err != nil {
		log.Println(err)
		return
	}

	// passing user to session
	session.Put(r.Context(), "loggedUser", user)

	RenderTemplate(w, r, "register.page.tmpl", models.TemplateData{
		Info: "Your account has successfully been createdd.",
	})

	log.Println(username, email, password, passwordConfirm)
}
