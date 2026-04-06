package filetrove

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// NSRLBloomURL is the download URL for the pre-built NSRL bloom filter.
// Update this constant when a new NSRL build is published to GitHub Releases.
const NSRLBloomURL = "https://github.com/steffenfritz/FileTrove/releases/download/nsrl-2026.03.1/nsrl.bloom"

// InstallFT creates necessary directories and databases
func InstallFT(installPath string, version string, initdate string) (error, error, error, error) {

	// Printing an additional newline
	fmt.Println()

	fmt.Println("Creating database and logfile directories.")
	dbdirerr := os.MkdirAll(filepath.Join(installPath, "db"), os.ModePerm)
	if dbdirerr != nil {
		return dbdirerr, nil, nil, nil
	}
	logsdirerr := os.MkdirAll(filepath.Join(installPath, "logs"), os.ModePerm)
	if logsdirerr != nil {
		return nil, logsdirerr, nil, nil
	}
	fmt.Println("Creating filetrove database.")
	trovedberr := CreateFileTroveDB(filepath.Join(installPath, "db"), version, initdate)
	if trovedberr != nil {
		return nil, nil, trovedberr, nil
	}
	var siegfriederr error
	sigPath := filepath.Join(installPath, "db", "siegfried.sig")
	if _, err := os.Stat(sigPath); err == nil {
		fmt.Println("Siegfried signature file already present.")
	} else {
		fmt.Println("Downloading signature database.")
		siegfriederr = GetSiegfriedDB(installPath)
	}

	// Try to find, copy, or download the NSRL bloom filter
	nsrlDst := filepath.Join(installPath, "db", "nsrl.bloom")
	if _, err := os.Stat(nsrlDst); os.IsNotExist(err) {
		if err := copyNSRLBloom(nsrlDst); err == nil {
			fmt.Println("Copied NSRL bloom filter to " + nsrlDst)
		} else {
			fmt.Println("Downloading NSRL bloom filter (~150 MB, this may take a while)...")
			if dlErr := DownloadNSRLBloom(nsrlDst); dlErr != nil {
				fmt.Println("\nNSRL bloom filter could not be downloaded: " + dlErr.Error())
				fmt.Println("Build it manually with: task nsrl:build-modern")
				fmt.Println("Or copy an existing nsrl.bloom file into the db/ directory.")
				fmt.Println("Scanning will work without it; NSRL checks will be skipped.")
			} else {
				fmt.Println("Downloaded NSRL bloom filter to " + nsrlDst)
			}
		}
	} else {
		fmt.Println("NSRL bloom filter already present.")
	}

	return dbdirerr, logsdirerr, trovedberr, siegfriederr
}

// DownloadNSRLBloom downloads the pre-built NSRL bloom filter from NSRLBloomURL.
func DownloadNSRLBloom(dst string) error {
	resp, err := http.Get(NSRLBloomURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("could not download NSRL bloom filter, server returned: " + resp.Status)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// copyNSRLBloom tries to find and copy nsrl.bloom from known locations
// into the destination path. It looks next to the binary and in db/.
func copyNSRLBloom(dst string) error {
	candidates := []string{
		// Relative to CWD (repo root or dist bundle)
		filepath.Join("db", "nsrl.bloom"),
		// Two levels up from cmd/ftrove/ to repo root
		filepath.Join("..", "..", "db", "nsrl.bloom"),
	}
	// Also check next to the running binary
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(exeDir, "db", "nsrl.bloom"),
			// Binary in cmd/ftrove/, bloom in repo root db/
			filepath.Join(exeDir, "..", "..", "db", "nsrl.bloom"),
		)
	}

	for _, src := range candidates {
		absSrc, _ := filepath.Abs(src)
		absDst, _ := filepath.Abs(dst)
		if absSrc == absDst {
			continue // already in place
		}
		if _, err := os.Stat(src); err == nil {
			return copyFile(src, dst)
		}
	}
	return fmt.Errorf("nsrl.bloom not found in any known location")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
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
			fmt.Println("ERROR: Legacy nsrl.db detected. Run 'task nsrl:build-all' or rebuild with admftrove --creatensrl to create nsrl.bloom.")
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
