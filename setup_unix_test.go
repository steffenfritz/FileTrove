//go:build unix

package filetrove

import (
	"fmt"
	"os"
	"syscall"
	"testing"
)

// TestMain creates fixtures that can't be tracked in git: a file with no read
// permission and a named pipe. Both are needed by TestCreateFileList to verify
// that CreateFileList records unreadable files and special files in the
// skipped list.
func TestMain(m *testing.M) {
	const (
		noaccess = "testdata/noaccess.rtf"
		pipe     = "testdata/testpipe"
	)

	// Remove any leftover from a prior failed run (may have 0o000 perms).
	_ = os.Chmod(noaccess, 0o600)
	_ = os.Remove(noaccess)

	if err := os.WriteFile(noaccess, []byte{}, 0o000); err != nil {
		fmt.Fprintf(os.Stderr, "test setup: create %s: %v\n", noaccess, err)
		os.Exit(1)
	}
	// WriteFile may not honor 0o000 on some umasks; force it.
	if err := os.Chmod(noaccess, 0o000); err != nil {
		fmt.Fprintf(os.Stderr, "test setup: chmod %s: %v\n", noaccess, err)
		os.Exit(1)
	}

	_ = os.Remove(pipe)
	if err := syscall.Mkfifo(pipe, 0o600); err != nil {
		fmt.Fprintf(os.Stderr, "test setup: mkfifo %s: %v\n", pipe, err)
		os.Exit(1)
	}

	code := m.Run()

	// Restore mode before removing so the file is actually deletable.
	_ = os.Chmod(noaccess, 0o600)
	_ = os.Remove(noaccess)
	_ = os.Remove(pipe)

	os.Exit(code)
}
