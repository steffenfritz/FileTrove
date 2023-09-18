package filetrove

import (
	"fmt"
	"os"
)

// InstallFT creates and downloads necessary directories and databases and copies them to installPath
func InstallFT(installPath string, version string, initdate string) (error, error, error, error, error) {
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
	siegfriederr := GetSiegfriedDB()
	GetNSRLDB()

	return dbdirerr, logsdirerr, trovedberr, siegfriederr, nil
}

// CheckInstall checks if all necessary file are available
func CheckInstall() error {
	_, err := os.Stat("db/siegfried.sig")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: siegfried signature file not installed.")
	}
	_, err = os.Stat("db/filetrove.db")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: filetrove database does not exist.")
	}
	_, err = os.Stat("db/nsrl.db")
	if os.IsNotExist(err) {
		fmt.Println("ERROR: nsrl database does not exist.")
	}

	if err != nil {
		fmt.Println("ERROR: Some or more checks failed, FileTrove is not ready. Did you run the installation?")
		return err
	}

	return nil
}
