package main

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
)

// checkUnexpectedArgs returns an error if positional arguments remain after flag parsing.
// ftrove accepts no positional arguments. They typically come from a shell-expanded glob
// such as `-i dir/*`, where -i only receives the first match and the rest would be ignored.
// inDir is the value of the -i flag and is used to suggest the directory to scan instead.
func checkUnexpectedArgs(args []string, inDir string) error {
	if len(args) == 0 {
		return nil
	}

	quoted := make([]string, len(args))
	for i, arg := range args {
		quoted[i] = strconv.Quote(arg)
	}

	var msg strings.Builder
	msg.WriteString("unexpected arguments: " + strings.Join(quoted, " ") + "\n")

	if len(inDir) > 0 {
		parent := filepath.Dir(filepath.Clean(inDir))
		msg.WriteString("-i expects exactly one directory. If these paths come from a shell glob, scan their directory instead:\n")
		msg.WriteString("  ftrove -i " + shellQuote(parent) + "\n")
		msg.WriteString("The directory is scanned recursively, including hidden files, which a shell glob does not match.")
	} else {
		msg.WriteString("ftrove accepts no positional arguments. Pass the directory to scan with -i:\n")
		msg.WriteString("  ftrove -i <directory>\n")
		msg.WriteString("Do not use a shell glob such as dir/*: it does not match hidden files.")
	}

	return errors.New(msg.String())
}

// shellQuote wraps path in double quotes if it contains characters a shell would split or interpret.
func shellQuote(path string) string {
	if strings.ContainsAny(path, " \t\"'$`\\*?[]") {
		return strconv.Quote(path)
	}
	return path
}
