package filetrove

import (
	"reflect"
	"testing"
)

func TestCreateFileList(t *testing.T) {
	type args struct {
		rootDir string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		want1   []string
		wantErr bool
	}{
		{"testdata file list", args{"testdata"}, []string{"testdata/.hiddenfile",
			"testdata/directory/A_PDF_File.pdf", "testdata/dublincore_ex.json", "testdata/emptyfile.txt", "testdata/images/screenshot_1.jpeg", "testdata/images/screenshot_1.jpeg_original", "testdata/images/screenshot_1.png", "testdata/images/screenshot_1.png_original", "testdata/images/screenshot_1.tiff", "testdata/images/screenshot_1.tiff_original", "testdata/noaccess.rtf", "testdata/noextension", "testdata/textfile.txt", "testdata/transparent.png", "testdata/white.jpg"},
			[]string{"testdata", "testdata/.hiddendir", "testdata/directory", "testdata/images"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := CreateFileList(tt.args.rootDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFileList() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateFileList() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
