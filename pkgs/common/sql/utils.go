package sql

import (
	"database/sql"
	"fmt"

	"github.com/friendsofgo/errors"
	_ "github.com/lib/pq"
)

type PsqlDatabaseConfiguration struct {
	Name     string
	User     string
	Password string
	Host     string
	SSLmode  string
}

func NewPostgresConnPool(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to PostgreSql")
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresConnPoolFromConf(conf PsqlDatabaseConfiguration) (*sql.DB, error) {
	return NewPostgresConnPool(fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=%s", conf.Name, conf.User, conf.Password, conf.Host, conf.SSLmode))
}
