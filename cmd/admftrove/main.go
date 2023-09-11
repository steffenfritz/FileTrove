package main

import (
	flag "github.com/spf13/pflag"
	ft "github.com/steffenfritz/FileTrove"
	"log/slog"
	"os"
)

// Version holds the version of ^FileTrove and is set by the build system
var Version string

// Build holds the sha1 fingerprint of the build and is set by the build system
var Build string

// logger is the structured logger that is used for all logging levels
var logger *slog.Logger

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Format of the source file MUST be a SHA1 hash per line
	createNSRL := flag.String("creatensrl", "", "Create a BoltDB file from a text file. A source file MUST be provided.")

	flag.Parse()

	if len(*createNSRL) != 0 {
		err := ft.CreateNSRLBoltDB(*createNSRL, "nsrl.db")
		if err != nil {
			logger.Error("Could not create BoltDB from NSRL text file", slog.String("error", err.Error()))
		}
		return
	}
}
