package filetrove

import (
	"reflect"
	"testing"
)

func TestExifDecode(t *testing.T) {
	type args struct {
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		want    ExifParsed
		wantErr bool
	}{
		{"EXIF screenshot_1.jpg", args{"testdata/images/screenshot_1.jpeg"}, ExifParsed{Artist: "\"Steffen Fritz\""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExifDecode(tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExifDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Artist, tt.want.Artist) {
				t.Errorf("ExifDecode() got = %v, want %v", got.Artist, tt.want.Artist)
			}
		})
	}
}
