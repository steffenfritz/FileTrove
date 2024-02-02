package filetrove

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// InstallFT creates and downloads necessary directories and databases and copies them to installPath
func InstallFT(installPath string, version string, initdate string) (error, error, error, error, error) {
	var choice string

	dbdirerr := os.Mkdir(installPath+"/db", os.ModePerm)
	if dbdirerr != nil {
		return dbdirerr, nil, nil, nil, nil
	}
	logsdirerr := os.Mkdir(installPath+"/logs", os.ModePerm)
	if logsdirerr != nil {
		return nil, logsdirerr, nil, nil, nil
	}
	trovedberr := CreateFileTroveDB(installPath+"/db", version, initdate)
	if trovedberr != nil {
		return nil, nil, trovedberr, nil, nil
	}
	siegfriederr := GetSiegfriedDB(installPath)

	fmt.Print("\nNext step is to download the NSRL database which is 3.5GB. Proceed? [y/n]: ")
	_, err := fmt.Scan(&choice)
	if err != nil {
		os.Exit(-1)
	}

	choice = strings.TrimSpace(choice)
	choice = strings.ToLower(choice)

	var nsrlerr error
	if choice == "y" {
		nsrlerr = GetNSRL()

	} else {
		log.Println("Skipping NSRL download. You have to copy an existing nsrl.db into the db directory.")
	}

	return dbdirerr, logsdirerr, trovedberr, siegfriederr, nsrlerr
}

// CheckInstall checks if all necessary file are available
func CheckInstall(version string) error {
	_, err := os.Stat("db/siegfried.sig")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: siegfried signature file not installed.")
	}
	_, err = os.Stat("db/filetrove.db")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: filetrove database does not exist.")
	}
	_, dberr := os.Stat("db/nsrl.db")
	if os.IsNotExist(dberr) {
		fmt.Println("ERROR: nsrl database does not exist.")
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
