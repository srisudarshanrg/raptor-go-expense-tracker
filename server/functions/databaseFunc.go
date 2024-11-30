package functions

import (
	"database/sql"
	"log"
	"time"

	"github.com/srisudarshanrg/go-expense-tracker/server/models"
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
