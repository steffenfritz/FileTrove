package filetrove

import (
	"fmt"
	"os"
	"path/filepath"
)

// InstallFT creates necessary directories and databases
func InstallFT(installPath string, version string, initdate string) (error, error, error, error) {

	// Printing an additional newline
	fmt.Println()

	fmt.Println("Creating database and logfile directories.")
	dbdirerr := os.Mkdir(filepath.Join(installPath, "db"), os.ModePerm)
	if dbdirerr != nil {
		return dbdirerr, nil, nil, nil
	}
	logsdirerr := os.Mkdir(filepath.Join(installPath, "logs"), os.ModePerm)
	if logsdirerr != nil {
		return nil, logsdirerr, nil, nil
	}
	fmt.Println("Creating filetrove database.")
	trovedberr := CreateFileTroveDB(filepath.Join(installPath, "db"), version, initdate)
	if trovedberr != nil {
		return nil, nil, trovedberr, nil
	}
	fmt.Println("Downloading signature database.")
	siegfriederr := GetSiegfriedDB(installPath)

	fmt.Println("\nNSRL bloom filter must be placed in the db/ directory as nsrl.bloom.")
	fmt.Println("Build it with: task nsrl:build-modern (or nsrl:build-mobile, nsrl:build-all)")
	fmt.Println("Or copy an existing nsrl.bloom file into the db/ directory.")

	return dbdirerr, logsdirerr, trovedberr, siegfriederr
}

// CheckInstall checks if all necessary files are available
func CheckInstall(version string) error {
	_, err := os.Stat(filepath.Join("db", "siegfried.sig"))
	if os.IsNotExist(err) {
		fmt.Println("ERROR: siegfried signature file not installed.")
	}
	_, err = os.Stat(filepath.Join("db", "filetrove.db"))
	if os.IsNotExist(err) {
		fmt.Println("ERROR: filetrove database does not exist.")
	}
	_, dberr := os.Stat(filepath.Join("db", "nsrl.bloom"))
	if os.IsNotExist(dberr) {
		// Check for legacy nsrl.db and provide migration hint
		if _, legacyErr := os.Stat(filepath.Join("db", "nsrl.db")); legacyErr == nil {
			fmt.Println("ERROR: Legacy nsrl.db detected. Run 'task nsrl:build-modern' or rebuild with admftrove --creatensrl to create nsrl.bloom.")
		} else {
			fmt.Println("ERROR: nsrl bloom filter does not exist.")
		}
	}

	if dberr == nil {
		ftdb, connerr := ConnectFileTroveDB("db")
		if connerr != nil {
			fmt.Println("Could not connect or open database. Error: " + connerr.Error())
			os.Exit(1)
		}

		compatible, dbversion, checkerr := CheckVersion(ftdb, version)
		if checkerr != nil {
			fmt.Println("Could not check database version. Error: " + checkerr.Error())
		}
		if !compatible {
			fmt.Println("Database not compatible with this Version of FileTrove. Database version: " + dbversion)
			os.Exit(1)
		}
	}

	if err != nil {
		fmt.Println("ERROR: Some or more checks failed, FileTrove is not ready. Did you run the installation?")
		return err
	}

	return nil
}
