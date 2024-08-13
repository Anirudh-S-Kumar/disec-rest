package db_server

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type ServerDB struct {
	db *sqlx.DB
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CreateDB(filename string, db_schema string, del_existing bool) (*ServerDB, error) {

	db_exists := fileExists(filename)

	if del_existing && db_exists {
		if err := os.Remove(filename); err != nil {
			return nil, fmt.Errorf("error deleting existing database: %+v", err)
		}
	}

	db_sql, err := sql.Open("sqlite3", fmt.Sprintf("file:%v?mode=rwc", filename))

	if err != nil {
		return nil, fmt.Errorf("cannot create sql.DB object: %+v", err)
	}

	DB := sqlx.NewDb(db_sql, "sqlite3")

	if del_existing || !db_exists {
		schema, err := os.ReadFile(db_schema)
		if err != nil {
			return nil, fmt.Errorf("cannot read db schema: %+v", err)
		}

		DB.MustExec(string(schema))
	}

	return &ServerDB{db: DB}, nil
}

func CreateTempDB(db_schema string) (*ServerDB, error) {
	DB := sqlx.MustOpen("sqlite3", ":memory:")
	schema, err := os.ReadFile(db_schema)

	if err != nil {
		return nil, fmt.Errorf("cannot read db schema: %+v", err)
	}

	DB.MustExec(string(schema))

	return &ServerDB{db: DB}, nil
}
