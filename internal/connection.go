package internal

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

var DBConn *sql.DB

func InitDb(environment string) error {
	path := "dev"
	if environment == "" || environment == "dev" || environment == "development" {
		path = "dev"
	}
	var err error
	// TODO: test to see if there is a file in the environment and create it from `TEMPLATE` if not.
	DBConn, err = sql.Open("sqlite3", "./internal/db/"+path+".sqlite")

	if err != nil {
		return err
	}
	return DBConn.Ping()
}
