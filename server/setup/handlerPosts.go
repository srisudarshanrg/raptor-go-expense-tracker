package setup

import (
	"log"
	"net/http"
	"strconv"

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
	if !check {
		RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Error: msg,
		})
		log.Println(err)
		return
	}

	session.Put(r.Context(), "loggedUser", user)

	log.Println("successfull")
	users, _ := session.Get(r.Context(), "loggedUser").(models.User)
	log.Println(users)

	http.Redirect(w, r, "/expenses?msg="+msg, http.StatusSeeOther)
}

// RegisterPost is the handler for the post requests from the register page
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

// ExpensesPost is the handler for the post requests from the expenses page
func ExpensesPost(w http.ResponseWriter, r *http.Request) {
	user, check := session.Get(r.Context(), "loggedUser").(*models.User)
	if user == nil || !check {
		log.Println("user not in session")
		return
	}

	if r.Method == "POST" {
		formName := r.Form.Get("formName")
		switch formName {
		case "addExpenseForm":
			name := r.Form.Get("expenseName")
			category := r.Form.Get("expenseCategory")
			amount := r.Form.Get("expenseAmount")
			color := r.Form.Get("expenseColor")
			amountConverted, _ := strconv.Atoi(amount)

			msg, err := functions.AddExpense(name, category, amountConverted, color, user.ID)
			if err != nil {
				log.Println(err)
				return
			}
			RenderTemplate(w, r, "expenses.page.tmpl", models.TemplateData{
				Info: msg,
			})
		case "searchExpenseForm":

		}
	}
}
