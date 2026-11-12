package main

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestSelectBloomPath(t *testing.T) {
	tests := []struct {
		name      string
		files     []string
		variant   string
		explicit  bool
		want      string
		wantFound bool
	}{
		{"only modern", []string{"nsrl-modern.bloom"}, "all", false, "nsrl-modern.bloom", true},
		{"only mobile", []string{"nsrl-mobile.bloom"}, "all", false, "nsrl-mobile.bloom", true},
		{"modern and all", []string{"nsrl-modern.bloom", "nsrl-all.bloom"}, "all", false, "nsrl-all.bloom", true},
		{"modern and mobile", []string{"nsrl-modern.bloom", "nsrl-mobile.bloom"}, "all", false, "nsrl-mobile.bloom", true},
		{"variant before legacy", []string{"nsrl.bloom", "nsrl-modern.bloom"}, "all", false, "nsrl-modern.bloom", true},
		{"explicit modern with modern and all", []string{"nsrl-modern.bloom", "nsrl-all.bloom"}, "modern", true, "nsrl-modern.bloom", true},
		{"explicit mobile with only modern", []string{"nsrl-modern.bloom"}, "mobile", true, "", false},
		{"explicit all falls back to legacy", []string{"nsrl.bloom", "nsrl-modern.bloom"}, "all", true, "nsrl.bloom", true},
		{"only legacy", []string{"nsrl.bloom"}, "all", false, "nsrl.bloom", true},
		{"only legacy with explicit variant", []string{"nsrl.bloom"}, "modern", true, "nsrl.bloom", true},
		{"unknown variant file ignored", []string{"nsrl-custom.bloom"}, "all", false, "", false},
		{"explicit custom variant", []string{"nsrl-custom.bloom"}, "custom", true, "nsrl-custom.bloom", true},
		{"empty directory", nil, "all", false, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbDir := t.TempDir()
			for _, f := range tt.files {
				if err := os.WriteFile(filepath.Join(dbDir, f), []byte("x"), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			want := ""
			if tt.want != "" {
				want = filepath.Join(dbDir, tt.want)
			}
			got, found := selectBloomPath(dbDir, tt.variant, tt.explicit)
			if got != want || found != tt.wantFound {
				t.Errorf("selectBloomPath() = %q, %v; want %q, %v", got, found, want, tt.wantFound)
			}
		})
	}
}

func TestSelectBloomPathIgnoresDirectory(t *testing.T) {
	dbDir := t.TempDir()
	if err := os.Mkdir(filepath.Join(dbDir, "nsrl-all.bloom"), 0o755); err != nil {
		t.Fatal(err)
	}
	if got, found := selectBloomPath(dbDir, "all", false); found {
		t.Errorf("selectBloomPath() = %q, true; want no match for a directory", got)
	}
}

func TestBloomCandidates(t *testing.T) {
	tests := []struct {
		name     string
		variant  string
		explicit bool
		want     []string
	}{
		{"auto", "all", false, []string{"db/nsrl-all.bloom", "db/nsrl-mobile.bloom", "db/nsrl-modern.bloom", "db/nsrl.bloom"}},
		{"explicit", "mobile", true, []string{"db/nsrl-mobile.bloom", "db/nsrl.bloom"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := make([]string, len(tt.want))
			for i, w := range tt.want {
				want[i] = filepath.FromSlash(w)
			}
			if got := bloomCandidates("db", tt.variant, tt.explicit); !slices.Equal(got, want) {
				t.Errorf("bloomCandidates() = %v, want %v", got, want)
			}
		})
	}
}

func TestReportNSRLDisabled(t *testing.T) {
	tests := []struct {
		name         string
		reason       string
		wantContains []string
	}{
		{"not found, auto", bloomNotFoundReason("db", "all", false),
			[]string{"none of db/nsrl-all.bloom, db/nsrl-mobile.bloom, db/nsrl-modern.bloom, db/nsrl.bloom found"}},
		{"not found, explicit", bloomNotFoundReason("db", "mobile", true),
			[]string{"none of db/nsrl-mobile.bloom, db/nsrl.bloom found"}},
		{"cannot be loaded", "db/nsrl-modern.bloom: could not decode NSRL bloom filter: unexpected EOF",
			[]string{"db/nsrl-modern.bloom", "unexpected EOF"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stderr bytes.Buffer
			reportNSRLDisabled(&stderr, tt.reason)
			got := stderr.String()
			if !strings.HasPrefix(got, "WARNING: NSRL bloom filter not available, NSRL checks disabled: ") {
				t.Errorf("stderr = %q, want WARNING prefix", got)
			}
			if !strings.HasSuffix(got, "\n") || strings.Count(got, "\n") != 1 {
				t.Errorf("stderr = %q, want exactly one line", got)
			}
			for _, want := range tt.wantContains {
				if want = filepath.FromSlash(want); !strings.Contains(got, want) {
					t.Errorf("stderr = %q, want it to contain %q", got, want)
				}
			}
		})
	}
}
