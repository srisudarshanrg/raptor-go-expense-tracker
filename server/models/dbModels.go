package models

import "time"

// User is the model object for a user from the database
type User struct {
	ID        int
	Username  string
	Email     string
	Password  string
	JoinDate  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Expense is the model object for an expense from the database
type Expense struct {
	ID        int
	Name      string
	Category  string
	Amount    int
	UserID    int
	Date      string
	CreatedAt string
	UpdatedAt string
}

// ExpenseCategory is the model object for distinct categories with total expenditure and color
type ExpenseCategory struct {
	Category         string
	TotalExpenditure int
	Color            string
}

// Budget is the model object for a budget from the database
type Budget struct {
	ID        int
	Category  string
	Amount    string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Colors is the model object for a category color from the database
type Colors struct {
	ID        int
	Color     string
	Category  string
	UserID    int
	CreatedAt time.Time
	UpdatedAt time.Time
}
