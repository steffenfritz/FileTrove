package filetrove

import (
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
