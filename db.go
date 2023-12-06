package filetrove

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"reflect"
	"strconv"
	"text/tabwriter"
)

// SessionMD holds the metadata written to table sessionsmd
type SessionMD struct {
	UUID           string
	Starttime      string
	Endtime        string
	Project        string
	Archivistname  string
	Mountpoint     string
	ExifFlag       string
	Dublincoreflag string
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
	Filesfregistry      string
	Filectime           string
	Filemtime           string
	Fileatime           string
	Filensrl            string
	Fileentropy         float64
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
					   	endtime TEXT,
					   	project TEXT,
					   	archivistname TEXT,
					   	mountpoint TEXT,
					   	exifflag TEXT,
					   	dublincoreflag TEXT
					   );
					   CREATE TABLE dublincore(uuid TEXT,
					    sessionuuid TEXT,
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
					   	filesize INTEGER,
					   	filemd5 TEXT,
					   	filesha1 TEXT,
					   	filesha256 TEXT,
					   	filesha512 TEXT,
					   	fileblake2b TEXT,
					   	filesffmt TEXT,
					   	filesfmime TEXT,
					   	filesfformatname TEXT,
					   	filesfformatversion TEXT,
					   	filesfidentnote TEXT,
					   	filesfidentproof TEXT,
					   	filectime TEXT,
					   	filemtime TEXT,
					   	fileatime TEXT,
					   	filensrl TEXT,
					   	fileentropy INTEGER
					   ); 
					   CREATE TABLE directories(diruuid TEXT,
					    sessionuuid TEXT,
					    dirname TEXT);
                       CREATE TABLE exif(exifuuid TEXT,
                         sessionuuid TEXT,
                         fileuuid TEXT,
                         exifversion TEXT,
                         datetime TEXT,
                         datetimeorig TEXT,
                         artist TEXT,
                         copyright TEXT,
                         make TEXT,
                         xptitle TEXT,
                         xpcomment TEXT,
                         xpauthor TEXT,
                         xpkeywords TEXT,
                         xpsubject TEXT
                         );`

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

// InsertSession adds session metadata to the database
func InsertSession(db *sql.DB, s SessionMD) error {
	_, err := db.Exec("INSERT INTO sessionsmd VALUES(?,?,?,?,?,?,?,?)", s.UUID, s.Starttime, nil, s.Project,
		s.Archivistname, s.Mountpoint, s.ExifFlag, s.Dublincoreflag)

	return err
}

// InsertDC adds DublinCore metadata to the database
func InsertDC(db *sql.DB, sessionuuid string, dcuuid string, dc DublinCore) error {
	_, err := db.Exec("INSERT INTO dublincore VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", dcuuid, sessionuuid,
		dc.Title, dc.Creator, dc.Contributor, dc.Publisher, dc.Subject, dc.Description, dc.Date, dc.Language,
		dc.Type, dc.Format, dc.Identifier, dc.Source, dc.Relation, dc.Rights, dc.Coverage)

	return err
}

// PrepInsertFile prepares a statement for the addition of a single file
func PrepInsertFile(db *sql.DB) (*sql.Stmt, error) {
	prepin, err := db.Prepare("INSERT INTO files VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)")

	return prepin, err
}

// PrepInsertDir prepares a statement for the addition of a single directory
func PrepInsertDir(db *sql.DB) (*sql.Stmt, error) {
	prepin, err := db.Prepare("INSERT INTO directories VALUES(?,?,?)")

	return prepin, err
}

// ListSessions lists all sessions from the FileTrove database
func ListSessions(db *sql.DB) error {
	rows, err := db.Query("SELECT rowid, uuid, starttime, COALESCE(endtime, '') AS endtime, " +
		"COALESCE(project, '') AS project, " +
		"COALESCE(archivistname,'') AS archivistname, " +
		"COALESCE(mountpoint,'') AS mountpoint FROM sessionsmd;")
	if err != nil {
		return err
	}
	defer rows.Close() // Schließen Sie die Zeilen am Ende.

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, '.', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "ROWID\tUUID\tStart Time\tEnd Time\tProject\tArchivist Name\tMount Point")

	for rows.Next() {
		var rowid, uuid, starttime, endtime, project, archivistname, mountpoint string
		if err := rows.Scan(&rowid, &uuid, &starttime, &endtime, &project, &archivistname, &mountpoint); err != nil {
			return err
		}
		fmt.Fprintln(w, rowid+"\t"+uuid+"\t"+starttime+"\t"+endtime+"\t"+project+"\t"+archivistname+"\t"+mountpoint)
	}
	w.Flush()

	return err
}

// ExportSessionSessionTSV exports all session metadata from a session to a TSV file.
// Filtering is done by session UUID.
func ExportSessionSessionTSV(sessionuuid string) ([]string, error) {
	db, err := sql.Open("sqlite3", "db/filetrove.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	outputFile, err := os.Create(sessionuuid + "_session.tsv")
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	tsvWriter := csv.NewWriter(outputFile)
	tsvWriter.Comma = '\t'
	defer tsvWriter.Flush()

	query := "SELECT * FROM sessionsmd WHERE uuid=\"" + sessionuuid + "\""
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if err := tsvWriter.Write(columns); err != nil {
		return nil, err
	}

	// Loop to create TSV headings from row names
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	var valueStrings []string
	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		//var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "")
			} else {
				//valueStrings = append(valueStrings, value.(string))
				switch v := value.(type) {
				case int64:
					valueStrings = append(valueStrings, strconv.FormatInt(v, 10))
				case string:
					valueStrings = append(valueStrings, string(v))
				case float64:
					valueStrings = append(valueStrings, strconv.FormatFloat(v, 'E', -1, 32))
				default:
					log.Printf("Unexpected Type: %s", reflect.TypeOf(value))
					valueStrings = append(valueStrings, "")
				}
			}
		}
		if err := tsvWriter.Write(valueStrings); err != nil {
			return nil, err
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return valueStrings, nil
}

// ExportSessionFilesTSV exports all file metadata from a session to a TSV file.
// Filtering is done by session UUID.
func ExportSessionFilesTSV(sessionuuid string) error {
	db, err := sql.Open("sqlite3", "db/filetrove.db")
	if err != nil {
		return err
	}
	defer db.Close()

	outputFile, err := os.Create(sessionuuid + "_files.tsv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tsvWriter := csv.NewWriter(outputFile)
	tsvWriter.Comma = '\t'
	defer tsvWriter.Flush()

	query := "SELECT * FROM files WHERE sessionuuid=\"" + sessionuuid + "\""
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if err := tsvWriter.Write(columns); err != nil {
		return err
	}

	// Loop to create TSV headings from row names
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}
		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "")
			} else {
				//valueStrings = append(valueStrings, value.(string))
				switch v := value.(type) {
				case int64:
					valueStrings = append(valueStrings, strconv.FormatInt(v, 10))
				case string:
					valueStrings = append(valueStrings, string(v))
				case float64:
					valueStrings = append(valueStrings, strconv.FormatFloat(v, 'E', -1, 32))
				default:
					log.Printf("Unexpected Type: %s", reflect.TypeOf(value))
					valueStrings = append(valueStrings, "")
				}
			}
		}
		if err := tsvWriter.Write(valueStrings); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// ExportSessionDirectoriesTSV exports all directory metadata from a session to a TSV file.
// Filtering is done by session UUID.
func ExportSessionDirectoriesTSV(sessionuuid string) error {
	db, err := sql.Open("sqlite3", "db/filetrove.db")
	if err != nil {
		return err
	}
	defer db.Close()

	outputFile, err := os.Create(sessionuuid + "_directories.tsv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tsvWriter := csv.NewWriter(outputFile)
	tsvWriter.Comma = '\t'
	defer tsvWriter.Flush()

	query := "SELECT * FROM directories WHERE sessionuuid=\"" + sessionuuid + "\""
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if err := tsvWriter.Write(columns); err != nil {
		return err
	}

	// Loop to create TSV headings from row names
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}
		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "")
			} else {
				//valueStrings = append(valueStrings, value.(string))
				switch v := value.(type) {
				case int64:
					valueStrings = append(valueStrings, strconv.FormatInt(v, 10))
				case string:
					valueStrings = append(valueStrings, string(v))
				case float64:
					valueStrings = append(valueStrings, strconv.FormatFloat(v, 'E', -1, 32))
				default:
					log.Printf("Unexpected Type: %s", reflect.TypeOf(value))
					valueStrings = append(valueStrings, "")
				}
			}
		}
		if err := tsvWriter.Write(valueStrings); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// ExportSessionEXIFTSV exports all exif metadata from a session to a TSV file. Filtering is done by session UUID.
func ExportSessionEXIFTSV(sessionuuid string) error {
	db, err := sql.Open("sqlite3", "db/filetrove.db")
	if err != nil {
		return err
	}
	defer db.Close()

	outputFile, err := os.Create(sessionuuid + "_exif.tsv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tsvWriter := csv.NewWriter(outputFile)
	tsvWriter.Comma = '\t'
	defer tsvWriter.Flush()

	query := "SELECT * FROM exif WHERE sessionuuid=\"" + sessionuuid + "\""
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if err := tsvWriter.Write(columns); err != nil {
		return err
	}

	// Loop to create TSV headings from row names
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}
		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "")
			} else {
				//valueStrings = append(valueStrings, value.(string))
				switch v := value.(type) {
				case int64:
					valueStrings = append(valueStrings, strconv.FormatInt(v, 10))
				case string:
					valueStrings = append(valueStrings, string(v))
				case float64:
					valueStrings = append(valueStrings, strconv.FormatFloat(v, 'E', -1, 32))
				default:
					log.Printf("Unexpected Type: %s", reflect.TypeOf(value))
					valueStrings = append(valueStrings, "")
				}
			}
		}
		if err := tsvWriter.Write(valueStrings); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// ExportSessionDCTSV exports all exif metadata from a session to a TSV file. Filtering is done by session UUID.
func ExportSessionDCTSV(sessionuuid string) error {
	db, err := sql.Open("sqlite3", "db/filetrove.db")
	if err != nil {
		return err
	}
	defer db.Close()

	outputFile, err := os.Create(sessionuuid + "_dublincore.tsv")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tsvWriter := csv.NewWriter(outputFile)
	tsvWriter.Comma = '\t'
	defer tsvWriter.Flush()

	query := "SELECT * FROM dublincore WHERE sessionuuid=\"" + sessionuuid + "\""
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if err := tsvWriter.Write(columns); err != nil {
		return err
	}

	// Loop to create TSV headings from row names
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(scanArgs...); err != nil {
			return err
		}
		var valueStrings []string
		for _, value := range values {
			if value == nil {
				valueStrings = append(valueStrings, "")
			} else {
				//valueStrings = append(valueStrings, value.(string))
				switch v := value.(type) {
				case int64:
					valueStrings = append(valueStrings, strconv.FormatInt(v, 10))
				case string:
					valueStrings = append(valueStrings, string(v))
				case float64:
					valueStrings = append(valueStrings, strconv.FormatFloat(v, 'E', -1, 32))
				default:
					log.Printf("Unexpected Type: %s", reflect.TypeOf(value))
					valueStrings = append(valueStrings, "")
				}
			}
		}
		if err := tsvWriter.Write(valueStrings); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// GetImageFiles queries all files that have mime type image from a session
func GetImageFiles(db *sql.DB, sessionuuid string) (map[string]string, error) {
	var filename string
	var fileuuid string

	query := "SELECT filename, fileuuid FROM files where sessionuuid=\"" + sessionuuid +
		"\" AND filesfmime=\"image/jpeg\" OR filesfmime=\"image/tiff\""

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	imagelist := make(map[string]string)

	for rows.Next() {
		if err := rows.Scan(&filename, &fileuuid); err != nil {
			return imagelist, err
		}
		imagelist[fileuuid] = filename
	}
	if err = rows.Err(); err != nil {
		return imagelist, err
	}
	return imagelist, nil
}

// InsertExif inserts exif metadata into the FileTrove database
func InsertExif(db *sql.DB, exifuuid string, sessionid string, fileuuid string, e ExifParsed) error {
	_, err := db.Exec("INSERT INTO exif VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", exifuuid, sessionid, fileuuid, e.ExifVersion, e.DateTime, e.DateTimeOrig, e.Artist, e.Copyright, e.Make, e.XPTitle, e.XPComment, e.XPAuthor, e.XPKeywords, e.XPSubject)

	return err
}
