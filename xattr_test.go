package filetrove

import (
	"path"
	"reflect"
	"testing"

	"github.com/pkg/xattr"
)

func TestGetXattr(t *testing.T) {
	testFile := path.Join("testdata", "textfile.txt")
	if err := xattr.Set(testFile, "user.testkey", []byte("testvalue")); err != nil {
		t.Skipf("Skipping: filesystem does not support xattr: %v", err)
	}

	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Test set Xattr", args{testFile}, "testvalue", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetXattr(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetXattr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got["user.testkey"], tt.want) {
				t.Errorf("GetXattr() got = %v, want %v", got, tt.want)
			}
		})
	}
}
