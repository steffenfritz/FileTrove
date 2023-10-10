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
		// TODO: Add test cases.
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
