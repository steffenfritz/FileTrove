package filetrove

import "fmt"

// PrintLicense prints a short license text
func PrintLicense(version string, build string) {
	fmt.Println("\n" +
		"FileTrove Copyright (C) 2023  Steffen Fritz <steffen@fritz.wtf> \n\n    " +
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
		"Version: " + version + " Build: " + build + "\n",
	)
}

// PrintBanner prints a pre-generated ascii banner with the program name
func PrintBanner() {
	fmt.Println("\no--o   o     o-O-o                   \n|    o |       |                     \nO-o    | o-o   |   o-o o-o o   o o-o \n|    | | |-'   |   |   | |  \\ /  |-' \no    | o o-o   o   o   o-o   o   o-o \n                                     \n                                     ")
}
