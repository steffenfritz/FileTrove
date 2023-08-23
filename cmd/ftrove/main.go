package main

import (
	"github.com/richardlehane/siegfried"
	flag "github.com/spf13/pflag"
	"log/slog"
	"os"
	"time"

	ft "github.com/steffenfritz/FileTrove"
)

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
	install := flag.StringP("install", "I", "", "Install FileTrove into the given directory.")
	inDir := flag.StringP("indir", "i", "", "Input directory to work on.")
	updateFT := flag.BoolP("update-all", "u", false, "Update FileTrove, siegfried and NSRL.")

	flag.Parse()

	if len(*install) > 0 {
		direrr, trovedberr, siegfriederr, nsrlerr := ft.InstallFT(*install)
		if direrr != nil {
			logger.Error("Could no create db directory.", slog.String("message", direrr.Error()))
			os.Exit(1)
		}
		if trovedberr != nil {
			logger.Error("Could no create FileTrove database.", slog.String("message", trovedberr.Error()))
			os.Exit(1)
		}
		if siegfriederr != nil {
			logger.Error("Could not download or create siegfried database.", slog.String("message", siegfriederr.Error()))
			os.Exit(1)
		}
		if nsrlerr != nil {
			logger.Error("Could not download or create NSRL database.", slog.String("message", nsrlerr.Error()))
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
	// if not suggest downloads or exit

	// Create file list
	// TODO: Add dirlist do output
	filelist, _, err := ft.CreateFileList(*inDir)
	if err != nil {
		logger.Error("An error occurred during the creation of the file list.", slog.String("message", err.Error()))
		os.Exit(1)
	}

	// Initialize siegfried database
	s, err := siegfried.Load("db/siegfried.sig")
	if err != nil {
		logger.Error("Could not read siegfried's database.", slog.String("message", err.Error()))
		os.Exit(1)
	}

	// Inspect every file in our list
	for _, file := range filelist {
		// Calculate all supported hash sums for each file
		supportedHashes := ft.ReturnSupportedHashes()
		for _, hash := range supportedHashes {
			hashsum, err := ft.Hashit(file, hash)
			if err != nil {

			}
			// debug
			println(hash, ":", hashsum)

		}
		// Get siegfried information for each file
		oneFile, err := ft.SiegfriedIdent(s, file)
		if err != nil {
			logger.Error("Could not identify file using siegfried", slog.String("message", err.Error()))
		}
		// DEBUG
		println(oneFile.FileName)
		println(oneFile.SizeInByte)
		println(oneFile.MIMEType)

		// -- PRONOM ID
		// -- file size
		// --check if in NSRL
		// write to DB
	}

}
