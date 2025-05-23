package filetrove

import (
	"reflect"
	"testing"
	"time"
)

func TestGetFileTimes(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    FileTime
		wantErr bool
	}{
		{"File Time white.jpg", args{"testdata/white.jpg"}, FileTime{Btime: time.Date(2025, time.March, 03, 13, 18, 34, 98689545, time.Local)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFileTimes(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileTimes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Btime, tt.want.Btime) {
				t.Errorf("GetFileTimes() got = %v, want %v", got.Btime, tt.want.Btime)
			}
		})
	}
}
