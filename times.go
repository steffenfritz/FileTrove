package filetrove

import (
	"time"

	"github.com/djherbis/times"
)

// FileTime holds all metadata times of a file
type FileTime struct {
	Atime time.Time
	Btime time.Time
	Ctime time.Time
	Mtime time.Time
}

// GetFileTimes returns a type that holds the access, change and birth time of a file if available.
func GetFileTimes(filename string) (FileTime, error) {
	var ft FileTime

	t, err := times.Stat(filename)
	if err != nil {
		return ft, err
	}

	ft.Atime = t.AccessTime()
	ft.Mtime = t.ModTime()

	if t.HasChangeTime() {
		ft.Ctime = t.ChangeTime()
	}

	if t.HasBirthTime() {
		ft.Btime = t.BirthTime()
	}

	return ft, nil
}
