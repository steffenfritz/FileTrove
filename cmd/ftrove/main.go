package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/richardlehane/siegfried"
	"github.com/schollz/progressbar/v3"
	flag "github.com/spf13/pflag"

	ft "github.com/steffenfritz/FileTrove"
)

// version holds the version of FileTrove and is set by the build system, i.e. goreleaser
var Version string

// commit holds the sha1 fingerprint of the build and is set by the build system, i.e. goreleaser
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
	archivistname := flag.StringP("archivist", "a", "", "The name of the person responsible for the scan.")
	// createThumbsImages :=
	// createStillsVideo :=
	// getTextfileIdea :=
	// grepYARA :=
	exifData := flag.BoolP("exifdata", "e", false, "Get some EXIF metadata from image files.")
	exportSessionToTSV := flag.StringP("export-tsv", "t", "", "Export a session from the database to a TSV file. Provide the session uuid.")
	inDir := flag.StringP("indir", "i", "", "Input directory to work on.")
	install := flag.StringP("install", "", "", "Install FileTrove into the given directory.")
	listSessions := flag.BoolP("list-sessions", "l", false, "List session information for all scans. Useful for exports.")
	projectname := flag.StringP("project", "p", "", "A name for the project or scan session.")

	// updateFT := flag.BoolP("update-all", "u", false, "Update FileTrove, siegfried and NSRL.")
	printversion := flag.BoolP("version", "v", false, "Show version and build.")
	verbose := flag.BoolP("verbose", "V", false, "Print messages also to the terminal (stdout).")

	flag.Parse()

	starttime := time.Now()

	var sessionmd ft.SessionMD
	sessionmd.Starttime = starttime.Format(time.RFC3339)
	sessionmd.UUID, _ = ft.CreateUUID()
	sessionmd.Archivistname = *archivistname
	sessionmd.Project = *projectname
	if *exifData {
		sessionmd.ExifFlag = "True"
	}

	ft.PrintBanner()

	if *printversion {
		ft.PrintLicense(Version, Build)
		return
	}

	if len(*install) > 0 {
		logger.Info("FileTrove installation started. Version: " + Version)
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

		// we return after the installation and quit the program with exit code 0
		return
	}

	// Check if ready for run.
	err := ft.CheckInstall(Version)
	if err != nil {
		logger.Error("FileTrove is not ready. Please check previous output.")
		os.Exit(-1)
	}

	// Change logger to MultiWriter for output to logfile and os.Stdout
	logfd, err := os.OpenFile("logs/filetrove.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("Could not open filetrove log file.", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if *verbose {
		logger.Info("Redirecting logs to stdout and logs/filetrove.log")
		logw := io.MultiWriter(os.Stdout, logfd)
		logger = slog.New(slog.NewTextHandler(logw, nil))
	} else {
		logger.Info("Redirecting logs to logs/filetrove.log")
		logw := io.Writer(logfd)
		logger = slog.New(slog.NewTextHandler(logw, nil))

	}
	logger.Info("FileTrove started at " + starttime.String())

	/*if *updateFT {
		// check local versions against web page/online resource
	}*/

	// Connect to FileTrove's database
	ftdb, err := ft.ConnectFileTroveDB("db")
	if err != nil {
		logger.Error("Could not connect to FileTrove's database. Quitting.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if *listSessions {
		err := ft.ListSessions(ftdb)
		if err != nil {
			logger.Error("Could not query last sessions.", slog.String("error", err.Error()))
		}
		return
	}

	if len(*exportSessionToTSV) != 0 {
		logger.Info("Export session " + *exportSessionToTSV + " to TSV files of the same name.")
		sessionValues, err := ft.ExportSessionSessionTSV(*exportSessionToTSV)
		if err != nil {
			logger.Error("Error while exporting files from session to TSV file.", slog.String("error", err.Error()))
			os.Exit(1)
		}
		// DOC: Value 6 MUST be the flag result of EXIF. We translate for clarity.
		exifFlagSet := sessionValues[6]

		err = ft.ExportSessionFilesTSV(*exportSessionToTSV)
		if err != nil {
			logger.Error("Error while exporting files from session to TSV file.", slog.String("error", err.Error()))
			os.Exit(1)
		}
		err = ft.ExportSessionDirectoriesTSV(*exportSessionToTSV)
		if err != nil {
			logger.Error("Error while exporting directories from session to TSV file.", slog.String("error", err.Error()))
			os.Exit(1)
		}

		if exifFlagSet == "True" {
			err = ft.ExportSessionEXIFTSV(*exportSessionToTSV)
			if err != nil {
				logger.Error("Error while exporting EXIF metadata from session to TSV file.", slog.String("error", err.Error()))
				os.Exit(1)
			}
		}

		logger.Info("Export successful.")
		return
	}

	// Add new session to database
	err = ft.InsertSession(ftdb, sessionmd)
	if err != nil {
		logger.Error("Could not add session to FileTrove database.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Prepare statement to add file scan results to database
	prepInsertFile, err := ft.PrepInsertFile(ftdb)
	if err != nil {
		logger.Error("Could not prepare an insert statement for file inserts.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Create file list
	filelist, dirlist, err := ft.CreateFileList(*inDir)
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

	// Inspect every file in filelist
	// Set up the progress bar
	bar := progressbar.Default(int64(len(filelist)))
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
		filemd.Filemd5 = hex.EncodeToString(hashsumsfile["md5"])
		filemd.Filesha1 = hex.EncodeToString(hashsumsfile["sha1"])
		filemd.Filesha256 = hex.EncodeToString(hashsumsfile["sha256"])
		filemd.Filesha512 = hex.EncodeToString(hashsumsfile["sha512"])
		filemd.Fileblake2b = hex.EncodeToString(hashsumsfile["blake2b-512"])

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
		if err != nil {
			logger.Warn("Could not get value from NSRL database.", slog.String("warn", err.Error()))
		}

		if fileIsInNSRL {
			filemd.Filensrl = "TRUE"
		} else {
			filemd.Filensrl = "FALSE"
		}

		// Calculate entropy of the file
		filemd.Fileentropy, err = ft.Entropy(file)
		if err != nil {
			logger.Warn("Could not calculate entropy for file "+file, slog.String("warning", err.Error()))
		}

		// Create a UUID for every file that is written to the database. This UUID is not stable over several runs!
		fileuuid, err := ft.CreateUUID()
		if err != nil {
			logger.Error("Could not create UUID for file "+file, slog.String("error", err.Error()))
			os.Exit(1)
		}
		_, err = prepInsertFile.Exec(fileuuid, sessionmd.UUID, filemd.Filename, filemd.Filesize,
			filemd.Filemd5, filemd.Filesha1, filemd.Filesha256, filemd.Filesha512, filemd.Fileblake2b,
			filemd.Filesffmt, filemd.Filesfmime, filemd.Filesfformatname, filemd.Filesfformatversion,
			filemd.Filesfidentnote, filemd.Filesfidentnote, filemd.Filectime, filemd.Filemtime, filemd.Fileatime,
			filemd.Filensrl, filemd.Fileentropy)

		if err != nil {
			logger.Warn("Could not add file entry into FileTrove database. File: "+file, slog.String("warn", err.Error()))
		}
		bar.Add(1)
	}

	// Add directory list to database. The metadata for directories is very limited so far.
	prepInsertDir, err := ft.PrepInsertDir(ftdb)
	if err != nil {
		logger.Error("Could not prepare an insert statement for directory inserts.", slog.String("error", err.Error()))
	}
	for _, direntry := range dirlist {
		diruuid, err := ft.CreateUUID()
		if err != nil {
			logger.Error("Could not create UUID for directory "+direntry, slog.String("error", err.Error()))
		}
		_, err = prepInsertDir.Exec(diruuid, sessionmd.UUID, direntry)

		if err != nil {
			logger.Warn("Could not add directory entry to FileTrove database.", slog.String("warn", err.Error()))
		}
	}

	// EXIF data for jpeg and tiff files
	if *exifData {
		imagelist, err := ft.GetImageFiles(ftdb, sessionmd.UUID)
		if err != nil {
			logger.Error("Could not get image list from database.", slog.String("error", err.Error()))
		}
		for fileuuid, imagepath := range imagelist {
			exifparsed, err := ft.ExifDecode(imagepath)
			if err != nil {
				logger.Error("Could not parse image for exif data. File: "+imagepath, slog.String("error", err.Error()))
			}

			exifuuid, err := ft.CreateUUID()
			if err != nil {
				logger.Error("Could not create UUID for exif entry.", slog.String("error", err.Error()))
			}
			err = ft.InsertExif(ftdb, exifuuid, sessionmd.UUID, fileuuid, exifparsed)
			if err != nil {
				logger.Error("Could not insert EXIF metadata into FileTrove database.", slog.String("error", err.Error()))
			}
		}
	}

	endtime := time.Now()

	// Short report after run, written to stdout if verbose and always to the log file
	fmt.Println()
	absPath, _ := filepath.Abs(*inDir)
	logger.Info("Finished indexing of " + absPath)
	logger.Info("Session UUID: " + sessionmd.UUID)
	logger.Info("Number of indexed files: " + strconv.Itoa(len(filelist)))
	logger.Info("Number of indexed directory names: " + strconv.Itoa(len(dirlist)))

	runtime := endtime.Sub(starttime)
	logger.Info("Indexing took: " + runtime.String())
	logger.Info("All results are written to the sqlite database db/filetrove.db")
	logger.Info("You can export the results with ./ftrove -t " + sessionmd.UUID)
	// End short report

	sessionmd.Endtime = endtime.Format(time.RFC3339)
	_, err = ftdb.Exec("UPDATE sessionsmd SET endtime=\"" + sessionmd.Endtime + "\" WHERE uuid=\"" + sessionmd.UUID + "\"RETURNING *;")
	if err != nil {
		logger.Error("Could not write endtime to database.", slog.String("error", err.Error()))
	}
	// Close database connection and quit FileTrove
	err = ftdb.Close()
	if err != nil {
		logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
	}
	err = logfd.Close()
	if err != nil {
		_ = fmt.Errorf("ERROR: Could not close error log file: " + err.Error())
	}
}
