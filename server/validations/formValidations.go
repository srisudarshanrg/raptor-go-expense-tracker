package validations

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/asaskevich/govalidator"
)

var errorList []string
var db *sql.DB
var session *scs.SessionManager

func DBAccessFormValidations(dbAccess *sql.DB) {
	db = dbAccess
}

// SessionAccessValidations provides the validations package with access to the sessions
func SessionAccessValidations(sessionAccess *scs.SessionManager) {
	session = sessionAccess
}

// Length validates if the input has more characters than the minimum length or not
func MinLength(toCheck string, requiredMinLength int) {
	if len(toCheck) < requiredMinLength {
		errorString := "Length of input should be equal or more than" + strconv.Itoa(requiredMinLength) + "characters."
		errorList = append(errorList, errorString)
	}
}

// Length validates if the input has less characters than the maximum length or not
func MaxLength(toCheck string, requiredMaxLength int) {
	if len(toCheck) > requiredMaxLength {
		errorString := "Length of input should be less than or equal to" + strconv.Itoa(requiredMaxLength) + "characters."
		errorList = append(errorList, errorString)
	}
}

// IsEmail checks if the email inputted by the user is valid or not
func ValidEmail(email string) {
	if !govalidator.IsEmail(email) {
		errorString := email + "is not a valid email address."
		errorList = append(errorList, errorString)
	}
}

// PasswordEqualConfirmPassword checks if the password entered is equal to the confirmed password
func PasswordEqualConfirmPassword(password string, confirmPassword string) {
	if password != confirmPassword {
		errorString := "Password entered should be equal to confirmed password."
		errorList = append(errorList, errorString)
	}
}

// UserExists checks if a user already exists with same username in database
func UsernameExists(username string) {
	query := `select * from users where username=$1`
	results, err := db.Exec(query, username)
	if err != nil {
		errorString := "Could not check database to validate unique username."
		errorList = append(errorList, errorString)
	}

	exists, err := results.RowsAffected()
	if err != nil {
		errorString := "Could not check database to validate unique username."
		errorList = append(errorList, errorString)
	}

	if exists > 0 {
		errorString := "This username already exists. Please choose another one."
		errorList = append(errorList, errorString)
	}
}

// EmailExists checks if a user already exists with same email or email in database
func EmailExists(email string) {
	query := `select * from users where email=$1`
	results, err := db.Exec(query, email)
	exists, _ := results.RowsAffected()

	if err != nil {
		log.Println(err)
		errorString := "Could not check database to validate unique email."
		errorList = append(errorList, errorString)
	}

	if exists > 0 {
		errorString := "This email already has an account."
		errorList = append(errorList, errorString)
	}
}

// GetErrorList passes the error list to the final validation
func ReturnErrorList(ctx context.Context) {
	session.Put(ctx, "errorList", errorList)
	errorList = nil
}
