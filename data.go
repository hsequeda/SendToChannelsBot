package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// Database properties
var (
	driver     = os.Getenv("DRIVER")
	dbhost     = os.Getenv("DB_HOST")
	dbuser     = os.Getenv("DB_USER")
	dbName     = os.Getenv("DB_NAME")
	dbpassword = os.Getenv("DB_PASSWORD")
	sslmode    = os.Getenv("DB_SSLMODE")
)

type stmtConfig struct {
	stmt *sql.Stmt
	q    string
}

// Statements
const (
	INSERT = "insert"
	UPDATE = "update"
	LIST   = "list"
	APPEND = "append"
)

type PostgresDatabase struct {
	Db    *sql.DB
	Stmts map[string]*stmtConfig
}

func InitDb() (*PostgresDatabase, error) {
	db, err := sql.Open(driver, fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		dbhost, dbuser, dbName, dbpassword, sslmode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	stmts := map[string]*stmtConfig{
		LIST:   {q: "select info from \"document\";"},
		UPDATE: {q: "update \"document\" set info=$1;"},
	}
	for k, v := range stmts {
		stmts[k].stmt, _ = db.Prepare(v.q)
	}

	return &PostgresDatabase{
		Db:    db,
		Stmts: stmts,
	}, nil
}

func (d *PostgresDatabase) Update(data map[string][]int64) error {
	insertUser := d.Stmts[UPDATE].stmt
	rawJson, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = insertUser.Exec(string(rawJson))
	if err != nil {
		return err
	}

	return nil
}

func (d *PostgresDatabase) List() (map[string][]int64, error) {
	listUser := d.Stmts[LIST].stmt

	row := listUser.QueryRow()

	var info string
	if err := row.Scan(&info); err != nil {
		return nil, err
	}

	var result = make(map[string][]int64)
	if err := json.Unmarshal([]byte(info), &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (d *PostgresDatabase) Close() error {
	for s := range d.Stmts {
		err := d.Stmts[s].stmt.Close()
		if err != nil {
			return err
		}
	}

	return d.Db.Close()
}
