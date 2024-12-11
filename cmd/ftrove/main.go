package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	yara_x "github.com/VirusTotal/yara-x/go"

	"github.com/richardlehane/siegfried"
	"github.com/schollz/progressbar/v3"
	flag "github.com/spf13/pflag"

	ft "github.com/steffenfritz/FileTrove"
)

// version holds the version of FileTrove. Due to different build systems and GH Actions set manually for now.
var Version string = "v1.0.0-BETA-3"

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
	grepYARA := flag.StringP("yararules", "y", "", "Provide a YARA rule file and scan all files for matches.")
	dublincore := flag.StringP("dublincore", "d", "", "Add DublinCore metadata as a JSON file for a session (not single files).")
	exifData := flag.BoolP("exifdata", "e", false, "Get some EXIF metadata from image files.")
	exportSessionToTSV := flag.StringP("export-tsv", "t", "", "Export a session from the database to a TSV file. Provide the session uuid.")
	inDir := flag.StringP("indir", "i", "", "Input directory to work on.")
	install := flag.StringP("install", "", "", "Install FileTrove into the given, existing directory.")
	listSessions := flag.BoolP("list-sessions", "l", false, "List session information for all scans. Useful for exports.")
	listSession := flag.StringP("list-session", "L", "", "List information about a single session.")
	projectname := flag.StringP("project", "p", "", "A name for the project or scan session.")
	resumeuuid := flag.StringP("resume", "r", "", "Resume an aborted session. Provide the session uuid.")
	timezone := flag.StringP("timezone", "z", "", "Set the time zone to a region in which the timestamps of files are to be translated. If this flag is not set, the local time zone is used. Example: Europe/Berlin")
	debug := flag.BoolP("debug", "D", false, "Enable debug mode. This creates the diagnostic file debug_ftrove.")

	// updateFT := flag.BoolP("update-all", "u", false, "Update FileTrove, siegfried and NSRL.")
	printversion := flag.BoolP("version", "v", false, "Show version and build.")
	verbose := flag.BoolP("verbose", "V", false, "Print messages also to the terminal (stdout).")

	flag.Parse()

	starttime := time.Now()

	var fddebug os.File
	var err error

	if *debug {
		fddebug, err = ft.DebugCreateDebugPackage()
		if err != nil {
			logger.Error("Could not create debug package:", slog.String("error", err.Error()))
			os.Exit(1)
		}
		defer fddebug.Close()

		err = ft.DebugHostinformation(fddebug)
		if err != nil {
			logger.Error("Could not create debug hostinformation:", slog.String("error", err.Error()))
		}
		err = ft.DebugCheckInstalled(fddebug)
		if err != nil {
			logger.Error("Could not check installed version:", slog.String("error", err.Error()))
		}
		err = ft.DebugWriteFlags(fddebug, flag.Args())
		if err != nil {
			logger.Error("Could not write flags:", slog.String("error", err.Error()))
		}

	}

	// Init new session with flags
	var sessionmd ft.SessionMD
	sessionmd.Starttime = starttime.Format(time.RFC3339)
	sessionmd.UUID, _ = ft.CreateUUID()
	sessionmd.Archivistname = *archivistname
	sessionmd.Project = *projectname
	sessionmd.Mountpoint, _ = filepath.Abs(*inDir)
	sessionmd.Pathseparator = string(os.PathSeparator)
	sessionmd.Goversion = runtime.Version()
	sessionmd.Filetroveversion = Version
	sessionmd.Filetrovedbversion = Version // this might be redundant due to the fact that the db aligns with FT's version
	sessionmd.Sfversion = ft.SiegfriedVersion

	if *exifData {
		sessionmd.ExifFlag = "True"
	}
	if len(*dublincore) > 0 {
		sessionmd.Dublincoreflag = "True"
	}
	if len(*grepYARA) > 0 {
		sessionmd.Yaraflag = "True"
		sessionmd.Yarasource = *grepYARA
	}

	// Print banner or version on startup
	ft.PrintBanner()

	if *printversion {
		//ft.PrintLicense(Version, Build)
		ft.PrintLicense(Version)
		return
	}

	// Start installation
	if len(*install) > 0 {
		if strings.HasSuffix(*install, "/") {
			*install = strings.TrimRight(*install, "/")
		}
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
			if strings.HasPrefix(nsrlerr.Error(), "Could not download siegfried") {
				logger.Info("Could not download siegfried signature file. You have to copy siegfried.sig into the db directory. See the documentation.")
			}
			os.Exit(1)
		}
		if nsrlerr != nil {
			logger.Error("Could not download or create NSRL database.", slog.String("error", nsrlerr.Error()))
			if strings.HasPrefix(nsrlerr.Error(), "Could not download NSRL") {
				logger.Info("Could not download NSRL database. You have to copy a nsrl.db into the db directory. See the documentation.")
			}
			os.Exit(1)
		}
		// We put an an extra newline here due to the mixed output from install function and the logging here
		fmt.Println()
		logger.Info("Created all necessary files and directories successfully.")
		fmt.Println()

		// we return after the installation and quit the program with exit code 0
		return
	}

	// Connect to FileTrove's database. We don't do this with the other ready checks because of the export usecase
	// without a full install
	ftdb, err := ft.ConnectFileTroveDB("db")
	if err != nil {
		logger.Error("Could not connect to FileTrove's database. Quitting.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// list all started sessions and quit
	if *listSessions {
		fmt.Println("SESSION OVERVIEW\n")
		err := ft.ListSessions(ftdb)
		if err != nil {
			logger.Error("Could not query last sessions.", slog.String("error", err.Error()))
		}
		return
	}

	// Print info about a single session
	if len(*listSession) > 0 {
		fmt.Println("SESSION INFORMATION\n")
		smd, err := ft.ListSession(ftdb, *listSession)
		if err != nil {
			logger.Error("Could not query single session.", slog.String("error", err.Error()))
		}

		fmt.Println("Session UUID:\t" + smd.Sessionmd.UUID)
		fmt.Println("Project:\t" + smd.Sessionmd.Project)
		fmt.Println("Archivist:\t" + smd.Sessionmd.Archivistname)
		fmt.Println("Mountpoint:\t" + smd.Sessionmd.Mountpoint)
		fmt.Println("File Count:\t" + strconv.Itoa(smd.Filecount))
		fmt.Println("NSRL Count:\t" + strconv.Itoa(smd.Nsrlcount))

		return
	}

	// Check if ready for run.
	err = ft.CheckInstall(Version)
	if err != nil {
		logger.Error("FileTrove is not ready. Please see previous output.")
		os.Exit(-1)
	}

	// Change logger to MultiWriter for output to logfile and os.Stdout
	logfd, err := os.OpenFile(filepath.Join("logs", "filetrove.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
	if len(*resumeuuid) > 0 {
		logger.Info("Resuming session " + *resumeuuid)
	}

	/*if *updateFT {
		// check local versions against web page/online resource
	}*/

	// Export a specific session to TSV files
	if len(*exportSessionToTSV) != 0 {
		logger.Info("Export session " + *exportSessionToTSV + " to TSV files of the same name.")
		sessionValues, err := ft.ExportSessionSessionTSV(*exportSessionToTSV)
		if err != nil {
			logger.Error("Error while exporting files from session to TSV file.", slog.String("error", err.Error()))
			os.Exit(1)
		}
		// DOC: Value 7 MUST be the flag result of EXIF. We translate for clarity.
		exifFlagSet := sessionValues[7]
		// DOC: Value 8 MUST be the flag result of DublinCore. We translate for clarity.
		dcFlagSet := sessionValues[8]
		// DOC: Value 9 MUST be the flag result of YARA. We translate for clarity.
		yaraFlagSet := sessionValues[9]

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

		if dcFlagSet == "True" {
			err = ft.ExportSessionDCTSV(*exportSessionToTSV)
			if err != nil {
				logger.Error("Error while exporting DublinCore metadata from session to TSV file.", slog.String("error", err.Error()))
				os.Exit(1)
			}

		}

		if yaraFlagSet == "True" {
			err = ft.ExportYaraTSV(*exportSessionToTSV)
			if err != nil {
				logger.Error("Error while exporting YARA identified files from session to TSV file.", slog.String("error", err.Error()))
				os.Exit(1)
			}

		}

		logger.Info("Export successful.")
		return
	}

	// Init type for resuming a session
	var ri ft.ResumeInfo
	// Set up the file counter
	filecount := 0
	// Set up the counter for files that are in NSRL. This is just relevant for the short summary and log file entry.
	nsrlcount := 0

	// Prepare BoltDB for reading hashes
	db, err := ft.ConnectNSRL(filepath.Join("db", "nsrl.db"))
	if err != nil {
		logger.Error("Could not connect to NSRL database", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// If we resume a session, the following steps will NOT be executed as they are used for new sessions
	if len(*resumeuuid) == 0 {
		// Add new session to database
		nsrlversion, err := ft.GetNSRLVersion(db)
		if err != nil {
			logger.Error("Could not get NSRL version from NSRL database.", slog.String("error", err.Error()))
		}
		sessionmd.Nsrlversion = nsrlversion

		err = ft.InsertSession(ftdb, sessionmd)
		if err != nil {
			logger.Error("Could not add session to FileTrove database.", slog.String("error", err.Error()))
			err = ftdb.Close()
			if err != nil {
				logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
			}
			os.Exit(1)
		}

		// Add DublinCore metadata from json file to the database.
		// The metadata is meant for a whole session not for single files, i.e. the resource is the "mountpoint"
		if len(*dublincore) > 0 {
			dc, err := ft.ReadDC(*dublincore)
			if err != nil {
				logger.Error("Could not read DublinCore JSON file.", slog.String("error", err.Error()))
				os.Exit(1)
			}
			dcuuid, _ := ft.CreateUUID()
			err = ft.InsertDC(ftdb, sessionmd.UUID, dcuuid, dc)
			if err != nil {
				logger.Error("Could not add DublinCore to FileTrove database.", slog.String("error", err.Error()))
				os.Exit(1)
			}
		}
	} else {
		// for session resuming: read files already processed. This list ist compared to already
		// processed files. The diff updates the input file list.
		// We also fetch information like processed files from the session that was cancelled.
		ri, err = ft.ResumeLatestEntry(ftdb, *resumeuuid)
		if err != nil {
			logger.Error("Could not get session information for resuming.", slog.String("error", err.Error()))
			err = ftdb.Close()
			if err != nil {
				logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
			}
			os.Exit(1)
		}

		// Set the flag inDir to the mountpoint we read from the database
		*inDir = ri.Mountpoint

		// Overwrite the new session uuid with the resumed session's uuid
		sessionmd.UUID = *resumeuuid

		// Overwrite filecount with already processed files
		filecount = ri.ProcessedFiles

		// Overwrite the NSRL counter
		nsrlcount = ri.NSRLFiles

		// ToDo: Get Yara information for resuming sessions
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
	if *debug {
		ft.DebugWriteFileList(fddebug, filelist, dirlist)
	}
	if err != nil {
		logger.Error("An error occurred during the creation of the file list.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// If we resumeuuid, we update the file list to start with the file after the last entry
	if len(*resumeuuid) > 0 {
		fileIndex := slices.Index(filelist, ri.LastFile)
		if fileIndex == -1 {
			logger.Error("Could not find the last indexed file in the new file list.", slog.String("error", err.Error()))
			os.Exit(1)
		}

		filelist = filelist[fileIndex+1:]
		if len(filelist) == 0 {
			logger.Info("Input list is empty, no files are left to process. Quitting.", slog.String("info", "Input file of resumed session is empty."))
			err = ftdb.Close()
			if err != nil {
				logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
			}
			return
		}
	}

	// Initialize siegfried database
	s, err := siegfried.Load(filepath.Join("db", "siegfried.sig"))
	if err != nil {
		logger.Error("Could not read siegfried's database.", slog.String("error", err.Error()))
		err = ftdb.Close()
		if err != nil {
			logger.Error("Could not close database connection to FileTrove.", slog.String("error", err.Error()))
		}
		os.Exit(1)
	}

	// Compile YARA rule if flag provided, used later
	var checkYara bool
	var yaraRules *yara_x.Rules
	prepYaraInsert, err := ft.PrepInsertYara(ftdb)
	if err != nil {
		logger.Error("Could not prepare an insert statement for YARA inserts.", slog.String("error", err.Error()))
	}

	if len(*grepYARA) > 0 {
		checkYara = true
		yaraRules, err = ft.YaraCompile(*grepYARA)
		if err != nil {
			logger.Error("Could not compile the YARA rules.", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}

	// Inspect every file in filelist
	// Set up the progress bar
	bar := progressbar.Default(int64(len(filelist)))
	for _, file := range filelist {
		var filemd ft.FileMD

		//filemd.Filename = file

		// Create hash sums for every file
		hashsumsfile := make(map[string][]byte)

		// Calculate all supported hash sums for each file
		supportedHashes := ft.ReturnSupportedHashes()

		// Mutex lock on write hashes to map and
		// create waitgroup to prevent outside of the
		// clojure access, i.e. the NSRL check
		var mutex = &sync.Mutex{}
		var wg sync.WaitGroup

		for _, hash := range supportedHashes {
			wg.Add(1)

			go func(hash string) {
				defer wg.Done()

				hashsum, err := ft.Hashit(file, hash)
				if err != nil {
					logger.Error("Could not hash file.", slog.String("error", err.Error()))
				}

				mutex.Lock()
				hashsumsfile[hash] = hashsum
				mutex.Unlock()
			}(hash)
		}

		wg.Wait()

		filemd.Filename = filepath.Base(file)
		filemd.Filepath = file
		// This is a workaround of the not-so-perfect handling of golang's filepath.Ext() function
		//see https://github.com/golang/go/issues/66814
		if filepath.Ext(file) != filepath.Base(file) {
			filemd.Filenameextension = filepath.Ext(file)
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

		// Use specific timezone for the translation of timestamps
		if len(*timezone) != 0 {
			// Check if timezone string is valid
			timeIn, err := time.LoadLocation(*timezone)
			if err != nil {
				logger.Error("Timezone string is not valid. Example: Europe/Berlin", slog.String("error", err.Error()))
			}
			filemd.Fileatime = filetime.Atime.In(timeIn).String()
			filemd.Filectime = filetime.Ctime.In(timeIn).String()
			filemd.Filemtime = filetime.Mtime.In(timeIn).String()
		} else {
			filemd.Fileatime = filetime.Atime.String()
			filemd.Filectime = filetime.Ctime.String()
			filemd.Filemtime = filetime.Mtime.String()
		}

		// Check if the hash sum of the file is in the NSRL.
		// We use the db connection created by ft.ConnectNSRL()
		fileIsInNSRL, err := ft.GetValueNSRL(db, []byte(hex.EncodeToString(hashsumsfile["sha1"])))
		if err != nil {
			logger.Warn("Could not get value from NSRL database.", slog.String("warn", err.Error()))
		}

		if fileIsInNSRL {
			filemd.Filensrl = "TRUE"
			nsrlcount += 1
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

		filehierarchy := strings.Count(file, string(os.PathSeparator))

		_, err = prepInsertFile.Exec(fileuuid, sessionmd.UUID, filemd.Filename, filemd.Filepath, filemd.Filenameextension,
			filemd.Filesize, filemd.Filemd5, filemd.Filesha1, filemd.Filesha256, filemd.Filesha512, filemd.Fileblake2b,
			filemd.Filesffmt, filemd.Filesfmime, filemd.Filesfformatname, filemd.Filesfformatversion,
			filemd.Filesfidentnote, filemd.Filesfidentproof, filemd.Filectime, filemd.Filemtime, filemd.Fileatime,
			filemd.Filensrl, filemd.Fileentropy, filehierarchy)

		if err != nil {
			logger.Warn("Could not add file entry into FileTrove database. File: "+file, slog.String("warn", err.Error()))
		}

		// Check for YARA rules
		if checkYara {
			matches, err := ft.YaraScan(yaraRules, file)
			if err != nil {
				logger.Warn("Could not scan YARA rules", slog.String("error", err.Error()))
			}

			// As matches could be nil we check for it here
			if matches != nil {
				if len(matches.MatchingRules()) > 0 {
					for _, match := range matches.MatchingRules() {
						yarauuid, err := ft.CreateUUID()
						if err != nil {
							logger.Warn("Could not create UUID for YARA rule", slog.String("error", err.Error()))
						}
						_, err = prepYaraInsert.Exec(yarauuid, sessionmd.UUID, fileuuid, match.Identifier())
						if err != nil {
							logger.Warn("Could not add YARA identification", slog.String("error", err.Error()))
						}
					}
				}
			}

		}

		filecount += 1
		bar.Add(1)
	}

	// Add directory list to database. The metadata for directories is very limited so far.
	prepInsertDir, err := ft.PrepInsertDir(ftdb)
	if err != nil {
		logger.Error("Could not prepare an insert statement for directory inserts.", slog.String("error", err.Error()))
	}
	for _, direntry := range dirlist {
		var dirmd ft.DirMD

		diruuid, err := ft.CreateUUID()
		if err != nil {
			logger.Error("Could not create UUID for directory "+direntry, slog.String("error", err.Error()))
		}

		dirtime, err := ft.GetFileTimes(direntry)
		if err != nil {
			logger.Error("Could not get timestamps of directory "+direntry, slog.String("error", err.Error()))
		}

		// Use specific timezone for the translation of timestamps
		if len(*timezone) != 0 {
			// Check if timezone string is valid
			timeIn, err := time.LoadLocation(*timezone)
			if err != nil {
				logger.Error("Timezone string is not valid. Example: Europe/Berlin", slog.String("error", err.Error()))
			}
			dirmd.Diratime = dirtime.Atime.In(timeIn).String()
			dirmd.Dirctime = dirtime.Ctime.In(timeIn).String()
			dirmd.Dirmtime = dirtime.Mtime.In(timeIn).String()
		} else {
			dirmd.Diratime = dirtime.Atime.String()
			dirmd.Dirctime = dirtime.Ctime.String()
			dirmd.Dirmtime = dirtime.Mtime.String()
		}

		dirnamelist := strings.Split(direntry, string(os.PathSeparator))
		dirname := dirnamelist[len(dirnamelist)-1]

		dirhierarchy := strings.Count(direntry, string(os.PathSeparator))

		_, err = prepInsertDir.Exec(diruuid, sessionmd.UUID, dirname, direntry,
			dirtime.Ctime.String(),
			dirtime.Mtime.String(),
			dirtime.Atime.String(),
			dirhierarchy)

		if err != nil {
			logger.Warn("Could not add directory entry to FileTrove database.", slog.String("warn", err.Error()))
		}
	}

	// EXIF data for jpeg and tiff
	if *exifData {
		imagelist, err := ft.GetImageFiles(ftdb, sessionmd.UUID)
		if err != nil {
			logger.Error("Could not get image list from database.", slog.String("error", err.Error()))
		}
		for fileuuid, imagepath := range imagelist {
			// debug
			println(imagepath)

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
	logger.Info("Number of indexed files: " + strconv.Itoa(filecount))
	logger.Info("Number of indexed directory names: " + strconv.Itoa(len(dirlist)))
	logger.Info("Number of known files (NSRL=True): " + strconv.Itoa(nsrlcount))

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
