package filetrove

import (
	"github.com/richardlehane/siegfried"
	"os"
	"strconv"
)

// GetSiegfriedDB downloads the signature db
func GetSiegfriedDB() {}

func siegfriedIdent(s *siegfried.Siegfried, inFile string) (bool, string) {
	var oneFile string
	var resultBool bool

	f, err := os.Open(inFile)
	if err != nil {
		return resultBool, err.Error()
	}

	defer f.Close()

	fi, _ := f.Stat()
	if fi.Size() == 0 {
		//return resultBool, "\"" + inFile + "\",0,,,,,,,"
		return true, "\"" + inFile + "\",0,,,,,,,"
	}

	ids, err := s.Identify(f, "", "")
	if err != nil {
		ret := inFile + " : " + err.Error()
		return resultBool, ret
	}

	for _, id := range ids {
		values := id.Values()
		for _, value := range values {
			oneFile += "\"" + value + "\"" + ","
		}

		oneFile = "\"" + inFile + "\",\"" + strconv.Itoa(int(fi.Size())) + "\"," + oneFile[:len(oneFile)-1] // remove last comma
	}

	return true, oneFile

}
