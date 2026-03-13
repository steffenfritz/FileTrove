package filetrove

import (
	"os"
	"reflect"
	"testing"
	"time"
)

func TestGetFileTimes(t *testing.T) {
	knownTime := time.Date(2025, time.March, 3, 13, 18, 34, 0, time.Local)
	if err := os.Chtimes("testdata/white.jpg", knownTime, knownTime); err != nil {
		t.Fatalf("Could not set test file times: %v", err)
	}

	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    FileTime
		wantErr bool
	}{
		{"File Time white.jpg", args{"testdata/white.jpg"}, FileTime{Mtime: knownTime}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileTimes(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Mtime, tt.want.Mtime) {
				t.Errorf("GetFileTimes() got = %v, want %v", got.Mtime, tt.want.Mtime)
			}
		})
	}
}
