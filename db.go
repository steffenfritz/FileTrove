package filetrove

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// CreateFileTroveDB creates a new an empty sqlite database for FileTrove.
// It contains information like configurations, sessions and db versions.
func CreateFileTroveDB(dbpath string, version string, initdate string) error {
	db, err := sql.Open("sqlite3", dbpath+"/filetrove.db")

	if err != nil {
		return err
	}

	defer db.Close()

	initstatements := `
		CREATE TABLE filetrove(version TEXT, initdate TEXT);
		INSERT INTO filetrove(version, initdate) VALUES(version, initdate);`

	_, err = db.Exec(initstatements)

	if err != nil {
		return err
	}

	return nil
}

// ConnectFileTroveDB creates a connection to an existing sqlite database.
func ConnectFileTroveDB(dbpath string) {

}
