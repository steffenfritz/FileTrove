package main

import (
	"encoding/hex"
	"fmt"
	"github.com/richardlehane/siegfried"
	flag "github.com/spf13/pflag"
	"log/slog"
	"os"
	"time"

	ft "github.com/steffenfritz/FileTrove"
)

// Version holds the version of FileTrove and is set by the build system
var Version string

// Build holds the sha1 fingerprint of the build and is set by the build system
var Build string

// tsStartedFormated is the formated timestamp when FileTrove was started
var tsStartedFormated string

// logger is the structured logger that is used for all logging levels
var logger *slog.Logger

func init() {
	tsStarted := time.Now()
	tsStartedFormated = tsStarted.Format("2006-01-02_15:04:05")
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func main() {
	// archivistname := flag.StringP("archivist", "a", "", "The name of the person responsible for the scan.")
	createNSRL := flag.String("creatensrl", "", "Create a BoltDB file from a text file. A source file MUST be provided.")
	// exportResultsToCSV
	inDir := flag.StringP("indir", "i", "", "Input directory to work on.")
	install := flag.StringP("install", "I", "", "Install FileTrove into the given directory.")
	// listSessions
	// projectname := flag.StringP("project", "p", "", "A name for the project or scan session.")

	updateFT := flag.BoolP("update-all", "u", false, "Update FileTrove, siegfried and NSRL.")
	version := flag.BoolP("version", "v", false, "Show version and build.")

	flag.Parse()

	starttime := time.Now()

	var sessionmd ft.SessionMD
	sessionmd.Starttime = starttime.Format(time.RFC3339)
	sessionmd.UUID, _ = ft.CreateUUID()

	if *version {
		ft.PrintLicense()
		fmt.Println("Version: " + Version + " Build: " + Build + "\n")
		return
	}

	if len(*createNSRL) != 0 {
		err := ft.CreateNSRLBoltDB(*createNSRL, "db/nsrl.db")
		if err != nil {
			logger.Error("Could not create BoltDB from NSRL text file", slog.String("error", err.Error()))
		}
		return
	}

	if len(*install) > 0 {
		direrr, logserr, trovedberr, siegfriederr, nsrlerr := ft.InstallFT(*install, Version, tsStartedFormated)
		if direrr != nil {
			logger.Error("Could not create db directory.", slog.String("error", direrr.Error()))
			os.Exit(1)
		}
		if logserr != nil {
			logger.Error("Could not create logs directory.", slog.String("error", direrr.Error()))
			os.Exit(1)
		}
		if trovedberr != nil {
			logger.Error("Could not create FileTrove database.", slog.String("error", trovedberr.Error()))
			os.Exit(1)
		}
		if siegfriederr != nil {
			logger.Error("Could not download or create siegfried database.", slog.String("error", siegfriederr.Error()))
			os.Exit(1)
		}
		if nsrlerr != nil {
			logger.Error("Could not download or create NSRL database.", slog.String("error", nsrlerr.Error()))
			os.Exit(1)
		}
		logger.Info("Created all necessary files and directories successfully.")

		// we return here and quit the program with exit code 0
		return
	}

	if *updateFT {
		// check local hashes against web page/online resource
	}
	// check if ready for run
	// if not suggest downloads and installation or exit

	// Connect to FileTrove's database
	ftdb, err := ft.ConnectFileTroveDB("db/filetrove.db")
	if err != nil {
		logger.Error("Could not connect to FileTrove's database.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Create file list
	// TODO: Add dirlist do output
	filelist, _, err := ft.CreateFileList(*inDir)
	if err != nil {
		logger.Error("An error occurred during the creation of the file list.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Initialize siegfried database
	s, err := siegfried.Load("db/siegfried.sig")
	if err != nil {
		logger.Error("Could not read siegfried's database.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Prepare BoltDB for reading hashes
	db, err := ft.ConnectNSRL("db/nsrl.db")
	if err != nil {
		logger.Error("Could not connect to NSRL database", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Inspect every file in our list
	for _, file := range filelist {
		var filemd ft.FileMD

		filemd.Filename = file

		// Create hash sums for every file
		hashsumsfile := make(map[string][]byte)

		// Calculate all supported hash sums for each file
		supportedHashes := ft.ReturnSupportedHashes()
		for _, hash := range supportedHashes {
			hashsum, err := ft.Hashit(file, hash)
			if err != nil {
				logger.Error("Could not hash file.", slog.String("error", err.Error()))
			}
			hashsumsfile[hash] = hashsum
		}

		// Add all hash sums to the filemd struct for writing into the file database
		filemd.Filemd5 = string(hashsumsfile["md5"])
		filemd.Filesha1 = string(hashsumsfile["sha1"])
		filemd.Filesha256 = string(hashsumsfile["sha256"])
		filemd.Filesha256 = string(hashsumsfile["sha512"])
		filemd.Fileblake2b = string(hashsumsfile["blake2b"])

		// Get siegfried information for each file. These are those in the type SiegfriedType
		oneFile, err := ft.SiegfriedIdent(s, file)
		if err != nil {
			logger.Error("Could not identify file using siegfried", slog.String("error", err.Error()))
		}
		filemd.Filesize = oneFile.SizeInByte
		filemd.Filesfmime = oneFile.MIMEType
		filemd.Filesfformatname = oneFile.FormatName
		filemd.Filesfformatversion = oneFile.FormatVersion
		filemd.Filesffmt = oneFile.FMT
		filemd.Filesfidentnote = oneFile.IdentificationNote
		filemd.Filesfidentproof = oneFile.IdentificationProof
		filemd.Filesfregistry = oneFile.Registry

		// Get file times
		filetime, err := ft.GetFileTimes(file)
		if err != nil {
			logger.Error("Could not get access, change or birth time for file.", slog.String("error", err.Error()))
		}

		filemd.Fileatime = filetime.Atime.String()
		filemd.Filectime = filetime.Ctime.String()
		filemd.Filemtime = filetime.Mtime.String()

		// Check if the hash sum of the file is in the NSRL.
		// We use the db connection created by ft.ConnectNSRL()
		fileIsInNSRL, err := ft.GetValueNSRL(db, []byte(hex.EncodeToString(hashsumsfile["sha1"])))

		if fileIsInNSRL {
			filemd.Filensrl = "TRUE"
		} else {
			filemd.Filensrl = "FALSE"
		}

		// write to DB
	}

	endtime := time.Now()
	sessionmd.Endtime = endtime.Format(time.RFC3339)
	// Close database connection and quit FileTrove
	err = ftdb.Close()
	if err != nil {
		logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
	}

}
