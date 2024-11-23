package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func CreateDatabaseConnExample() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=localhost port=5432 dbname=expense_tracker user=postgres password=")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
