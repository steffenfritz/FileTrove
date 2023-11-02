package wikiprov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// structs for wikiprov
//
// The primary JSON data we're interested in is the info struct from
// the API endpoint.
//
/*
	{
		"continue": {
			"rvcontinue": "20200221144033|1120067133",
			"continue": "||"
		},
		"query": {
			"pages": {
				"5147078": {
					"pageid": 5147078,
					"ns": 0,
					"title": "Q5381415",
					"revisions": [{
							"revid": 1247209137,
							"parentid": 1247208427,
							"user": "Beet keeper",
							"timestamp": "2020-08-04T23:41:27Z",
							"sha1": "4fa4f3344e2db600c11273028e63ba21976ede80",
							"comment": ... wbsetclaim-update:2||1 ... [[Property:P4152]]: B297E169"
						},
						{
							"revid": 1247208427,
							"parentid": 1120067133,
							"user": "Beet keeper",
							"timestamp": "2020-08-04T23:40:10Z",
							"sha1": "88a134dc3b112584e143003cadf0fdf3a4503dfe",
							"comment": ... wbsetclaim-update:2||1 ... [[Property:P4152]]: 325E1010"
						}
					]
				}
			}
		}
	}
*/
//

type revision struct {
	RevisionID int    `json:"revid"`
	ParentID   int    `json:"parentid"`
	User       string `json:"user"`
	Timestamp  string `json:"timestamp"`
	SHA1       string `json:"sha1"`
	Comment    string `json:"comment"`
}

// String creates a simple rendition of the revision history until we
// know what else we want to do with it.
func (rev revision) String() string {
	return fmt.Sprintf("%s (oldid: %d): '%s' edited: '%s'", rev.Timestamp, rev.RevisionID, rev.User, rev.Comment)
}

type revisions struct {
	PageID    int    `json:"pageid"`
	NS        int    `json:"ns"`
	Title     string `json:"title"`
	Revisions []revision
}

type page map[string]revisions

type pages struct {
	Pages page `json:"pages"`
}

type wdRevisions struct {
	Query pages `json:"query"`
}

// normalize simplifies the wdInfo structure so it can be easily used by
// the caller.
//
// {
//  	"Title": "Q27229608",
//  	"Revision": 784082439,
//  	"Modified": "2018-11-07T16:26:11Z",
//  	"Permalink": "https://www.wikidata.org/w/index.php?format=json&oldid=0&title="
//      "History": [ ... ]
// }
//
func (revisions *wdRevisions) normalize() Provenance {

	var prov Provenance

	revMap := revisions.Query.Pages

	var key string
	for k := range revMap {
		key = k
		break
	}

	revs := revMap[key]
	if len(revs.Revisions) < 1 {
		return Provenance{}
	}

	firstRecord := revs.Revisions[0]

	prov.Title = revs.Title
	prov.Revision = firstRecord.RevisionID
	prov.Modified = firstRecord.Timestamp
	prov.Permalink = prov.buildPermalink()

	for _, value := range revs.Revisions {
		prov.History = append(prov.History, fmt.Sprintf("%s", value))
	}

	return prov
}

// Provenance provides simplified provenance information about a
// Wikidata record.
type Provenance struct {
	Title     string
	Revision  int
	Modified  string
	Permalink string
	History   []string
	Error     error `json:"-"`
}

// buildPermalink creates a permalink based on the title and revision
// values being set in the Provenance structure.
func (prov *Provenance) buildPermalink() string {
	const paramTitle = "title"
	const paramOldID = "oldid"
	req, _ := http.NewRequest("GET", wikibasePermalinkBase, nil)
	query := req.URL.Query()
	title := prov.Title
	oldid := prov.Revision
	query.Set(paramTitle, title)
	query.Set(paramOldID, fmt.Sprintf("%d", oldid))
	req.URL.RawQuery = query.Encode()
	return fmt.Sprintf("%s", req.URL)
}

// String creates a human readable representation of the provenance
// struct.
func (prov Provenance) String() string {

	str, err := json.MarshalIndent(prov, "", "  ")
	if err != nil {
		return ""
	}

	// THe encoder now escapes these values, this is for browser
	// compatibility, and I don't think it matters to us too much.
	//
	//    * https://stackoverflow.com/a/24657016
	//
	str = bytes.Replace(str, []byte("\\u003c"), []byte("<"), -1)
	str = bytes.Replace(str, []byte("\\u003e"), []byte(">"), -1)
	str = bytes.Replace(str, []byte("\\u0026"), []byte("&"), -1)

	return fmt.Sprintf("%s", str)
}
