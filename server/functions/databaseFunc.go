package functions

import (
	"database/sql"
	"log"
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}
		expenseList = append(expenseList, expense)
	}

	return expenseList, nil
}

// AddExpense adds a new expense in the database
func AddExpense(name string, category string, amount int, color string, userID int) (string, error) {
	currentDate := time.Now().Format("02-01-2006")
	addExpenseQuery := `insert into expenses(name, category, amount, date, user_id, created_at, updated_at) values($1, $2, $3, $4, $5, $6, $7)`
	_, err := db.Exec(addExpenseQuery, name, category, amount, currentDate, userID, time.Now(), time.Now())
	if err != nil {
		log.Println(err)
		return "", err
	}
	msg := "Expense" + name + "has been successfully added"
	return msg, nil
}
