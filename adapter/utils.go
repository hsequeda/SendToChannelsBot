package adapter

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func NewPostgresConnPool() (*sql.DB, error) {
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbSslmode := os.Getenv("DB_SSLMODE")
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=%s", dbName, dbUser, dbPassword, dbHost, dbSslmode))
	if err != nil {
		return nil, fmt.Errorf("cannot connect to PostgreSql: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
