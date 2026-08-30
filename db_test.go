package filetrove

import (
	"os"
	"testing"
)

func TestCreateFileTroveDB(t *testing.T) {
	type args struct {
		dbpath   string
		version  string
		initdate string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create database", args{dbpath: ".", version: "TEST", initdate: "23.05.1949"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateFileTroveDB(tt.args.dbpath, tt.args.version, tt.args.initdate); (err != nil) != tt.wantErr {
				t.Errorf("CreateFileTroveDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Remove file after testing
	err := os.Remove("filetrove.db")
	if err != nil {
		println("Could not remove test file: filetrove.db")
	}
}

// TestFileRoundTrip inserts a file row through PrepInsertFile with a distinct,
// column-name-derived value in every field and reads each column back by name.
// This guards against silent column misalignment in the positional
// VALUES(?,?,?...) INSERT statement, which has no compiler-checked mapping
// between argument position and column name.
func TestFileRoundTrip(t *testing.T) {
	dir := t.TempDir()
	if err := CreateFileTroveDB(dir, "TEST", "23.05.1949"); err != nil {
		t.Fatalf("CreateFileTroveDB() error = %v", err)
	}

	db, err := ConnectFileTroveDB(dir)
	if err != nil {
		t.Fatalf("ConnectFileTroveDB() error = %v", err)
	}
	defer db.Close()

	prepin, err := PrepInsertFile(db)
	if err != nil {
		t.Fatalf("PrepInsertFile() error = %v", err)
	}
	defer prepin.Close()

	want := map[string]string{
		"fileuuid":            "val_fileuuid",
		"sessionuuid":         "val_sessionuuid",
		"filename":            "val_filename",
		"filepath":            "val_filepath",
		"filenameextension":   "val_filenameextension",
		"filemd5":             "val_filemd5",
		"filesha1":            "val_filesha1",
		"filesha256":          "val_filesha256",
		"filesha512":          "val_filesha512",
		"fileblake2b":         "val_fileblake2b",
		"filesffmt":           "val_filesffmt",
		"filesfmime":          "val_filesfmime",
		"filesfformatname":    "val_filesfformatname",
		"filesfformatversion": "val_filesfformatversion",
		"filesfidentnote":     "val_filesfidentnote",
		"filesfidentproof":    "val_filesfidentproof",
		"filectime":           "val_filectime",
		"filemtime":           "val_filemtime",
		"fileatime":           "val_fileatime",
		"filebtime":           "val_filebtime",
		"filensrl":            "val_filensrl",
	}

	_, err = prepin.Exec(
		want["fileuuid"], want["sessionuuid"], want["filename"], want["filepath"], want["filenameextension"],
		int64(12345), want["filemd5"], want["filesha1"], want["filesha256"], want["filesha512"], want["fileblake2b"],
		want["filesffmt"], want["filesfmime"], want["filesfformatname"], want["filesfformatversion"],
		want["filesfidentnote"], want["filesfidentproof"], want["filectime"], want["filemtime"], want["fileatime"],
		want["filebtime"], want["filensrl"], 3.14, 7,
	)
	if err != nil {
		t.Fatalf("PrepInsertFile exec error = %v", err)
	}

	for col, wantVal := range want {
		var got string
		if err := db.QueryRow("SELECT "+col+" FROM files WHERE fileuuid = ?", want["fileuuid"]).Scan(&got); err != nil {
			t.Fatalf("reading column %q back: %v", col, err)
		}
		if got != wantVal {
			t.Errorf("column %q = %q, want %q (columns are likely misaligned in the INSERT statement)", col, got, wantVal)
		}
	}

	var gotSize int64
	if err := db.QueryRow("SELECT filesize FROM files WHERE fileuuid = ?", want["fileuuid"]).Scan(&gotSize); err != nil {
		t.Fatalf("reading filesize back: %v", err)
	}
	if gotSize != 12345 {
		t.Errorf("filesize = %d, want 12345", gotSize)
	}

	var gotEntropy float64
	if err := db.QueryRow("SELECT fileentropy FROM files WHERE fileuuid = ?", want["fileuuid"]).Scan(&gotEntropy); err != nil {
		t.Fatalf("reading fileentropy back: %v", err)
	}
	if gotEntropy != 3.14 {
		t.Errorf("fileentropy = %v, want 3.14", gotEntropy)
	}

	var gotHierarchy int
	if err := db.QueryRow("SELECT hierarchy FROM files WHERE fileuuid = ?", want["fileuuid"]).Scan(&gotHierarchy); err != nil {
		t.Fatalf("reading hierarchy back: %v", err)
	}
	if gotHierarchy != 7 {
		t.Errorf("hierarchy = %d, want 7", gotHierarchy)
	}
}

// FuzzCreateFileTroveDB fuzzes three input string values
func FuzzCreateFileTroveDB(f *testing.F) {
	f.Fuzz(func(t *testing.T, a string, b string, c string) {
		CreateFileTroveDB(a, b, c)
		err := os.Remove(a)
		if err != nil {
			println("Could not remove test file")
		}

	})
}
