package filetrove

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// generateTestHashes creates deterministic SHA1 hashes from sequential strings
func generateTestHashes(n int) []string {
	hashes := make([]string, n)
	for i := 0; i < n; i++ {
		h := sha1.Sum([]byte(fmt.Sprintf("test-file-%d", i)))
		hashes[i] = hex.EncodeToString(h[:])
	}
	return hashes
}

// writeHashFile writes hashes to a temporary file, one per line
func writeHashFile(t *testing.T, hashes []string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "nsrl-hashes-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	for _, h := range hashes {
		fmt.Fprintln(f, h)
	}
	f.Close()
	return f.Name()
}

func TestCreateAndLoadBloom(t *testing.T) {
	hashes := generateTestHashes(1000)
	hashFile := writeHashFile(t, hashes)
	bloomFile := filepath.Join(t.TempDir(), "test.bloom")

	err := CreateNSRLBloom(hashFile, "test-v1", bloomFile, 1000, 0.0001)
	if err != nil {
		t.Fatalf("CreateNSRLBloom failed: %v", err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatalf("LoadNSRL failed: %v", err)
	}

	// Verify all inserted hashes are found (zero false negatives)
	for i, h := range hashes {
		if !nf.Contains(h) {
			t.Errorf("hash %d not found in bloom filter (false negative): %s", i, h)
		}
	}
}

func TestBloomVersionMetadata(t *testing.T) {
	hashes := generateTestHashes(100)
	hashFile := writeHashFile(t, hashes)
	bloomFile := filepath.Join(t.TempDir(), "test.bloom")

	err := CreateNSRLBloom(hashFile, "RDS_2026.03.1-modern", bloomFile, 100, 0.0001)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatal(err)
	}

	if nf.Version != "RDS_2026.03.1-modern" {
		t.Errorf("expected version 'RDS_2026.03.1-modern', got '%s'", nf.Version)
	}
	if nf.HashType != "sha1" {
		t.Errorf("expected hash type 'sha1', got '%s'", nf.HashType)
	}
	if nf.Items != 100 {
		t.Errorf("expected 100 items, got %d", nf.Items)
	}
	if nf.FPR != 0.0001 {
		t.Errorf("expected FPR 0.0001, got %f", nf.FPR)
	}
}

func TestBloomFalsePositiveRate(t *testing.T) {
	hashes := generateTestHashes(10000)
	hashFile := writeHashFile(t, hashes)
	bloomFile := filepath.Join(t.TempDir(), "test.bloom")

	err := CreateNSRLBloom(hashFile, "test-fpr", bloomFile, 10000, 0.001)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatal(err)
	}

	// Test with hashes that were NOT inserted
	falsePositives := 0
	testCount := 100000
	for i := 10000; i < 10000+testCount; i++ {
		h := sha1.Sum([]byte(fmt.Sprintf("test-file-%d", i)))
		hash := hex.EncodeToString(h[:])
		if nf.Contains(hash) {
			falsePositives++
		}
	}

	observedFPR := float64(falsePositives) / float64(testCount)
	// Allow up to 5x the target FPR as tolerance
	if observedFPR > 0.005 {
		t.Errorf("observed FPR %.4f exceeds 5x target (0.001), got %d false positives out of %d", observedFPR, falsePositives, testCount)
	}
	t.Logf("FPR: %.6f (%d/%d)", observedFPR, falsePositives, testCount)
}

