package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testUUID = "926be141-ab75-4106-8236-34edfcf102f2"

func tsvNames(tables ...string) []string {
	names := make([]string, len(tables))
	for i, table := range tables {
		names[i] = tsvFileName(testUUID, table)
	}
	return names
}

func TestTSVFileName(t *testing.T) {
	if got, want := tsvFileName(testUUID, "files"), testUUID+"_files.tsv"; got != want {
		t.Errorf("tsvFileName() = %q, want %q", got, want)
	}
}

func TestReportTSVExport(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name   string
		tables []string
		header string
	}{
		{"without optional features", []string{"session", "files", "directories"},
			"Exported session " + testUUID + " to 3 TSV files:\n"},
		{"all features", []string{"session", "files", "directories", "exif", "dublincore", "yara", "xattr"},
			"Exported session " + testUUID + " to 7 TSV files:\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			reportTSVExport(&stdout, &stderr, testUUID, tsvNames(tt.tables...))

			var want strings.Builder
			for _, table := range tt.tables {
				want.WriteString(filepath.Join(wd, testUUID+"_"+table+".tsv") + "\n")
			}
			if stdout.String() != want.String() {
				t.Errorf("stdout = %q, want %q", stdout.String(), want.String())
			}
			if stderr.String() != tt.header {
				t.Errorf("stderr = %q, want %q", stderr.String(), tt.header)
			}
		})
	}
}

func TestReportTSVExportFailure(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	exportErr := errors.New("permission denied")

	tests := []struct {
		name         string
		failed       string
		written      []string
		wantContains []string
		wantMissing  []string
	}{
		{"failure after two files", "directories", []string{"session", "files"},
			[]string{
				"TSV export of session " + testUUID + " is incomplete",
				"Writing " + filepath.Join(wd, testUUID+"_directories.tsv") + " failed",
				"permission denied",
				"Files written before the failure:",
				"  " + filepath.Join(wd, testUUID+"_session.tsv") + "\n",
				"  " + filepath.Join(wd, testUUID+"_files.tsv") + "\n",
			},
			[]string{"No TSV file was written"}},
		{"failure before the first file", "session", nil,
			[]string{
				"TSV export of session " + testUUID + " is incomplete",
				"Writing " + filepath.Join(wd, testUUID+"_session.tsv") + " failed",
				"permission denied",
				"No TSV file was written.",
			},
			[]string{"Files written before the failure"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			reportTSVExportFailure(&stderr, testUUID, tsvFileName(testUUID, tt.failed), tsvNames(tt.written...), exportErr)

			for _, want := range tt.wantContains {
				if !strings.Contains(stderr.String(), want) {
					t.Errorf("stderr = %q, want it to contain %q", stderr.String(), want)
				}
			}
			for _, unwanted := range tt.wantMissing {
				if strings.Contains(stderr.String(), unwanted) {
					t.Errorf("stderr = %q, want it not to contain %q", stderr.String(), unwanted)
				}
			}
		})
	}
}
