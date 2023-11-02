package filetrove

import (
	"encoding/json"
	"os"
)

// DublinCore is a struct that holds 15 core elements of DC
// https://datatracker.ietf.org/doc/html/rfc5013
type DublinCore struct {
	Title       string `json:"title"`
	Creator     string `json:"creator"`
	Contributor string `json:"contributor"`
	Publisher   string `json:"publisher"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
	Date        string `json:"date"`
	Language    string `json:"language"`
	Type        string `json:"type"`
	Format      string `json:"format"`
	Identifier  string `json:"identifier"`
	Source      string `json:"source"`
	Relation    string `json:"relation"`
	Rights      string `json:"rights"`
	Coverage    string `json:"coverage"`
}

// ReadDC reads a json file and unmarshals it into the DublinCore struct
func ReadDC(dcjson string) (DublinCore, error) {
	var dc DublinCore

	filecontent, err := os.ReadFile(dcjson)
	if err != nil {
		return dc, err
	}

	err = json.Unmarshal(filecontent, &dc)

	return dc, err
}
