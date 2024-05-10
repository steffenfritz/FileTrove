package filetrove

import (
	"reflect"
	"testing"
)

func TestReadDC(t *testing.T) {
	type args struct {
		dcjson string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Title", args{dcjson: "testdata/dublincore_ex.json"}, "The Dublin Core™ Resource Type (DC.Type) element is used to describe the category or genre of the content of the resource.", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadDC(tt.args.dcjson)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadDC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Type, tt.want) {
				t.Errorf("ReadDC() got = %v, want %v", got.Title, tt.want)
			}
		})
	}
}
