package filetrove

import (
	"database/sql"
	"testing"
)

// seedFile inserts one file row with distinct, column-name-derived values so
// a misaligned SELECT/Scan pair in webdb.go produces an obviously wrong value
// instead of silently passing.
func seedFile(t *testing.T) (db *sql.DB, fileuuid, sessionuuid string) {
	t.Helper()

	dir := t.TempDir()
	if err := CreateFileTroveDB(dir, "TEST", "23.05.1949"); err != nil {
		t.Fatalf("CreateFileTroveDB() error = %v", err)
	}

	conn, err := ConnectFileTroveDB(dir)
	if err != nil {
		t.Fatalf("ConnectFileTroveDB() error = %v", err)
	}
	t.Cleanup(func() { conn.Close() })

	prepin, err := PrepInsertFile(conn)
	if err != nil {
		t.Fatalf("PrepInsertFile() error = %v", err)
	}
	defer prepin.Close()

	fileuuid = "web_fileuuid"
	sessionuuid = "web_sessionuuid"

	_, err = prepin.Exec(
		fileuuid, sessionuuid, "web_filename", "web_filepath", "web_filenameextension",
		int64(999), "web_filemd5", "web_filesha1", "web_filesha256", "web_filesha512", "web_fileblake2b",
		"web_filesffmt", "web_filesfmime", "web_filesfformatname", "web_filesfformatversion",
		"web_filesfidentnote", "web_filesfidentproof", "web_filectime", "web_filemtime", "web_fileatime",
		"web_filebtime", "web_filensrl", 1.5, 2,
	)
	if err != nil {
		t.Fatalf("PrepInsertFile exec error = %v", err)
	}

	return conn, fileuuid, sessionuuid
}

func TestQueryFilesFieldAlignment(t *testing.T) {
	db, fileuuid, sessionuuid := seedFile(t)

	files, total, err := QueryFiles(db, sessionuuid, FileFilters{})
	if err != nil {
		t.Fatalf("QueryFiles() error = %v", err)
	}
	if total != 1 || len(files) != 1 {
		t.Fatalf("QueryFiles() returned %d/%d rows, want 1/1", len(files), total)
	}

	f := files[0]
	checks := map[string]struct{ got, want string }{
		"FileUUID":  {f.FileUUID, fileuuid},
		"Filectime": {f.Filectime, "web_filectime"},
		"Filemtime": {f.Filemtime, "web_filemtime"},
		"Fileatime": {f.Fileatime, "web_fileatime"},
		"Filebtime": {f.Filebtime, "web_filebtime"},
		"Filensrl":  {f.Filensrl, "web_filensrl"},
	}
	for field, c := range checks {
		if c.got != c.want {
			t.Errorf("QueryFiles() field %s = %q, want %q (SELECT/Scan likely misaligned)", field, c.got, c.want)
		}
	}
}

func TestGetFileDetailFieldAlignment(t *testing.T) {
	db, fileuuid, _ := seedFile(t)

	detail, err := GetFileDetail(db, fileuuid)
	if err != nil {
		t.Fatalf("GetFileDetail() error = %v", err)
	}

	checks := map[string]struct{ got, want string }{
		"Filectime": {detail.File.Filectime, "web_filectime"},
		"Filemtime": {detail.File.Filemtime, "web_filemtime"},
		"Fileatime": {detail.File.Fileatime, "web_fileatime"},
		"Filebtime": {detail.File.Filebtime, "web_filebtime"},
		"Filensrl":  {detail.File.Filensrl, "web_filensrl"},
	}
	for field, c := range checks {
		if c.got != c.want {
			t.Errorf("GetFileDetail() field %s = %q, want %q (SELECT/Scan likely misaligned)", field, c.got, c.want)
		}
	}
}
