package filetrove

import (
	"os"
)

// InstallFT creates and downloads necessary directories and databases and copies them to installPath
func InstallFT(installPath string, version string, initdate string) (error, error, error, error) {
	direrr := os.Mkdir(installPath+"/db", os.ModePerm)
	if direrr != nil {
		return direrr, nil, nil, nil
	}
	CreateFileTroveDB(installPath+"/db", version, initdate)
	siegfriederr := GetSiegfriedDB()
	GetNSRLDB()

	return direrr, nil, siegfriederr, nil
}
