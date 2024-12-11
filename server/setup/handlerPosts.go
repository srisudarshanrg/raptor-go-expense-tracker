package setup

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// ExpensesPost is the handler for the post requests from the expenses page
func ExpensesPost(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		log.Println("user not in session", user, check)
		return
	}

	link := session.Get(r.Context(), "link").(string)
	linkFilePath := session.Get(r.Context(), "linkFilePath").(string)

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	name := r.Form.Get("expenseName")
	searchKey := r.Form.Get("searchExpenseKey")
	deleteExpense := r.Form.Get("deleteExpenseID")

	if name != "" {
		category := r.Form.Get("expenseCategory")
		amount := r.Form.Get("expenseAmount")
		color := r.Form.Get("expenseColor")
		amountConverted, _ := strconv.Atoi(amount)

		err := functions.AddExpense(name, category, amountConverted, color, user.ID)
		if err != nil {
			log.Println(err)
			return
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	} else if searchKey != "" {
		searchResults, err := functions.SearchExpense(strings.ToLower(searchKey), user.ID)
		if err != nil {
			log.Println(err)
		}

		var postData = map[string]interface{}{}
		postData["searchResults"] = searchResults
		postData["searchResultsLength"] = len(searchResults)

		RenderTemplate(w, r, linkFilePath, models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if deleteExpense != "" {
		id, err := strconv.Atoi(deleteExpense)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteExpense(id)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, link, http.StatusSeeOther)
	}

	session.Remove(r.Context(), "link")
	session.Remove(r.Context(), "linkFilePath")
}

// TrackerPost is the handler for the post requests from the tracker page
func TrackerPost(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		log.Println("user not in session", user, check)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	category := r.Form.Get("category")
	searchExpenseKey := r.Form.Get("searchExpenseKey")
	deleteExpenseID := r.Form.Get("deleteExpenseID")
	date := r.Form.Get("date")
	dateRangeStart := r.Form.Get("dateRangeStart")
	dateRangeEnd := r.Form.Get("dateRangeEnd")

	if category != "" {
		categoryExpensesList, err := functions.GetExpensesByCategory(category, user.ID)
		if err != nil {
			log.Println(err)
		}
		session.Put(r.Context(), "categoryExpensesList", categoryExpensesList)
		session.Put(r.Context(), "categoryName", category)

		http.Redirect(w, r, "/tracker-category", http.StatusSeeOther)
	} else if searchExpenseKey != "" {
		searchResults, err := functions.SearchExpense(searchExpenseKey, user.ID)
		if err != nil {
			log.Println(err)
		}

		var postData = map[string]interface{}{}
		postData["searchResults"] = searchResults
		postData["searchResultsLength"] = len(searchResults)

		RenderTemplate(w, r, "/tracker.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	} else if deleteExpenseID != "" {
		deleteExpenseIDConverted, err := strconv.Atoi(deleteExpenseID)
		if err != nil {
			log.Println(err)
		}
		err = functions.DeleteExpense(deleteExpenseIDConverted)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/tracker", http.StatusSeeOther)
	} else if date != "" {
		dateTimeConverted, err := time.Parse("2006-01-02", date)
		if err != nil {
			log.Println(err)
		}

		dateConvertedLayout := dateTimeConverted.Format("02-01-2006")
		dateSearchResults, err := functions.SearchExpensesByDate(dateConvertedLayout, user.ID)
		if err != nil {
			log.Println(err)
		}

		newPostData := map[string]interface{}{}
		if len(dateSearchResults) == 0 {
			msg := "No expenses found on " + dateConvertedLayout
			RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{
				Data: data,
				Info: msg,
			})

			return
		}

		newPostData["dateSearchResults"] = dateSearchResults
		newPostData["dateSearchResultsLength"] = len(dateSearchResults)
		newPostData["dateSearchResultsDate"] = dateConvertedLayout
		RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: newPostData,
		})
	} else if dateRangeStart != "" && dateRangeEnd != "" {
		dateRangeStartConverted, err := time.Parse("2006-01-02", dateRangeStart)
		if err != nil {
			log.Println(err)
		}
		dateRangeEndConverted, err := time.Parse("2006-01-02", dateRangeEnd)
		if err != nil {
			log.Println(err)
		}

		log.Println(dateRangeStartConverted)
		log.Println(dateRangeEndConverted)

		searchResults, err := functions.SearchExpensesByDateRange(dateRangeStartConverted, dateRangeEndConverted, user.ID)
		if err != nil {
			log.Println(err)
		}

		log.Println(searchResults)

		if len(searchResults) == 0 {
			msg := "No expenses found between " + dateRangeStart + " and " + dateRangeEnd
			RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{
				Info: msg,
				Data: data,
			})
			return
		}

		postData := map[string]interface{}{}
		postData["dateRangeSearchResults"] = searchResults
		postData["dateRangeSearchResultsLength"] = len(searchResults)
		postData["dateRangeSearchResultsStart"] = dateRangeStart
		postData["dateRangeSearchResultsEnd"] = dateRangeEnd

		RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{
			Data:     data,
			PostData: postData,
		})
	}
}

// TrackerCategoryPost is the handler for the post requests from the tracker category page
func TrackerCategoryPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	deleteExpenseID := r.Form.Get("deleteExpenseID")

	deleteExpenseIDConverted, err := strconv.Atoi(deleteExpenseID)
	if err != nil {
		log.Println(err)
	}

	if deleteExpenseID != "" {
		err = functions.DeleteExpense(deleteExpenseIDConverted)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/tracker-category", http.StatusSeeOther)
	}
}

// BudgetPost is the handler for the post requests from the budget page
func BudgetPost(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		log.Println("user not in session", user, check)
		return
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	budgetCategory := r.Form.Get("budgetCategory")
	budgetCategoryDelete := r.Form.Get("budgetCategoryDelete")

	if budgetCategory != "" {
		budgetAmount := r.Form.Get("budgetAmount")
		budgetAmountConverted, err := strconv.Atoi(budgetAmount)
		if err != nil {
			log.Println(err)
		}

		err = functions.AddBudget(budgetCategory, budgetAmountConverted, user.ID)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, "/budget", http.StatusSeeOther)
	} else if budgetCategoryDelete != "" {
		err = functions.DeleteBudget(budgetCategoryDelete, user.ID)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/budget", http.StatusSeeOther)
	}
}
