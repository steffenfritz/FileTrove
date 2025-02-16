package wikiprov

// Consts and variables used internally to request or create the data
// that we want.

import (
	"fmt"
	"runtime/debug"
	"strings"
)

const contactInfo string = "(https://github.com/ross-spencer/wikiprov; all.along.the.watchtower+github@gmail.com)"

var agent string = fmt.Sprintf("wikiprov/0.0.0 %s", contactInfo)

// TODO: this probably needs to be overwritten for another
// Wikibases/Wikidata instance.
const defaultBaseURI = "https://www.wikidata.org/"
const indexPage = "w/index.php"
const apiPage = "w/api.php"

var wikibaseAPI = ""
var wikibasePermalinkBase = ""

var format = "json"
var action = "query"
var prop = "revisions"

var revisionPropertiesDefault = [...]string{"ids", "user", "comment", "timestamp", "sha1"}

func init() {
	agent = getVersionFromBuildFlags()
	wikibaseAPI = constructWikibaseAPIURL(defaultBaseURI)
	wikibasePermalinkBase = constructWikibaseIndexURL(defaultBaseURI)
}

// getVersionFromBuildFlags returns a version string from the build
// flags if the build flags have been set.
func getVersionFromBuildFlags() string {
	var buildInfo *debug.BuildInfo
	var ok bool
	if buildInfo, ok = debug.ReadBuildInfo(); !ok {
		return agent
	}
	for _, prop := range buildInfo.Settings {
		if prop.Key != "-ldflags" {
			continue
		}
		setting := strings.Split(prop.Value, "-X")
		var version string
		for _, property := range setting {
			if strings.Contains(property, "main.version") {
				version = strings.TrimSpace(strings.Split(property, "=")[1])
			}
		}
		agent = fmt.Sprintf(
			"wikiprov/%s %s",
			version,
			contactInfo,
		)
	}
	return agent
}

// constructWikibaseAPIURL will create a URL for connecting to the
// Wikimedia API.
func constructWikibaseAPIURL(baseURL string) string {
	if strings.HasSuffix(baseURL, "/") {
		return fmt.Sprintf("%s%s", baseURL, apiPage)
	}
	return fmt.Sprintf("%s/%s", baseURL, apiPage)
}

// constructWikibaseIndexURL will create a URL pointing to the index
// page of the given Wikimedia/Wikibase instance, e.g. for resolution
// of permalinks.
func constructWikibaseIndexURL(baseURL string) string {
	if strings.HasSuffix(baseURL, "/") {
		return fmt.Sprintf("%s%s", baseURL, indexPage)
	}
	return fmt.Sprintf("%s/%s", baseURL, indexPage)
}
