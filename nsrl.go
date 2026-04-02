package filetrove

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
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
// estimatedItems is the expected number of hashes (e.g., 40_000_000).
// fpr is the target false positive rate (e.g., 0.0001 for 0.01%).
func CreateNSRLBloom(nsrlsourcefile string, nsrlversion string, nsrloutfile string, estimatedItems uint, fpr float64) error {
	filter := bloom.NewWithEstimates(estimatedItems, fpr)

	file, err := os.Open(nsrlsourcefile)
	if err != nil {
		return err
	}
	defer file.Close()

	var count uint
	scanner := bufio.NewScanner(file)
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

	fmt.Printf("Bloom filter stats: estimated items: %d, actual items inserted: %d, target FPR: %.6f\n",
		estimatedItems, count, fpr)
	if count > estimatedItems {
		fmt.Printf("WARNING: actual item count (%d) exceeds estimated items (%d). "+
			"The real false positive rate will be significantly higher than the target %.6f. "+
			"Re-create the filter with --nsrl-estimate >= %d.\n",
			count, estimatedItems, fpr, count)
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
