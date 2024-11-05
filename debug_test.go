package filetrove

import (
	"os"
	"reflect"
	"testing"
)

func TestCreateDebugPackage(t *testing.T) {
	tests := []struct {
		name    string
		want    os.File
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DebugCreateDebugPackage()
			if (err != nil) != tt.wantErr {
				t.Errorf("DebugCreateDebugPackage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DebugCreateDebugPackage() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHostinformation(t *testing.T) {
	type args struct {
		fd os.File
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DebugHostinformation(tt.args.fd); (err != nil) != tt.wantErr {
				t.Errorf("DebugHostinformation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
