package filetrove

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bits-and-blooms/bloom/v3"
)

// NSRLFilter wraps a Bloom filter with NSRL metadata
type NSRLFilter struct {
	Filter   *bloom.BloomFilter
	Version  string   // NSRL RDS version (e.g., "2026.03.1-modern")
	HashType string   // "sha1" (future: "sha256")
	FPR      float64  // target false positive rate
	Items    uint     // number of hashes inserted
	Subsets  []string // e.g., ["modern"], ["modern", "android", "ios"]
}

// Contains checks if a given SHA1 hash is present in the NSRL Bloom filter
func (nf *NSRLFilter) Contains(sha1hash string) bool {
	return nf.Filter.TestString(strings.ToLower(sha1hash))
}

// CreateNSRLBloom reads a newline-delimited SHA1 hash file and creates a Bloom filter.
// nsrlsourcefile may be "-" to read from stdin, in which case estimatedItems must be > 0.
// estimatedItems is a hint for filter sizing. If 0, the file is pre-scanned to count
// the actual number of hashes, which guarantees the target FPR is met.
// fpr is the target false positive rate (e.g., 0.0001 for 0.01%).
func CreateNSRLBloom(nsrlsourcefile string, nsrlversion string, nsrloutfile string, estimatedItems uint, fpr float64) error {
	var r io.Reader

	if nsrlsourcefile == "-" {
		if estimatedItems == 0 {
			return fmt.Errorf("--nsrl-estimate must be provided when reading from stdin")
		}
		r = os.Stdin
	} else {
		// If no estimate provided, count actual lines first so the filter is correctly sized.
		if estimatedItems == 0 {
			n, err := countNonEmptyLines(nsrlsourcefile)
			if err != nil {
				return fmt.Errorf("counting hashes: %w", err)
			}
			estimatedItems = n
		}
		f, err := os.Open(nsrlsourcefile)
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	// Ensure at least 1 to avoid zero-size filter
	if estimatedItems == 0 {
		estimatedItems = 1
	}

	filter := bloom.NewWithEstimates(estimatedItems, fpr)

	var count uint
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		hash := strings.TrimSpace(scanner.Text())
		if len(hash) == 0 {
			continue
		}
		filter.AddString(strings.ToLower(hash))
		count++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	nf := NSRLFilter{
		Filter:   filter,
		Version:  nsrlversion,
		HashType: "sha1",
		FPR:      fpr,
		Items:    count,
		Subsets:  []string{},
	}

	outFile, err := os.Create(nsrloutfile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := gob.NewEncoder(outFile)
	if err := encoder.Encode(&nf); err != nil {
		return err
	}

	return nil
}

// countNonEmptyLines counts non-empty lines in a file (single pass).
func countNonEmptyLines(path string) (uint, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var n uint
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.TrimSpace(scanner.Text()) != "" {
			n++
		}
	}
	return n, scanner.Err()
}

// LoadNSRL loads a serialized NSRLFilter from a .bloom file into memory
func LoadNSRL(nsrlbloomfile string) (*NSRLFilter, error) {
	file, err := os.Open(nsrlbloomfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var nf NSRLFilter
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&nf); err != nil {
		return nil, errors.New("could not decode NSRL bloom filter: " + err.Error())
	}

	if nf.Filter == nil {
		return nil, errors.New("NSRL bloom filter is empty or corrupt")
	}

	return &nf, nil
}
