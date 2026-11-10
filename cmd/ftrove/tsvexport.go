package main

import (
	"fmt"
	"io"
	"path/filepath"
)

// tsvFileName returns the name of the TSV file that the ft.Export*TSV functions in db.go
// create for a session and table. It MUST match the names used there.
func tsvFileName(sessionuuid, table string) string {
	return sessionuuid + "_" + table + ".tsv"
}

// absPath returns the absolute path of name, or name itself if it cannot be resolved.
func absPath(name string) string {
	if abs, err := filepath.Abs(name); err == nil {
		return abs
	}
	return name
}

// reportTSVExport prints the absolute paths of the written TSV files to stdout, one per line,
// and a summary line to stderr. Stdout carries only the paths so that it can be piped.
func reportTSVExport(stdout, stderr io.Writer, sessionuuid string, files []string) {
	fmt.Fprintf(stderr, "Exported session %s to %d TSV files:\n", sessionuuid, len(files))
	for _, f := range files {
		fmt.Fprintln(stdout, absPath(f))
	}
}

// reportTSVExportFailure tells the user on stderr that the TSV export is incomplete, which
// file failed and which files were written before the failure.
func reportTSVExportFailure(stderr io.Writer, sessionuuid, failed string, files []string, err error) {
	fmt.Fprintf(stderr, "Error: TSV export of session %s is incomplete.\n", sessionuuid)
	fmt.Fprintf(stderr, "Writing %s failed (the file may be partially written): %v\n", absPath(failed), err)
	if len(files) == 0 {
		fmt.Fprintln(stderr, "No TSV file was written.")
		return
	}
	fmt.Fprintln(stderr, "Files written before the failure:")
	for _, f := range files {
		fmt.Fprintln(stderr, "  "+absPath(f))
	}
}
