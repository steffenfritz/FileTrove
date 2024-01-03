package main

import (
	"log/slog"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	ft "github.com/steffenfritz/FileTrove"
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
	createNSRL := flag.String("creatensrl", "", "Create or update a BoltDB file from a text file. A source file MUST be provided.")
	updateDB := flag.String("updatedb", "", "Update a filetrove sqlite database to the next version. Expects the directory of the database file.")

	flag.Parse()

	if len(*createNSRL) != 0 {
		err := ft.CreateNSRLBoltDB(*createNSRL, "nsrl.db")
		if err != nil {
			logger.Error("Could not create BoltDB from NSRL text file", slog.String("error", err.Error()))
		}
		return
	}

	if len(*updateDB) != 0 {
		var instversion string
		ftdb, err := ft.ConnectFileTroveDB(*updateDB)
		if err != nil {
			logger.Error("Could not connect to FileTrove database", slog.String("error", err.Error()))
		}

		resrow := ftdb.QueryRow("SELECT version FROM filetrove")
		err = resrow.Scan(&instversion)
		if err != nil {
			logger.Error("Could not read the version of the FileTrove database", slog.String("error", err.Error()))
		}

		if instversion == "1.0.0-DEV-6" {
			logger.Info("You are at the latest version. No update possible.")
			return
		}
		logger.Info("You are at version " + instversion + ". You can upgrade to the next version.")

		// Update version 1.0.0-DEV-5 --> 1.0.0-DEV-6
		if instversion == "1.0.0-DEV-5" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-6' where version = '1.0.0-DEV-5'")
			if err != nil {
				logger.Error("Could not update database", slog.String("error", err.Error()))
				return
			}
			updatetime := time.Now().Format(time.RFC3339)
			_, err = ftdb.Exec("UPDATE filetrove SET lastupdate = ?", updatetime)
			if err != nil {
				logger.Error("Could not update last update time.", slog.String("error", err.Error()))
				return
			}
			logger.Info("FileTrove database updated to version 1.0.0-DEV-6.")
			return
		}

		// Update version 1.0.0-DEV-6 --> 1.0.0-DEV-7
		if instversion == "1.0.0-DEV-6" {
			logger.Info("There is no update path from 1.0.0-DEV-6 to 1.0.0.-DEV-7. Please backup database and recreate with --install flag.")
			return
		}

		// Update version 1.0.0-DEV-7 --> 1.0.0-DEV-8
		if instversion == "1.0.0-DEV-7" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-8' where version = '1.0.0-DEV-7'")
			if err != nil {
				logger.Error("Could not update database", slog.String("error", err.Error()))
				return
			}
			updatetime := time.Now().Format(time.RFC3339)
			_, err = ftdb.Exec("UPDATE filetrove SET lastupdate = ?", updatetime)
			if err != nil {
				logger.Error("Could not update last update time.", slog.String("error", err.Error()))
				return
			}
			logger.Info("FileTrove database updated to version 1.0.0-DEV-8.")
			return
		}

	}
}
