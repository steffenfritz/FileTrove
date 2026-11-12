package filetrove

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// captureOutput runs f and returns what it wrote to os.Stdout and os.Stderr.
func captureOutput(t *testing.T, f func()) (stdout, stderr string) {
	t.Helper()
	outR, outW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	errR, errW, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	origOut, origErr := os.Stdout, os.Stderr
	os.Stdout, os.Stderr = outW, errW
	defer func() { os.Stdout, os.Stderr = origOut, origErr }()

	f()

	outW.Close()
	errW.Close()
	outBytes, _ := io.ReadAll(outR)
	errBytes, _ := io.ReadAll(errR)
	return string(outBytes), string(errBytes)
}

func TestCheckInstallWritesToStderr(t *testing.T) {
	t.Chdir(t.TempDir())
	if err := os.Mkdir("db", 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"siegfried.sig", "filetrove.db"} {
		if err := os.WriteFile(filepath.Join("db", name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	stdout, stderr := captureOutput(t, func() {
		_ = CheckInstall("dev")
	})

	if stdout != "" {
		t.Errorf("CheckInstall() stdout = %q, want empty", stdout)
	}
	if !strings.Contains(stderr, "nsrl bloom filter does not exist") {
		t.Errorf("CheckInstall() stderr = %q, want it to name the missing bloom filter", stderr)
	}
}
