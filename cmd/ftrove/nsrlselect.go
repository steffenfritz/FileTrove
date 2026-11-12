package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// bloomVariantOrder lists the NSRL variants from most to least comprehensive. Without an
// explicit --nsrl-variant, a scan uses the first installed one.
var bloomVariantOrder = []string{"all", "mobile", "modern"}

// legacyBloomName is the bloom filter name used before variant-specific file names.
const legacyBloomName = "nsrl.bloom"

// bloomCandidates returns the bloom filter paths below dbDir in the order a scan tries them.
// An explicit variant is never replaced by another variant, only by the legacy file.
func bloomCandidates(dbDir, variant string, explicit bool) []string {
	variants := bloomVariantOrder
	if explicit {
		variants = []string{variant}
	}
	candidates := make([]string, 0, len(variants)+1)
	for _, v := range variants {
		candidates = append(candidates, filepath.Join(dbDir, "nsrl-"+v+".bloom"))
	}
	return append(candidates, filepath.Join(dbDir, legacyBloomName))
}

// selectBloomPath returns the first existing bloom filter from bloomCandidates and true,
// or an empty path and false if none exists.
func selectBloomPath(dbDir, variant string, explicit bool) (string, bool) {
	for _, candidate := range bloomCandidates(dbDir, variant, explicit) {
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() {
			return candidate, true
		}
	}
	return "", false
}

// bloomNotFoundReason explains which bloom filter files a scan looked for without success.
func bloomNotFoundReason(dbDir, variant string, explicit bool) string {
	return "none of " + strings.Join(bloomCandidates(dbDir, variant, explicit), ", ") + " found"
}

// reportNSRLDisabled tells the user on stderr that the scan runs without NSRL checks and why.
// It is printed independent of -V, so a scan is never silently missing NSRL results.
func reportNSRLDisabled(stderr io.Writer, reason string) {
	fmt.Fprintln(stderr, "WARNING: NSRL bloom filter not available, NSRL checks disabled: "+reason)
}
