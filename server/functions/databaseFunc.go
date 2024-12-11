package functions

import (
	"database/sql"
	"encoding/json"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/srisudarshanrg/go-expense-tracker/server/models"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// DBAccessFunctions provides the functions package with access to the database
func DBAccessFunctions(dbAccess *sql.DB) {
	db = dbAccess
}

// CreateNewUser creates a new user in the database
func CreateNewUser(username string, email string, password string) (models.User, error) {
	hashPassword, err := HashPassword(password)
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	currentDate := time.Now().Format("02-01-2006")

	createUserQuery := `insert into users (username, email, password, join_date, created_at, updated_at) values($1, $2, $3, $4, $5, $6)`
	_, err = db.Exec(createUserQuery, username, email, hashPassword, currentDate, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return models.User{}, err
	}

	var usernameComplete, emailComplete, passwordComplete string
	var id int
	var createdAt, updatedAt time.Time

	getFullDetailsQuery := `select * from users where username=$1`
	rows := db.QueryRow(getFullDetailsQuery, username)
	rows.Scan(&id, &usernameComplete, &emailComplete, &passwordComplete, &createdAt, &updatedAt)

	completeUser := models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return completeUser, err
}

func AuthenticateUser(credential string, passwordInput string) (bool, models.User, string, error) {
	getUserQuery := `select * from users where username=$1 or email=$1`
	result, err := db.Exec(getUserQuery, credential)
	if err != nil {
		log.Println(err)
		return false, models.User{}, "", err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return false, models.User{}, "", err
	}
	if affected == 0 {
		return false, models.User{}, "Invalid Credentials", nil
	}

	row := db.QueryRow(getUserQuery, credential)

	var id int
	var username, email, password, joinDate string
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &username, &email, &password, &joinDate, &createdAt, &updatedAt)
	if err != nil {
		log.Println(err)
		return false, models.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordInput))
	if err != nil {
		log.Println(err)
		return false, models.User{}, "Invalid Credentials", err
	}

	user := models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		Password:  password,
		JoinDate:  joinDate,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return true, user, "You have been successfully logged in!", nil
}

// GetExpenses gets all the user's expenses for a given user id
func GetExpenses(userID int) ([]models.Expense, error) {
	var expenseList []models.Expense

	getExpensesQuery := `select * from expenses where user_id=$1`
	rows, err := db.Query(getExpensesQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, userID, amount int
		var name, category, date string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &category, &amount, &date, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expense := models.Expense{
			ID:        id,
			Name:      name,
			Category:  category,
			Amount:    amount,
			Date:      date,
			UserID:    userID,
			CreatedAt: createdAt.Format("15:04"),
			UpdatedAt: updatedAt.Format("02-01-2006 15:04"),
		}
		expenseList = append(expenseList, expense)
	}

	return expenseList, nil
}

