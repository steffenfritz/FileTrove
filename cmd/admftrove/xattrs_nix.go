package main

import (
	"github.com/pkg/xattr"
)

// XattrsCheck checks if a file has xattr and return a list. This check depends on the filesystem.
// Filesystems like ext3/4, btrfs, xfs support xattr.
func XattrsCheck(path string) (bool, []string, error) {
	var list []string
	var err error
	if list, err = xattr.List(path); err != nil {
		return false, list, err
	}

	if len(list) <= 0 {
		return false, list, nil
	}

	return true, list, nil
}

// ADSCheckOnLinux checks if a file has an alternative data stream. This check depends on the filesystem.
// Filesystems like NTFS support ADSs.
// Note: It is difficult to have this check implemented on *Nix and the same way on Windows systems.
//
//	Some FS on *Nix accept ":" and don't treat it as a separator. Even drivers handle that differently.
//	Windows provides an API syscall here that would be the best solution. However, that does not work
//	on Linux.
//	Probably we need to use build tags here and implement it for at least two "flavours".
func ADSCheckOnLinux(path string) (bool, []string, error) {
	return false, nil, nil
}
