package setup

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/srisudarshanrg/go-expense-tracker/server/functions"
	"github.com/srisudarshanrg/go-expense-tracker/server/models"
)

var db *sql.DB
var session *scs.SessionManager
var data = make(map[string]interface{})

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
	status := r.URL.Query().Get("status")
	loggedOut := r.URL.Query().Get("loggedOut")
	if status != "" {
		err := RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Error: status,
		})
		if err != nil {
			log.Println(err)
		}
		return
	} else if loggedOut != "" {
		err := RenderTemplate(w, r, "login.page.tmpl", models.TemplateData{
			Info: loggedOut,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

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
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}

	session.Put(r.Context(), "link", "expenses")
	session.Put(r.Context(), "linkFilePath", "expenses.page.tmpl")

	expenseListRecieved, err := functions.GetExpenses(user.ID)
	if err != nil {
		log.Println(err)
	}

	expenseList := functions.ReverseSliceExpenseStruct(expenseListRecieved)

	expenseCategoryList, expenseCategories, expenditureAmounts, colorList, err := functions.GetExpenseCategories(user.ID)
	if err != nil {
		log.Println(err)
	}

	expenseCategoriesNew, err := json.Marshal(expenseCategories)
	if err != nil {
		log.Println(err)
	}
	expenditureAmountsNew, err := json.Marshal(expenditureAmounts)
	if err != nil {
		log.Println(err)
	}
	colorListNew, err := json.Marshal(colorList)
	if err != nil {
		log.Println(err)
	}

	data["expenseList"] = expenseList
	data["expenseCategoryList"] = expenseCategoryList
	data["expenseCategories"] = string(expenseCategoriesNew)
	data["expenditureAmounts"] = string(expenditureAmountsNew)
	data["colorList"] = string(colorListNew)

	// do the msg url checking after getting all the database data and doing all the logic
	msg := r.URL.Query().Get("msg")
	if msg != "" {
		err := RenderTemplate(w, r, "expenses.page.tmpl", models.TemplateData{
			Info: msg,
			Data: data,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = RenderTemplate(w, r, "expenses.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

func ExpensesTable(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}

	session.Put(r.Context(), "link", "expenses-table")
	session.Put(r.Context(), "linkFilePath", "expenses-table.page.tmpl")

	expenseListRecieved, err := functions.GetExpenses(user.ID)
	if err != nil {
		log.Println(err)
	}

	expenseList := functions.ReverseSliceExpenseStruct(expenseListRecieved)

	expenseCategoryList, expenseCategories, expenditureAmounts, colorList, err := functions.GetExpenseCategories(user.ID)
	if err != nil {
		log.Println(err)
	}

	expenseCategoriesNew, err := json.Marshal(expenseCategories)
	if err != nil {
		log.Println(err)
	}
	expenditureAmountsNew, err := json.Marshal(expenditureAmounts)
	if err != nil {
		log.Println(err)
	}
	colorListNew, err := json.Marshal(colorList)
	if err != nil {
		log.Println(err)
	}

	data["expenseList"] = expenseList
	data["expenseCategoryList"] = expenseCategoryList
	data["expenseCategories"] = string(expenseCategoriesNew)
	data["expenditureAmounts"] = string(expenditureAmountsNew)
	data["colorList"] = string(colorListNew)

	// do the msg url checking after getting all the database data and doing all the logic
	msg := r.URL.Query().Get("msg")
	if msg != "" {
		err := RenderTemplate(w, r, "expenses-table.page.tmpl", models.TemplateData{
			Info: msg,
			Data: data,
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	err = RenderTemplate(w, r, "expenses-table.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// Tracker is the handler for the register page
func Tracker(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}

	expenseCategories, labelsCategory, valuesCategory, _, err := functions.GetExpenseCategories(user.ID)
	if err != nil {
		log.Println(err)
	}

	labelsDate, valuesDate, err := functions.GetTotalExpenditureByDate(user.ID)
	if err != nil {
		log.Println(err)
	}

	labelsCategoryNew, err := json.Marshal(labelsCategory)
	if err != nil {
		log.Println(err)
	}
	valuesCategoryNew, err := json.Marshal(valuesCategory)
	if err != nil {
		log.Println(err)
	}

	data["expenseCategoryList"] = expenseCategories
	data["labelsDate"] = labelsDate
	data["valuesDate"] = valuesDate
	data["labelsCategory"] = string(labelsCategoryNew)
	data["valuesCategory"] = string(valuesCategoryNew)

	err = RenderTemplate(w, r, "tracker.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

func TrackerCategory(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}
	log.Println(user)

	categoryExpensesListInterface := session.Get(r.Context(), "categoryExpensesList")
	categoryNameInterface := session.Get(r.Context(), "categoryName")

	categoryExpensesList, checkCategoryExpensesList := categoryExpensesListInterface.([]models.Expense)
	categoryName, checkCategoryName := categoryNameInterface.(string)
	if !checkCategoryExpensesList || !checkCategoryName {
		http.Redirect(w, r, "tracker.page.tmpl", http.StatusSeeOther)
		return
	}

	data["categoryExpensesList"] = categoryExpensesList
	data["categoryName"] = categoryName

	RenderTemplate(w, r, "tracker-category.page.tmpl", models.TemplateData{
		Data: data,
	})
}

// Budget is the handler for the register page
func Budget(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}

	budgetList, categoriesList, expenditureAmount, budgetAmount, totalExpenditure, totalBudget, totalDifference, err := functions.GetBudgets(user.ID)
	if err != nil {
		log.Println(err)
	}

	data["budgetList"] = budgetList
	data["categoriesList"] = categoriesList
	data["expenditureAmount"] = expenditureAmount
	data["budgetAmount"] = budgetAmount
	data["totalExpenditure"] = totalExpenditure
	data["totalBudget"] = totalBudget
	data["totalDifference"] = totalDifference

	err = RenderTemplate(w, r, "budget.page.tmpl", models.TemplateData{
		Data: data,
	})
	if err != nil {
		log.Println(err)
	}
}

// Profile is the handler for the register page
func Profile(w http.ResponseWriter, r *http.Request) {
	userInterface := session.Get(r.Context(), "loggedUser")
	user, check := userInterface.(models.User)
	if !check {
		notLogged := "You have to be logged in first to access this page"
		http.Redirect(w, r, "/login?status="+notLogged, http.StatusSeeOther)
		return
	}

	err := RenderTemplate(w, r, "profile.page.tmpl", models.TemplateData{
		Data: user,
	})
	if err != nil {
		log.Println(err)
	}
}

// Logout is the handler to logout of the web app
func Logout(w http.ResponseWriter, r *http.Request) {
	session.Destroy(r.Context())
	loggedOutMsg := "You have been logged out successfully"
	http.Redirect(w, r, "/login?loggedOut="+loggedOutMsg, http.StatusSeeOther)
	log.Println("logged out")
}