// GetExpenseCategories gets all the distinct expense categories along with their total expenditure and color
func GetExpenseCategories(userID int) ([]models.ExpenseCategory, []string, []int, []string, error) {
	getDistinctCategories := `select distinct category from expenses where user_id=$1`
	rows, err := db.Query(getDistinctCategories, userID)
	if err != nil {
		log.Println(err)
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var expenseCategoryList []models.ExpenseCategory
	var expenseCategories []string
	var expenditureAmounts []int
	var colorList []string

	for rows.Next() {
		var category string
		rows.Scan(&category)

		expenseCategories = append(expenseCategories, category)

		getAllExpensesFromCategoryQuery := `select amount from expenses where category=$1 and user_id=$2`
		row, err := db.Query(getAllExpensesFromCategoryQuery, category, userID)
		if err != nil {
			log.Println(err)
			return nil, nil, nil, nil, err
		}

		getCategoryColorQuery := `select color from colors where category=$1 and user_id=$2`
		colorRow := db.QueryRow(getCategoryColorQuery, category, userID)
		var color string
		colorRow.Scan(&color)
		colorList = append(colorList, color)

		totalExpenditure := 0
		for row.Next() {
			var amount int
			err = row.Scan(&amount)
			if err != nil {
				log.Println(err)
				return nil, nil, nil, nil, err
			}
			totalExpenditure += amount
		}
		expenditureAmounts = append(expenditureAmounts, totalExpenditure)

		categoryExpense := models.ExpenseCategory{
			Category:         category,
			TotalExpenditure: totalExpenditure,
			Color:            color,
		}
		expenseCategoryList = append(expenseCategoryList, categoryExpense)
	}

	return expenseCategoryList, expenseCategories, expenditureAmounts, colorList, nil
}

// AddExpense adds a new expense in the database and updates the category colors table
func AddExpense(name string, category string, amount int, color string, userID int) error {
	categoryUpper := strings.ToUpper(category)
	currentDate := time.Now().Format("02-01-2006")
	addExpenseQuery := `insert into expenses(name, category, amount, date, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(addExpenseQuery, name, categoryUpper, amount, currentDate, userID, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}

	searchCategoryColorQuery := `select * from colors where category=$1 and user_id=$2`
	result, err := db.Exec(searchCategoryColorQuery, category, userID)
	if err != nil {
		log.Println(err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if rowsAffected == 0 {
		insertCategoryColorQuery := `insert into colors(color, category, user_id, created_at, updated_at) values($1, $2, $3, $4, $5)`
		_, err = db.Exec(insertCategoryColorQuery, color, categoryUpper, userID, time.Now(), time.Now())
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// SearchExpense searches for a expense in the database based on a given key
func SearchExpense(key string, userID int) ([]models.Expense, error) {
	searchExpenseQuery := `select * from expenses where lower(name) like $1 and user_id=$2`
	searchExpenseArg := "%" + key + "%"
	rows, err := db.Query(searchExpenseQuery, searchExpenseArg, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var searchResults []models.Expense
	for rows.Next() {
		var name, category, date string
		var id, userID, amount int
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &category, &amount, &date, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expense := models.Expense{
			ID:        id,
			Name:      name,
			Category:  category,
			Amount:    amount,
			Date:      date,
			UserID:    userID,
			CreatedAt: createdAt.Format("02-01-2006 15:04"),
			UpdatedAt: updatedAt.Format("02-01-2006 15:04"),
		}
		searchResults = append(searchResults, expense)
	}

	return searchResults, nil
}

// DeleteExpense deletes an expense in the database
func DeleteExpense(id int) error {
	deleteExpenseQuery := `delete from expenses where id=$1`
	_, err := db.Exec(deleteExpenseQuery, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// GetExpenseByCategory gets all the expenses of a given category
func GetExpensesByCategory(category string, userID int) ([]models.Expense, error) {
	categoryConverted := strings.ToLower(category)
	getExpensesByCategoryQuery := `select * from expenses where lower(category)=$1 and user_id=$2`
	rows, err := db.Query(getExpensesByCategoryQuery, categoryConverted, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var expensesList []models.Expense

	for rows.Next() {
		var id, amount, userID int
		var name, category, date string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &category, &amount, &date, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expense := models.Expense{
			ID:        id,
			Name:      name,
			Category:  category,
			Amount:    amount,
			Date:      date,
			UserID:    userID,
			CreatedAt: createdAt.Format("15:04"),
			UpdatedAt: updatedAt.Format("02-01-2006 15:04"),
		}

		expensesList = append(expensesList, expense)
	}

	return expensesList, nil
}

// GetTotalExpenditureByDate returns total expenditure of each day from past 10 days
func GetTotalExpenditureByDate(userID int) (string, string, error) {
	var labels []string
	var values []int

	for n := 10; n >= 0; n-- {
		currentDate := time.Now()
		dayBefore := currentDate.AddDate(0, 0, -n).Format("02-01-2006")
		labels = append(labels, dayBefore)

		getRowsByDateQuery := `select amount from expenses where date=$1 and user_id=$2`
		rows, err := db.Query(getRowsByDateQuery, dayBefore, userID)
		if err != nil {
			log.Println(err)
			return "", "", err
		}
		defer rows.Close()

		var totalAmount int

		for rows.Next() {
			var amount int

			err = rows.Scan(&amount)
			if err != nil {
				log.Println(err)
				return "", "", err
			}
			totalAmount += amount
		}
		values = append(values, totalAmount)
	}

	newLabels, err := json.Marshal(labels)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	newValues, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
		return "", "", err
	}

	return string(newLabels), string(newValues), nil
}

// SearchExpensesByDate gets all the expenses for a given date
func SearchExpensesByDate(date string, userID int) ([]models.Expense, error) {
	searchExpenseByDateQuery := `select * from expenses where date=$1 and user_id=$2`
	rows, err := db.Query(searchExpenseByDateQuery, date, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var expensesList []models.Expense

	for rows.Next() {
		var id, amount, userID int
		var name, category, date string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &category, &amount, &date, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expense := models.Expense{
			ID:        id,
			Name:      name,
			Category:  category,
			Amount:    amount,
			Date:      date,
			UserID:    userID,
			CreatedAt: createdAt.Format("02-01-2006"),
			UpdatedAt: updatedAt.Format("02-01-2006"),
		}

		expensesList = append(expensesList, expense)
	}

	return expensesList, nil
}

func SearchExpensesByDateRange(startDate, endDate time.Time, userID int) ([]models.Expense, error) {
	searchExpensesByDateRangeQuery := `select * from expenses where created_at>=$1 and created_at<=$2 and user_id=$3`
	rows, err := db.Query(searchExpensesByDateRangeQuery, startDate, endDate, userID)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var searchResults []models.Expense
	for rows.Next() {
		var id, userID, amount int
		var name, category, date string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &name, &category, &amount, &date, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		expense := models.Expense{
			ID:        id,
			Name:      name,
			Category:  category,
			Amount:    amount,
			Date:      date,
			UserID:    userID,
			CreatedAt: createdAt.Format("02-01-2006"),
			UpdatedAt: updatedAt.Format("02-01-2006"),
		}

		searchResults = append(searchResults, expense)
	}

	return searchResults, nil
}

// GetBudgets gets all the budgets defined by the user from the database
func GetBudgets(userID int) ([]models.Budget, string, string, string, int, int, int, error) {
	getBudgetQuery := `select * from budget where user_id=$1`
	rows, err := db.Query(getBudgetQuery, userID)
	if err != nil {
		log.Println(err)
		return nil, "", "", "", 0, 0, 0, err
	}
	defer rows.Close()

	var budgetList []models.Budget
	var categoriesList []string
	var budgetAmount []int
	var expenditureAmount []int

	for rows.Next() {
		var id, amount, userID int
		var category string
		var createdAt, updatedAt time.Time

		err = rows.Scan(&id, &category, &amount, &userID, &createdAt, &updatedAt)
		if err != nil {
			log.Println(err)
			return nil, "", "", "", 0, 0, 0, err
		}

		getTotalExpenditureQuery := `select category, amount from expenses where category=$1 and user_id=$2`
		rowsExpenditure, err := db.Query(getTotalExpenditureQuery, category, userID)
		if err != nil {
			log.Println(err)
			return nil, "", "", "", 0, 0, 0, err
		}

		totalAmount := 0
		for rowsExpenditure.Next() {
			var expenditureAmount int
			var expenditureCategory string

			err = rowsExpenditure.Scan(&expenditureCategory, &expenditureAmount)
			if err != nil {
				log.Println(err)
				return nil, "", "", "", 0, 0, 0, err
			}
			totalAmount += expenditureAmount
			if !slices.Contains(categoriesList, expenditureCategory) {
				categoriesList = append(categoriesList, expenditureCategory)
			}
		}
		expenditureAmount = append(expenditureAmount, totalAmount)
		budgetAmount = append(budgetAmount, amount)

		difference := totalAmount - amount

		var color string
		if difference <= 0 {
			color = "#198754"
		} else {
			color = "#dc3545"
		}

		budget := models.Budget{
			ID:          id,
			Category:    category,
			Amount:      amount,
			Expenditure: totalAmount,
			Difference:  amount - totalAmount,
			Color:       color,
			UserID:      userID,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}

		budgetList = append(budgetList, budget)
	}

	totalExpenditure := 0
	totalBudget := 0

	for _, i := range expenditureAmount {
		totalExpenditure += i
	}

	for _, i := range budgetAmount {
		totalBudget += i
	}

	categoriesListJSON, _ := json.Marshal(categoriesList)
	expenditureAmountJSON, _ := json.Marshal(expenditureAmount)
	budgetAmountJSON, _ := json.Marshal(budgetAmount)

	categoriesListConverted := string(categoriesListJSON)
	expenditureAmountConverted := string(expenditureAmountJSON)
	budgetAmountConverted := string(budgetAmountJSON)

	return budgetList, categoriesListConverted, expenditureAmountConverted, budgetAmountConverted, totalExpenditure, totalBudget, totalBudget - totalExpenditure, nil
}

// AddBudget adds a budget to the database
func AddBudget(category string, amount int, userID int) error {
	categoryUpper := strings.ToUpper(category)
	alreadyExistsQuery := `select * from budget where category=$1 and user_id=$2`
	result, err := db.Exec(alreadyExistsQuery, categoryUpper, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		log.Println(err)
		return err
	}

	if affected > 0 {
		updateBudgetQuery := `update budget set amount=$1, updated_at=$2 where category=$3 and user_id=$4`
		_, err = db.Exec(updateBudgetQuery, amount, time.Now(), categoryUpper, userID)
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	}

	createBudgetQuery := `insert into budget(category, amount, user_id, created_at, updated_at) values($1, $2, $3, $4, $5)`
	_, err = db.Exec(createBudgetQuery, categoryUpper, amount, userID, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// DeleteBudget deletes a budget from the database
func DeleteBudget(category string, userID int) error {
	categoryConverted := strings.ToUpper(category)
	deleteBudgetQuery := `delete from budget where category=$1 and user_id=$2`
	_, err := db.Exec(deleteBudgetQuery, categoryConverted, userID)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
