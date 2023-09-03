package filetrove

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// SessionMD holds the metadata written to table sessionsmd
type SessionMD struct {
	UUID          string
	Starttime     string
	endtime       string
	project       string
	Archivistname string
	Mountpoint    string
}

// FileMD holds the metadata for each inspected file and that is written to the table files
type FileMD struct {
	Filename            string
	Filesize            int64
	Filemd5             string
	Filesha1            string
	Filesha256          string
	Filesha512          string
	Fileblake2b         string
	Filesffmt           string
	Filesfmime          string
	Filesfformatname    string
	Filesfformatversion string
	Filesfidentnote     string
	Filesfidentproof    string
	Filectime           string
	Filemtime           string
	Fileatime           string
}

// CreateFileTroveDB creates a new an empty sqlite database for FileTrove.
// It contains information like configurations, sessions and db versions.
func CreateFileTroveDB(dbpath string, version string, initdate string) error {
	db, err := sql.Open("sqlite3", dbpath+"/filetrove.db")

	if err != nil {
		return err
	}

	defer db.Close()

	initstatements := `CREATE TABLE filetrove(version TEXT, initdate TEXT, lastupdate TEXT);
					   CREATE TABLE sessionsmd(uuid TEXT, 
					   	starttime TEXT,
					   	endtime TEXT.
					   	project TEXT,
					   	archivistname TEXT,
					   	mountpoint TEXT,
					   	
					   );
					   CREATE TABLE dublincore(uuid TEXT,
					   	title TEXT,
					   	creator TEXT,
					   	contributor TEXT,
					   	publisher TEXT,
					   	subject TEXT,
					   	description TEXT,
					   	date TEXT,
					   	language TEXT,
					   	type TEXT,
					   	format TEXT,
					   	identifier TEXT,
					   	source TEXT,
					   	relation TEXT,
					   	rights TEXT,
					   	coverage TEXT
					   );
					   CREATE TABLE files(fileuuid TEXT,
					   	sessionuuid TEXT,
					   	filename TEXT,
					   	filesize INT,
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
func ConnectFileTroveDB(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbpath+"/filetrove.db")

	if err != nil {
		return nil, err
	}

	return db, nil
}

// InsertSession adds session metadata to database
func InsertSession(db *sql.DB) {

}
