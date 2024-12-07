package functions

import (
	"log"

	"github.com/srisudarshanrg/go-expense-tracker/server/models"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes an input password
func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return string(hashPassword), nil
}

// CheckPasswordHash checks if input password is equal to the hashed password in the database
func CheckPasswordHash(password string, passwordHash string) bool {
	check := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return check == nil
}

// ReverseSlice reverses a given slice of type models.Expense
func ReverseSliceExpenseStruct(slice []models.Expense) []models.Expense {
	for i, j := 0, len(slice)-1; i < j; {
		slice[i], slice[j] = slice[j], slice[i]
		i++
		j--
	}
	return slice
}
