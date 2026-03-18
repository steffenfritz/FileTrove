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
		name        string
		args        args
		want        []string
		want1       []string
		wantSkipped []string
		wantErr     bool
	}{
		{
			name: "testdata file list",
			args: args{"testdata"},
			want: []string{
				"testdata/.hiddendir/.gitkeep",
				"testdata/.hiddenfile",
				"testdata/directory/A_PDF_File.pdf",
				"testdata/dublincore_ex.json",
				"testdata/emptyfile.txt",
				"testdata/images/screenshot_1.jpeg",
				"testdata/images/screenshot_1.jpeg_original",
				"testdata/images/screenshot_1.png",
				"testdata/images/screenshot_1.png_original",
				"testdata/images/screenshot_1.tiff",
				"testdata/images/screenshot_1.tiff_original",
				"testdata/noextension",
				"testdata/textfile.txt",
				"testdata/transparent.png",
				"testdata/white.jpg",
				"testdata/yara/testrule.yara",
			},
			want1: []string{
				"testdata",
				"testdata/.hiddendir",
				"testdata/directory",
				"testdata/images",
				"testdata/yara",
			},
			wantSkipped: nil,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, gotSkipped, err := CreateFileList(tt.args.rootDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFileList() files = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("CreateFileList() dirs = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(gotSkipped, tt.wantSkipped) {
				t.Errorf("CreateFileList() skipped = %v, want %v", gotSkipped, tt.wantSkipped)
			}
		})
	}
}
