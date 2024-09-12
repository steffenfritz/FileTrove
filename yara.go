package filetrove

import (
	yarax "github.com/VirusTotal/yara-x/go"
	"os"
)

// YaraCompile compiles a string that is provided via a flag from the main function
func YaraCompile(ruleFile string) (*yarax.Rules, error) {
	var rules *yarax.Rules

	ruleBytes, err := os.ReadFile(ruleFile)
	if err != nil {
		return nil, err
	}
	rules, err = yarax.Compile(string(ruleBytes))

	return rules, err
}

// YaraScan receives pre-compiled rules and checks if one or more rules match on the input file
// For that check it has to read files into []byte. While YARA itself is fast this might become a bottleneck.
func YaraScan(rules *yarax.Rules, inFile string) (*yarax.ScanResults, error) {
	fileRead, err := os.ReadFile(inFile)
	if err != nil {
		return nil, err
	}

	matchedRules, err := rules.Scan(fileRead)

	return matchedRules, err
}
