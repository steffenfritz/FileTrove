package filetrove

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExportSessionJSONL(t *testing.T) {
	tmpDir := t.TempDir()
	dbDir := filepath.Join(tmpDir, "db")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := CreateFileTroveDB(dbDir, "TEST", "2026-01-01"); err != nil {
		t.Fatal(err)
	}

	// Temporarily change working directory so ExportSessionJSONL can find db/filetrove.db
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(origDir)

	db, err := ConnectFileTroveDB("db")
	if err != nil {
		t.Fatal(err)
	}

	sessionuuid := "test-session-uuid-1234"
	s := SessionMD{
		UUID:               sessionuuid,
		Starttime:          "2026-01-01T00:00:00Z",
		Project:            "testproject",
		Archivistname:      "tester",
		Mountpoint:         "/tmp/test",
		Pathseparator:      "/",
		Filetroveversion:   "TEST",
		Filetrovedbversion: "TEST",
		Goversion:          "go1.0",
	}
	if err := InsertSession(db, s); err != nil {
		t.Fatal(err)
	}

	prepFile, err := PrepInsertFile(db)
	if err != nil {
		t.Fatal(err)
	}
	_, err = prepFile.Exec(
		"file-uuid-1", sessionuuid,
		"test.txt", "/tmp/test/test.txt", ".txt",
		42, "md5hash", "sha1hash", "sha256hash", "sha512hash", "blake2bhash",
		"fmt/111", "text/plain", "Plain Text", "1.0",
		"", "", "2026-01-01", "2026-01-01", "2026-01-01",
		"FALSE", 0.5, 1,
	)
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	var buf bytes.Buffer
	if err := ExportSessionJSONL(sessionuuid, &buf); err != nil {
		t.Fatalf("ExportSessionJSONL() error = %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least 2 JSONL lines, got %d", len(lines))
	}

	typesSeen := make(map[string]bool)
	for _, line := range lines {
		var rec JSONLRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			t.Fatalf("failed to unmarshal line %q: %v", line, err)
		}
		typesSeen[rec.Type] = true
	}

	for _, want := range []string{"session", "file"} {
		if !typesSeen[want] {
			t.Errorf("expected type %q in JSONL output, but not found", want)
		}
	}
}
