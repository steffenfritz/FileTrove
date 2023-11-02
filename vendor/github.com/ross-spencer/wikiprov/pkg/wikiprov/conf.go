package wikiprov

// Consts and variables used internally to request or create the data
// that we want.

import (
	"fmt"
	"strings"
)

const agent string = "wikiprov/0.2.0 (https://github.com/ross-spencer/wikiprov/; all.along.the.watchtower+github@gmail.com)"

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
	wikibaseAPI = constructWikibaseAPIURL(defaultBaseURI)
	wikibasePermalinkBase = constructWikibaseIndexURL(defaultBaseURI)
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
