package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func CreateDatabaseConn() (*sql.DB, error) {
	db, err := sql.Open("pgx", "host=postgresql-raptor.alwaysdata.net port=5432 dbname=raptor_expenses user=raptor password=velociraptor4796")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return db, nil
}
