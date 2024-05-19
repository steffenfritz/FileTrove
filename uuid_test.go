package filetrove

import "testing"

func TestCreateUUID(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{"UUID length", 36, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateUUID()
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.want {
				t.Errorf("CreateUUID() got = %v, want %v", len(got), tt.want)
			}
		})
	}
}
