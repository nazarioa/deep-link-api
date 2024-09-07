package internal

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
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

func GetLinksByFingerprint(fingerprint string) ([]LinkResponse, error) {

	// Prepare
	statement, err := DBConn.Prepare("SELECT * FROM link WHERE finger_print = ? AND deleted_at IS NULL")
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

	var possibleLinks []LinkResponse
	for rows.Next() {
		var ll = new(LinkResponse)
		err = rows.Scan(&ll.ID, &ll.Destination)
		if err != nil {
			fmt.Println(err)
		} else {
			possibleLinks = append(possibleLinks, *ll)
		}
	}

	return possibleLinks, nil
}

func GetLinksByMemberIdHash(memberIdHash string) ([]LinkResponse, error) {

	// Prepare
	statement, err := DBConn.Prepare("SELECT id, destination FROM link WHERE member_id_hash = ? AND deleted_at IS NULL")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Execute
	rows, err := statement.Query(memberIdHash)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var possibleLinks []LinkResponse
	for rows.Next() {
		var ll = new(LinkResponse)
		err = rows.Scan(&ll.ID, &ll.Destination)
		if err != nil {
			fmt.Println(err)
		} else {
			possibleLinks = append(possibleLinks, *ll)
		}
	}

	return possibleLinks, nil
}

func SaveLink(l *LinkStoreRequest) error {
	if l.MemberIdHash == "" && l.Fingerprint == "" {
		err := fmt.Errorf("missing required property")
		return err
	}

	statement, err := DBConn.Prepare("INSERT INTO link (destination, finger_print, member_id_hash, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		fmt.Println(err)
		return err
	}

	currentTime := time.Now().Format(time.RFC3339)
	_, err = statement.Exec(l.Destination, l.Fingerprint, l.MemberIdHash, currentTime)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ConsumeDeeplink(id int) (ok bool, error error) {
	statement, err := DBConn.Prepare("UPDATE link SET deleted_at = ? WHERE id = ?")
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = statement.Exec(time.Now().Format(time.RFC3339), id)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	return true, nil
}
