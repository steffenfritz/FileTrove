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

	initstatements := `CREATE TABLE filetrove(version TEXT, initdate TEXT);
					   CREATE TABLE sessionsmd(uuid TEXT, 
					   	starttime TEXT, 
					   	project TEXT,
					   	archivistname TEXT,
					   	title TEXT,
					   	creator TEXT,
					   	subject TEXT,
					   	description TEXT,
					   	dateavailable TEXT,
					   	type TEXT,
					   	format TEXT,
					   	identifier TEXT,
					   	source TEXT,
					   	relations TEXT,
					   	rights TEXT,
					   	mountpoint TEXT
					   );
					   CREATE TABLE files(fileuuid TEXT,
					   	sessionuuid TEXT,
					   	filename TEXT,
					   	filesize TEXT,
					   	filemd5 TEXT,
					   	filesha1 TEXT,
					   	filesha256 TEXT,
					   	filesha512 TEXT
					   	fileblake2b TEXT,
					   	filesffmt TEXT,
					   	filesfmime TEXT,
					   	filesfformatname TEXT,
					   	filesfformatversion TEXT,
					   	filesfidentnote TEXT,
					   	filesfidentproof TEXT,
					   	filectime TEXT,
					   	filemtime TEXT,
					   	fileatime TEXT
					   ); `

	_, err = db.Exec(initstatements)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO filetrove(version, initdate) values(?,?)", version, initdate)
	if err != nil {
		return err
	}

	return nil
}

// ConnectFileTroveDB creates a connection to an existing sqlite database.
func ConnectFileTroveDB(dbpath string) {

}
