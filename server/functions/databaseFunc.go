package functions

import (
	"database/sql"
	"log"
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
