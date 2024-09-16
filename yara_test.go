package filetrove

import (
	yarax "github.com/VirusTotal/yara-x/go"
	"os"
	"reflect"
	"testing"
)

// TestYaraCompile checks if a rule from a text file can be compiled, so
// we just check for the wanted error, which is false.
func TestYaraCompile(t *testing.T) {
	type args struct {
		ruleString string
	}
	tests := []struct {
		name string
		args args
		// want    *yarax.Rules
		wantErr bool
	}{
		{"Successful compile", args{"testdata/yara/testrule.yara"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//got, err := YaraCompile(tt.args.ruleString)
			_, err := YaraCompile(tt.args.ruleString)
			if (err != nil) != tt.wantErr {
				t.Errorf("YaraCompile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			/*if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("YaraCompile() got = %v, want %v", got, tt.want)
			}*/
		})
	}
}

func TestYaraScan(t *testing.T) {
	type args struct {
		rules  *yarax.Rules
		inFile string
	}

	// Compile rules for match testing.
	rules, err := YaraCompile("testdata/yara/testrule.yara")
	if err != nil {
		println("yara compilation error:", err)
		os.Exit(1)
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Yara Matching Test Rule", args{rules, "testdata/transparent.png"}, "TestPNGRule", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := YaraScan(tt.args.rules, tt.args.inFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("YaraScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			matchedRule := got.MatchingRules()[0]

			if !reflect.DeepEqual(matchedRule.Identifier(), tt.want) {
				t.Errorf("YaraScan() got = %v, want %v", got, tt.want)
			}
		})
	}
}
