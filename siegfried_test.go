package filetrove

import (
	"github.com/richardlehane/siegfried"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestSiegfriedIdent(t *testing.T) {
	// Prepare siegfried signature for SiegfriedIdent test
	s, err := siegfried.Load(filepath.Join("resources", "siegfried.sig"))
	if err != nil {
		println("Could not read siegfried's database.")
		os.Exit(1)
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
