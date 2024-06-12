package main

import (
	"log/slog"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	ft "github.com/steffenfritz/FileTrove"
)

// Version holds the version of FileTrove and is set by the build system
var Version string = "v1.0.0-DEV-16"

// Build is not used anymore since DEV-11
// Build holds the sha1 fingerprint of the build and is set by the build system
// var Build string

// logger is the structured logger that is used for all logging levels
var logger *slog.Logger

func main() {
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Format of the source file MUST be a SHA1 hash per line
	createNSRL := flag.String("creatensrl", "", "Create or update a BoltDB file from a text file. A source file MUST be provided.")
	nsrlversion := flag.String("nsrlversion", "", "NSRL version flag. This string will be used for ftrove's session information.")
	updateDB := flag.String("updatedb", "", "Update a filetrove sqlite database to the next version. Expects the directory of the database file.")

	flag.Parse()

	if len(*createNSRL) != 0 {
		err := ft.CreateNSRLBoltDB(*createNSRL, *nsrlversion, "nsrl.db")
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
			logger.Info("You are at a very early version. No update possible, sorry.")
			return
		}
		logger.Info("You are at version " + instversion + ". Checking for next possible version.")

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

		// Update version 1.0.0-DEV-8 --> 1.0.0-DEV-9
		if instversion == "1.0.0-DEV-8" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-9' where version = '1.0.0-DEV-8'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-9.")
			return
		}

		// Update version 1.0.0-DEV-9 --> 1.0.0-DEV-10
		if instversion == "1.0.0-DEV-9" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-10' where version = '1.0.0-DEV-9'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-10.")
			return
		}

		// Update version 1.0.0-DEV-10 --> 1.0.0-DEV-11
		if instversion == "1.0.0-DEV-10" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-11' where version = '1.0.0-DEV-10'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-11.")
			return
		}

		// Update version 1.0.0-DEV-11 --> 1.0.0-DEV-12
		if instversion == "1.0.0-DEV-11" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-12' where version = '1.0.0-DEV-11'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-12.")
			return
		}

		// Update version 1.0.0-DEV-12 --> 1.0.0-DEV-13
		if instversion == "1.0.0-DEV-12" {
			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-13' where version = '1.0.0-DEV-12'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-13.")
			return
		}

		// Update version 1.0.0-DEV-13 --> 1.0.0-DEV-14
		if instversion == "1.0.0-DEV-13" {
			_, err = ftdb.Exec("ALTER TABLE files ADD hierarchy INTEGER")
			if err != nil {
				logger.Error("Could not update database", slog.String("error", err.Error()))
				return
			}
			_, err = ftdb.Exec("ALTER TABLE directories ADD hierarchy INTEGER")
			if err != nil {
				logger.Error("Could not update database", slog.String("error", err.Error()))
				return
			}

			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-14' where version = '1.0.0-DEV-13'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-14.")
			return
		}

		// Update version 1.0.0-DEV-14 --> 1.0.0-DEV-15
		if instversion == "1.0.0-DEV-14" {
			_, err = ftdb.Exec("ALTER TABLE directories RENAME TO directories_temp")
			if err != nil {
				logger.Error("Could not create temporary database for migration", slog.String("error", err.Error()))
				return
			}

			_, err = ftdb.Exec("CREATE TABLE directories(diruuid TEXT, sessionuuid TEXT, dirname TEXT, dircttime TEXT, dirmtime TEXT, diratime TEXT, hierarchy INTEGER)")
			if err != nil {
				logger.Error("Could not create new directories table for migration", slog.String("error", err.Error()))
				return
			}

			_, err = ftdb.Exec("INSERT INTO directories(diruuid, sessionuuid, dirname, hierarchy) SELECT diruuid, sessionuuid, dirname, hierarchy FROM directories_temp;")
			if err != nil {
				logger.Error("Could not copy old directories table to new one for migration", slog.String("error", err.Error()))
				return
			}

			_, err = ftdb.Exec("DROP TABLE directories_temp")
			if err != nil {
				logger.Error("Could not delete directories_temp table after migration", slog.String("error", err.Error()))
				return
			}

			_, err = ftdb.Exec("UPDATE filetrove SET version = '1.0.0-DEV-15' where version = '1.0.0-DEV-14'")
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
			logger.Info("FileTrove database updated to version 1.0.0-DEV-15.")
			return
		}

		// Update version 1.0.0-DEV-15 --> 1.0.0-DEV-16
		if instversion == "1.0.0-DEV-15" {
			logger.Info("FileTrove database cannot be updated to version 1.0.0-DEV-16 from your version. See changelog.")
			return
		}

		// Update version 1.0.0-DEV-16 --> 1.0.0-BETA-1
		if instversion == "1.0.0-DEV-16" {
			logger.Info("FileTrove database cannot be updated to version 1.0.0-BETA-1 from your version. See changelog.")
			return
		}

	}
}
