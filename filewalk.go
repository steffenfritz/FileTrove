package filetrove

import (
	"io/fs"
	"os"
	"path/filepath"
)

// CreateFileList walks rootDir and returns three lists: regular files, directories, and skipped paths.
//
// Skipped paths include symlinks (not followed), special files (sockets, devices, FIFOs),
// and any path that could not be accessed (e.g. permission denied, stale network mount).
// The walk continues past inaccessible entries rather than aborting.
//
// Note: filepath.WalkDir crosses filesystem boundaries, including mounted network shares.
// Callers that need to stay within a single device should compare the device ID of each
// entry (via DirEntry.Info().Sys()) against the root device.
func CreateFileList(rootDir string) ([]string, []string, []string, error) {
	var fileList []string
	var dirList []string
	var skippedList []string

	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			// Inaccessible entry (permission denied, stale NFS mount, etc.).
			// Record and continue rather than aborting the whole walk.
			skippedList = append(skippedList, path)
			return nil
		}

		switch {
		case info.Type().IsRegular():
			f, openErr := os.Open(path)
			if openErr != nil {
				skippedList = append(skippedList, path)
				break
			}
			f.Close()
			fileList = append(fileList, path)
		case info.IsDir():
			dirList = append(dirList, path)
		case info.Type()&fs.ModeSymlink != 0:
			// Symlinks are not followed by filepath.WalkDir. Record them so
			// callers are aware they exist on the filesystem.
			skippedList = append(skippedList, path)
		default:
			// Special files: sockets, named pipes, device nodes.
			skippedList = append(skippedList, path)
		}

		return nil
	})

	return fileList, dirList, skippedList, err
}
