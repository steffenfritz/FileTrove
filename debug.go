package filetrove

import (
	info "github.com/elastic/go-sysinfo"
	"os"
	"runtime"
	"strconv"
)

// DebugCreateDebugPackage creates the file for compiling information into a debug package
func DebugCreateDebugPackage() (os.File, error) {
	file, err := os.OpenFile("debug_ftrove", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return os.File{}, err
	}
	// We don't close the file handle here, this must be done after all information are gathered
	file.WriteString("BEGIN FILETROVE DEBUG\n")

	return *file, err
}

// DebugCheckInstalled checks if FileTrove is installed by checking if the database exists
func DebugCheckInstalled(fd os.File) error {
	filePath := "db/filetrove.db"
	_, err := fd.WriteString("\nBEGIN INSTALLED\n")

	_, err = os.Stat(filePath)
	if err == nil {
		fd.WriteString("FileTrove installed: True")
		fd.WriteString("\nEND INSTALLED\n")
		return nil
	}
	if os.IsNotExist(err) {
		fd.WriteString("FileTrove installed: False")
		fd.WriteString("\nEND INSTALLED\n")
		return nil
	}

	return err
}

// DebugHostinformation writes host stats and returns on error
func DebugHostinformation(fd os.File) error {
	_, err := fd.WriteString("\nBEGIN HOST INFORMATION\n")
	infoOnHost, err := info.Host()

	if err != nil {
		return err
	}

	fd.WriteString("OSName: " + infoOnHost.Info().OS.Name + "\n")
	fd.WriteString("OSVersion: " + infoOnHost.Info().OS.Version + "\n")
	fd.WriteString("Platform: " + infoOnHost.Info().OS.Platform + "\n")
	mem, err := infoOnHost.Memory()
	if err != nil {
		fd.WriteString("Mem: " + err.Error())
		return err
	}
	fd.WriteString("Memory Total: " + strconv.Itoa(int(mem.Total)) + "\n")
	fd.WriteString("Memory Used: " + strconv.Itoa(int(mem.Used)) + "\n")
	fd.WriteString("Go Runtime: " + runtime.Version() + "\n")

	fd.WriteString("END HOST INFORMATION\n")

	return nil
}

// DebugWriteFlags takes parsed flags from main and writes them to the diag file
func DebugWriteFlags(fd os.File, args []string) error {
	_, err := fd.WriteString("\nBEGIN FLAGS\n")
	if err != nil {
		return err
	}
	for k, v := range args {
		_, err = fd.WriteString(strconv.Itoa(k) + ": " + v + "\n")
	}

	fd.WriteString("END FLAGS\n")

	return err
}

func DebugWriteFileList(fd os.File, filelist []string, dirlist []string) error {
	fd.WriteString("\nBEGIN FILELIST\n")
	for _, file := range filelist {
		_, err := fd.WriteString(file + "\n")
		if err != nil {
			return err
		}
	}
	fd.WriteString("END FILELIST\n")

	fd.WriteString("\nBEGIN DIRLIST\n")
	for _, dir := range dirlist {
		_, err := fd.WriteString(dir + "\n")
		if err != nil {
			return err
		}
	}
	fd.WriteString("END DIRLIST\n")

	return nil
}

// DebugCollectLogs collects log files
