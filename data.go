package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type ChannelMessage struct {
	ChannelId int64 `json:"channel_id"`
	MessageId int64 `json:"message_id"`
}

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
	stmt  *sql.Stmt
	query string
}

// Statements
const (
	INSERT          = "insert"
	UPDATE_HASHTAG  = "update"
	LIST_HASHTAG    = "list"
	LIST_MESSAGES   = "list_messages"
	UPDATE_MESSAGES = "update_messages"
	APPEND          = "append"

	tableName = "documentb"
)

type PostgresDatabase struct {
	Db    *sql.DB
	Stmts map[string]*stmtConfig
}

func InitDb() (*PostgresDatabase, error) {
	var db, err = sql.Open(driver, fmt.Sprintf(
		"host=%s user=%s dbname=%s password=%s sslmode=%s",
		dbhost, dbuser, dbName, dbpassword, sslmode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	var stmts = map[string]*stmtConfig{
		LIST_HASHTAG:    {query: fmt.Sprintf(`select info::jsonb->'hashtags' from "%s";`, tableName)},
		UPDATE_HASHTAG:  {query: fmt.Sprintf(`update "%s" set info=jsonb_set(info, '{hashtags, $1}', $2) returning info::jsonb->'hashtags';`, tableName)},
		LIST_MESSAGES:   {query: fmt.Sprintf(`select info::jsonb->'messages' from "%s";`, tableName)},
		UPDATE_MESSAGES: {query: fmt.Sprintf(`update "%s" set info=jsonb_set(info, '{messages, $1}', $2) returning info::jsonb->'messages';`, tableName)},
	}
	for key, value := range stmts {
		stmts[key].stmt, _ = db.Prepare(value.query)
	}

	return &PostgresDatabase{
		Db:    db,
		Stmts: stmts,
	}, nil
}

func (d *PostgresDatabase) UpdateHashtags(hashtag string, channelIdList []int64) (map[string][]int64, error) {
	var insertUser = d.Stmts[UPDATE_HASHTAG].stmt
	var channelIdListJson, err = json.Marshal(channelIdList)
	if err != nil {
		return nil, err
	}

	var RawHashtags string
	err = insertUser.QueryRow(hashtag, string(channelIdListJson)).Scan(&RawHashtags)
	if err != nil {
		return nil, err
	}

	var Hashtags = make(map[string][]int64)
	if err := json.Unmarshal([]byte(RawHashtags), &Hashtags); err != nil {
		return nil, err
	}

	return Hashtags, nil
}

func (d *PostgresDatabase) ListHashtags() (map[string][]int64, error) {
	var listUser = d.Stmts[LIST_HASHTAG].stmt

	var RawHashtags string
	if err := listUser.QueryRow().Scan(&RawHashtags); err != nil {
		return nil, err
	}

	var Hashtags = make(map[string][]int64)
	if err := json.Unmarshal([]byte(RawHashtags), &Hashtags); err != nil {
		return nil, err
	}

	return Hashtags, nil
}

func (d *PostgresDatabase) UpdateMessages(messageId string, channelRefList []ChannelMessage) (map[string][]ChannelMessage, error) {
	var insertUser = d.Stmts[UPDATE_MESSAGES].stmt
	var channelRefListJson, err = json.Marshal(channelRefList)
	if err != nil {
		return nil, err
	}

	var RawMessages string
	err = insertUser.QueryRow(messageId, string(channelRefListJson)).Scan(&RawMessages)
	if err != nil {
		return nil, err
	}

	var Messages = make(map[string][]ChannelMessage)
	if err := json.Unmarshal([]byte(RawMessages), &Messages); err != nil {
		return nil, err
	}

	return Messages, nil
}

func (d *PostgresDatabase) ListMessages() (map[string][]ChannelMessage, error) {
	var listUser = d.Stmts[LIST_MESSAGES].stmt

	var RawMessages string
	if err := listUser.QueryRow().Scan(&RawMessages); err != nil {
		return nil, err
	}

	var Messages = make(map[string][]ChannelMessage)
	if err := json.Unmarshal([]byte(RawMessages), &Messages); err != nil {
		return nil, err
	}

	return Messages, nil
}
func (d *PostgresDatabase) Close() error {
	for s := range d.Stmts {
		if err := d.Stmts[s].stmt.Close(); err != nil {
			return err
		}
	}

	return d.Db.Close()
}
