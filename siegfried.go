package filetrove

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/richardlehane/siegfried"
)

// SiegfriedType is a struct for all the strings siegfried returns
type SiegfriedType struct {
	FileName            string
	SizeInByte          int64
	Registry            string
	FMT                 string
	FormatName          string
	FormatVersion       string
	MIMEType            string
	IdentificationNote  string
	IdentificationProof string
	SiegOutput          string
}

const SiegfriedVersion = "1_11"

// GetSiegfriedDB downloads the signature db
func GetSiegfriedDB(installPath string) error {
	sigurl := "https://www.itforarchivists.com/siegfried/latest/" + SiegfriedVersion + "/default"
	// We download siegfried's database derived from DROID here, PRONOM based
	// TODO: Check license note with Richard
	resp, err := http.Get(sigurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Could not download siegfried signature. Server returned: " + resp.Status)
	}

	// Create the signature file in the db subdirectory
	out, err := os.Create(filepath.Join(installPath, "db", "siegfried.sig"))
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to the signature file
	_, err = io.Copy(out, resp.Body)
	return err

}

// SiegfriedIdent gets PRONOM metadata and the size of a single file
func SiegfriedIdent(s *siegfried.Siegfried, inFile string) (SiegfriedType, error) {
	var MetaOneFile SiegfriedType
	var oneFile string

	f, err := os.Open(inFile)
	if err != nil {
		return MetaOneFile, err
	}
	defer f.Close()

	fi, _ := f.Stat()
	MetaOneFile.SizeInByte = fi.Size()
	if fi.Size() == 0 {
		return MetaOneFile, err
	}

	ids, err := s.Identify(f, "", "")
	if err != nil {
		return MetaOneFile, err
	}

	for _, id := range ids {
		values := id.Values()
		MetaOneFile.Registry = values[0]
		MetaOneFile.FMT = values[1]
		MetaOneFile.FormatName = values[2]
		MetaOneFile.FormatVersion = values[3]
		MetaOneFile.MIMEType = values[4]
		MetaOneFile.IdentificationNote = values[5]
		MetaOneFile.IdentificationProof = values[6]

	}

	MetaOneFile.FileName = inFile
	MetaOneFile.SizeInByte = fi.Size()
	MetaOneFile.SiegOutput = oneFile

	return MetaOneFile, nil

}
