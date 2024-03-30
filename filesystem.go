//go:build linux
// +build linux

package filetrove

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// GetFileSystemType returns the filesystem type of the specified mount point or volume.
func GetFileSystemType(path string) (string, error) {
	if runtime.GOOS == "linux" {
		return getLinuxFileSystemType(path)
	}
	if runtime.GOOS == "windows" {
		return getWindowsFileSystemType(path)
	}
	if runtime.GOOS == "darwin" {
		return getMacFileSystemType(path)
	} else {
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getLinuxFileSystemType returns the filesystem type of the specified mount point on Linux.
func getLinuxFileSystemType(path string) (string, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 3 && fields[1] == path {
			return fields[2], nil
		}
	}
	return "", fmt.Errorf("mount point not found: %s", path)
}

// getWindowsFileSystemType returns the filesystem type of the specified volume on Windows.
func getWindowsFileSystemType(volumeName string) (string, error) {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getVolumeInformation := kernel32.NewProc("GetVolumeInformationW")

	var volumeSerialNumber, maximumComponentLength, fileSystemFlags uint32
	var fileSystemNameBuffer [255]uint16

	volumeNameBuffer := syscall.StringToUTF16Ptr(volumeName)

	_, _, err := getVolumeInformation.Call(
		uintptr(unsafe.Pointer(volumeNameBuffer)),
		uintptr(unsafe.Pointer(&fileSystemNameBuffer[0])),
		uintptr(len(fileSystemNameBuffer)),
		uintptr(unsafe.Pointer(&volumeSerialNumber)),
		uintptr(unsafe.Pointer(&maximumComponentLength)),
		uintptr(unsafe.Pointer(&fileSystemFlags)),
		0,
		0,
	)
	if err != nil {
		return "", err
	}

	fileSystemName := syscall.UTF16ToString(fileSystemNameBuffer[:])
	return fileSystemName, nil
}

// getMacFileSystemType returns the filesystem type of the specified mount point on macOS.
func getMacFileSystemType(path string) (string, error) {
	var stat syscall.Statfs_t

	err := syscall.Statfs(path, &stat)
	if err != nil {
		return "", err
	}

	fsType := string(stat.Fstypename[:])
	return fsType, nil
}
