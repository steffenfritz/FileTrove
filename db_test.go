package filetrove

import (
	"os"
	"testing"
)

func TestCreateFileTroveDB(t *testing.T) {
	type args struct {
		dbpath   string
		version  string
		initdate string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Create database", args{dbpath: ".", version: "TEST", initdate: "23.05.1949"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateFileTroveDB(tt.args.dbpath, tt.args.version, tt.args.initdate); (err != nil) != tt.wantErr {
				t.Errorf("CreateFileTroveDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	// Remove file after testing
	err := os.Remove("filetrove.db")
	if err != nil {
		println("Could not remove test file: filetrove.db")
	}
}