func TestBloomCaseInsensitivity(t *testing.T) {
	hashes := []string{
		"da39a3ee5e6b4b0d3255bfef95601890afd80709",
		"AABBCCDD11223344556677889900AABBCCDDEEFF",
	}
	hashFile := writeHashFile(t, hashes)
	bloomFile := filepath.Join(t.TempDir(), "test.bloom")

	err := CreateNSRLBloom(hashFile, "test-case", bloomFile, 10, 0.0001)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatal(err)
	}

	// Both upper and lower case should match
	if !nf.Contains("da39a3ee5e6b4b0d3255bfef95601890afd80709") {
		t.Error("lowercase hash not found")
	}
	if !nf.Contains("DA39A3EE5E6B4B0D3255BFEF95601890AFD80709") {
		t.Error("uppercase hash not found")
	}
	if !nf.Contains("aabbccdd11223344556677889900aabbccddeeff") {
		t.Error("originally uppercase hash not found via lowercase query")
	}
}

func TestBloomCorruptFile(t *testing.T) {
	corruptFile := filepath.Join(t.TempDir(), "corrupt.bloom")
	os.WriteFile(corruptFile, []byte("not a valid bloom filter"), 0644)

	_, err := LoadNSRL(corruptFile)
	if err == nil {
		t.Error("expected error loading corrupt file, got nil")
	}
}

func TestBloomEmptyFile(t *testing.T) {
	hashFile := writeHashFile(t, []string{})
	bloomFile := filepath.Join(t.TempDir(), "empty.bloom")

	err := CreateNSRLBloom(hashFile, "empty", bloomFile, 1, 0.0001)
	if err != nil {
		t.Fatal(err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatal(err)
	}

	if nf.Items != 0 {
		t.Errorf("expected 0 items, got %d", nf.Items)
	}
}

func TestBloomWithRealNSRL(t *testing.T) {
	bloomFile := "db/nsrl.bloom"
	if _, err := os.Stat(bloomFile); os.IsNotExist(err) {
		t.Skip("db/nsrl.bloom not present; run 'task nsrl:build-all' first")
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatalf("LoadNSRL failed: %v", err)
	}

	// Known NSRL hashes from testdata (extracted from RDS 2026.03.1)
	knownFile := "testdata/nsrl_known_hashes.txt"
	f, err := os.Open(knownFile)
	if err != nil {
		t.Fatalf("open %s: %v", knownFile, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	line := 0
	for scanner.Scan() {
		line++
		hash := scanner.Text()
		if hash == "" {
			continue
		}
		if !nf.Contains(hash) {
			t.Errorf("line %d: known NSRL hash not found: %s", line, hash)
		}
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("reading %s: %v", knownFile, err)
	}
	t.Logf("all %d known hashes found in bloom filter (version: %s, items: %d)", line, nf.Version, nf.Items)
}

func TestBloomMissingFile(t *testing.T) {
	_, err := LoadNSRL("/nonexistent/path/nsrl.bloom")
	if err == nil {
		t.Error("expected error loading nonexistent file, got nil")
	}
}

func TestBloomAutoCount(t *testing.T) {
	hashes := generateTestHashes(500)
	hashFile := writeHashFile(t, hashes)
	bloomFile := filepath.Join(t.TempDir(), "autocount.bloom")

	// estimatedItems=0 triggers auto-count (two-pass: count lines then scan)
	err := CreateNSRLBloom(hashFile, "test-autocount", bloomFile, 0, 0.0001)
	if err != nil {
		t.Fatalf("CreateNSRLBloom with auto-count failed: %v", err)
	}

	nf, err := LoadNSRL(bloomFile)
	if err != nil {
		t.Fatalf("LoadNSRL failed: %v", err)
	}

	if nf.Items != 500 {
		t.Errorf("expected 500 items, got %d", nf.Items)
	}
	for i, h := range hashes {
		if !nf.Contains(h) {
			t.Errorf("hash %d not found after auto-count build: %s", i, h)
		}
	}
}

func TestBloomStdinRequiresEstimate(t *testing.T) {
	bloomFile := filepath.Join(t.TempDir(), "stdin.bloom")

	err := CreateNSRLBloom("-", "test-stdin", bloomFile, 0, 0.0001)
	if err == nil {
		t.Fatal("expected error when reading from stdin without --nsrl-estimate, got nil")
	}
}
