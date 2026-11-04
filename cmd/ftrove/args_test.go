package main

import (
	"strings"
	"testing"
)

func TestCheckUnexpectedArgs(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		inDir        string
		wantErr      bool
		wantContains []string
	}{
		{"no args", nil, "corpus/", false, nil},
		{"no args without -i", []string{}, "", false, nil},
		{"shell glob", []string{"corpus/b", "corpus/c"}, "corpus/a", true,
			[]string{`"corpus/b" "corpus/c"`, "-i expects exactly one directory", "ftrove -i corpus\n", "hidden files"}},
		{"shell glob with trailing slash", []string{"corpus/b/", "corpus/c/"}, "corpus/a/", true,
			[]string{`"corpus/b/" "corpus/c/"`, "ftrove -i corpus\n"}},
		{"glob in current directory", []string{"b", "c"}, "a", true,
			[]string{"ftrove -i .\n"}},
		{"parent with space", []string{"my corpus/b"}, "my corpus/a", true,
			[]string{`ftrove -i "my corpus"`}},
		{"args without -i", []string{"corpus/"}, "", true,
			[]string{`"corpus/"`, "Pass the directory to scan with -i", "hidden files"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkUnexpectedArgs(tt.args, tt.inDir)
			if (err != nil) != tt.wantErr {
				t.Fatalf("checkUnexpectedArgs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				return
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("checkUnexpectedArgs() message = %q, want it to contain %q", err.Error(), want)
				}
			}
		})
	}
}
