package filetrove

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/richardlehane/siegfried"
)

func TestSiegfriedIdent(t *testing.T) {
	// Prepare siegfried signature for SiegfriedIdent test
	// Try db/ first (where --install and dist:bundle place it), fall back to legacy resources/
	sigPaths := []string{
		filepath.Join("db", "siegfried.sig"),
		filepath.Join("resources", "siegfried.sig"),
	}
	var s *siegfried.Siegfried
	var err error
	for _, p := range sigPaths {
		s, err = siegfried.Load(p)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Skip("Skipping: could not read siegfried's database:", err)
	}

	type args struct {
		s      *siegfried.Siegfried
		inFile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test file 1", args{s: s, inFile: "testdata/transparent.png"}, "fmt/11", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SiegfriedIdent(tt.args.s, tt.args.inFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SiegfriedIdent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.FMT, tt.want) {
				t.Errorf("SiegfriedIdent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
