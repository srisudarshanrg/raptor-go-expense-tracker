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

	credential := r.Form.Get("credential")
	password := r.Form.Get("password")

	check, user, msg, err := functions.AuthenticateUser(credential, password)
	if err != nil {
		log.Println(err)
		return
	}
	if !check {
		RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Error: msg,
		})
		return
	}

	session.Put(r.Context(), "loggedUser", user)

	log.Println("successfull")

	err = RenderTemplate(w, r, "expenses.page.tmpl", models.TemplateData{
		Info: msg,
	})

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
		log.Println("validation problem")
		return
	}

	log.Println("user created")

	// create user
	user, err := functions.CreateNewUser(username, email, password)
	if err != nil {
		log.Println(err)
		return
	}

	// passing user to session
	session.Put(r.Context(), "loggedUser", user)

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
