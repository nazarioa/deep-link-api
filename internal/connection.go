package internal

import (
	"database/sql"
	"fmt"
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

func GetLinksByFingerprint(fingerprint string) ([]Link, error) {

	// Prepare
	statement, err := DBConn.Prepare("SELECT * FROM link WHERE finger_print = ?")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Execute
	rows, err := statement.Query(fingerprint)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var possibleLinks []Link
	for rows.Next() {
		var ll = new(Link)
		err = rows.Scan(&ll.ID, &ll.Destination, &ll.Fingerprint, &ll.CreatedAt)
		if err != nil {
			fmt.Println(err)
		} else {
			possibleLinks = append(possibleLinks, *ll)
		}
	}

	return possibleLinks, nil
}
