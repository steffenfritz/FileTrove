package filetrove

import (
	"io/fs"
	"path/filepath"
)

// CreateFileList creates a list of file paths and a directory listing
func CreateFileList(rootDir string) ([]string, []string, error) {
	var fileList []string
	var dirList []string
	err := filepath.WalkDir(rootDir, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Type().IsRegular() {
			fileList = append(fileList, path)
		} else if info.IsDir() {
			dirList = append(dirList, path)
		}
		return nil
	})
	if err != nil {
		return fileList, dirList, err
	}
	return fileList, dirList, nil
}
