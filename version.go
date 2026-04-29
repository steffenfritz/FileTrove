package filetrove

import (
	"database/sql"
	"fmt"
	"strings"
)

// CheckVersion checks if the binary version is compatible with the database.
// Only the base version (the part before '+') is compared, so that builds
// from different commits but the same release (e.g. 1.0.0-BETA-4+abc vs
// 1.0.0-BETA-4+def) are treated as compatible.
func CheckVersion(db *sql.DB, version string) (bool, string, error) {
	var dbversion string

	query := "SELECT version FROM filetrove;"
	resultrow := db.QueryRow(query)

	err := resultrow.Scan(&dbversion)
	if err != nil {
		return false, "", err
	}

	baseVersion := strings.SplitN(version, "+", 2)[0]
	baseDbVersion := strings.SplitN(dbversion, "+", 2)[0]

	if baseDbVersion == baseVersion {
		return true, "", nil
	}

	return false, dbversion, nil
}

// PrintLicense prints a short license text
// func PrintLicense(version string, build string) {
func PrintLicense(version string) {
	fmt.Println("\n" +
		"FileTrove Copyright (C) 2023-2026  Steffen Fritz <steffen@fritz.wtf> \n\n    " +
		"This program is free software: you can redistribute it and/or modify\n    " +
		"it under the terms of the GNU Affero General Public License as published\n    " +
		"by the Free Software Foundation, either version 3 of the License, or\n    " +
		"(at your option) any later version.\n\n    " +
		"This program is distributed in the hope that it will be useful,\n    " +
		"but WITHOUT ANY WARRANTY; without even the implied warranty of\n    " +
		"MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the\n    " +
		"GNU Affero General Public License for more details.\n\n    " +
		"You should have received a copy of the GNU Affero General Public License\n    " +
		"along with this program.  If not, see <https://www.gnu.org/licenses/>.\n\n" +
		"Version: " + version + "\n",
	)
}

// PrintBanner prints a pre-generated ascii banner with the program name
func PrintBanner() {
	fmt.Println("\no--o   o     o-O-o                   \n|    o |       |                     \nO-o    | o-o   |   o-o o-o o   o o-o \n|    | | |-'   |   |   | |  \\ /  |-' \no    | o o-o   o   o   o-o   o   o-o \n                                     \n                                     ")
}
