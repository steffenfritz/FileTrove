package filetrove

import (
	"path"
	"reflect"
	"testing"
)

func TestGetXattr(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test set Xattr", args{path.Join("testdata", "textfile.txt")}, "testvalue", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetXattr(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetXattr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			if !reflect.DeepEqual(got["testkey"], tt.want) {
				t.Errorf("GetXattr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
