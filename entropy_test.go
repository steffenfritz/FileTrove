package filetrove

import "testing"

func TestEntropy(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		wantEntropy float64
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEntropy, err := Entropy(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Entropy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotEntropy != tt.wantEntropy {
				t.Errorf("Entropy() gotEntropy = %v, want %v", gotEntropy, tt.wantEntropy)
			}
		})
	}
}
