package filetrove

import (
	"io"
	"net/http"
	"os"

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

// GetSiegfriedDB downloads the signature db
func GetSiegfriedDB(installPath string) error {
	sigurl := "https://www.itforarchivists.com/siegfried/latest/1_11/default"
	// We download siegfried's database derived from DROID here
	// TODO: Check license note with Richard
	resp, err := http.Get(sigurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the signature file in the db subdirectory
	out, err := os.Create(installPath + "/db/siegfried.sig")
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to the signature file
	_, err = io.Copy(out, resp.Body)
	return err

}

// SiegfiredIdent gets PRONOM metadata and the size of a single file
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

	//oneFile = "\"" + inFile + "\",\"" + strconv.Itoa(int(fi.Size())) + "\"," + oneFile[:len(oneFile)-1] // remove last comma
	MetaOneFile.FileName = inFile
	MetaOneFile.SizeInByte = fi.Size()
	MetaOneFile.SiegOutput = oneFile

	return MetaOneFile, nil

}
