//go:build darwin || linux

package filetrove

import (
	"fmt"
	"os"
	"syscall"
)

// getFileOwner gets the owner and group ids of a file on Unix and Unix like systems
func GetFileOwner(path string) (*FileOwnerInfo, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	// Cast  in syscall.Stat_t casten
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("cannot cast file info to Stat_t")
	}

	uid := stat.Uid
	gid := stat.Gid

	return &FileOwnerInfo{
		Owner: fmt.Sprintf("%d", uid),
		Group: fmt.Sprintf("%d", gid),
	}, nil
}
