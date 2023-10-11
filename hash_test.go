package filetrove

import (
	"reflect"
	"testing"
)

func TestHashit(t *testing.T) {
	type args struct {
		inFile  string
		hashalg string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"md5", args{"testdata/emptyfile.txt", "md5"}, []byte{212, 29, 140, 217, 143, 0, 178, 4, 233, 128, 9, 152, 236, 248, 66, 126}, false},
		{"sha1", args{"testdata/emptyfile.txt", "sha1"}, []byte{218, 57, 163, 238, 94, 107, 75, 13, 50, 85, 191, 239, 149, 96, 24, 144, 175, 216, 7, 9}, false},
		{"sha256", args{"testdata/emptyfile.txt", "sha256"}, []byte{227, 176, 196, 66, 152, 252, 28, 20, 154, 251, 244, 200, 153, 111, 185, 36, 39, 174, 65, 228, 100, 155, 147, 76, 164, 149, 153, 27, 120, 82, 184, 85}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Hashit(tt.args.inFile, tt.args.hashalg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Hashit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hashit() got = %v, want %v", got, tt.want)
			}
		})
	}
}
