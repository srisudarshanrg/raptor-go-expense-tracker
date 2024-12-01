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

	createUserQuery := `insert into users (username, email, password, created_at, updated_at) values($1, $2, $3, $4, $5)`
	_, err = db.Exec(createUserQuery, username, email, hashPassword, time.Now(), time.Now())
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
	var username, email, password string
	var createdAt, updatedAt time.Time

	err = row.Scan(&id, &username, &email, &password, &createdAt, &updatedAt)
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
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return true, user, "You have been successfully logged in!", nil
}
