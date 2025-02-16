// Package spargo is a Wrapper for the generic spargo package:
//
//   - github.com/ross-spencer/spargo/pkg/spargo
//
// The package exists to enable to inclusion of Wikibase provenance in
// those results. Where spargo is a generic package this version is
// specific to Wikidata implementations on-top of Wikibase.
package spargo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/ross-spencer/spargo/pkg/spargo"
	"github.com/ross-spencer/wikiprov/pkg/wikiprov"
)

// DefaultAgent as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
const DefaultAgent = spargo.DefaultAgent

// Binding as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type Binding = spargo.Binding

// Item as it exists in the spargo package exported to enable dropping
// this package into host packages/executables.
type Item = spargo.Item

// SPARQLClient as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type SPARQLClient = spargo.SPARQLClient

// SPARQLResult as it exists in the spargo package exported to enable
// dropping this package into host packages/executables.
type SPARQLResult = spargo.SPARQLResult

// WikiProv wraps spargo's standard results so that we can attach
// provenance without attempting to modify the generic capabilities of
// the wikiprov's sister package.
type WikiProv struct {
	Head       map[string]interface{} `json:"head"`
	Binding    `json:"results"`
	Provenance []wikiprov.Provenance `json:"provenance,omitempty"`
}

// maxChannels determines the number of channels to use in requests to
// Wikibase for its provenance data. Ostensibly it's a throttle.
// Wikidata will return an error if we ask for too much too quickly, 20
// caused an error previously for over 1000 records. 10 seems to work
// fairly well. Without requesting this information in threads,
// processing can be pretty slow.
var maxChannels = 10

// fixKey makes it a little easier to work with the library by making
// sure that if a parameter is specified with '?' in it, the results
// will still be returned to the caller if a matching key is found.
// fixkey will strip the leading '?'.
func fixKey(key string) string {
	if key[0] == '?' {
		return key[1:]
	}
	return key
}

// SPARQLWithProv is used to query the Wikidata query service and attach
// Wikibase provenance. History can be configured as well as the number
// of threads used to connect to Wikibase. The key provided this
// function must exist as a parameter in the SPARQL query, e.g. SELECT
// `?uri` where `?uri` is the key. This parameter must also be a
// Wikidata IRI from which the QID will be returned. The QID is then
// used to grab the provenance information for the record. If key is
// empty then provenance functions will not be called.
func SPARQLWithProv(
	endpoint string,
	queryString string,
	param string,
	lenHistory int,
	threads int,
) (WikiProv, error) {
	sparqlMe := SPARQLClient{}
	sparqlMe.ClientInit(endpoint, queryString)
	res, err := sparqlMe.SPARQLGo()
	if err != nil {
		return WikiProv{}, err
	}
	provResults := WikiProv{}
	provResults.Head = res.Head
	provResults.Binding = res.Results
	if param == "" || lenHistory < 1 {
		return provResults, nil
	}
	param = fixKey(param)
	if threads > maxChannels {
		threads = maxChannels
	}
	err = provResults.attachProvenance(param, lenHistory, threads)
	if err != nil {
		return WikiProv{}, err
	}
	return provResults, nil
}

// validateIRI will attempt to perform some basic validation on IRI's
// we're trying to retrieve provenance information for. We need to build
// up a set of rules.
//
// NB. There was a URL validation rule here, e.g. did the value from
// the ?param have a sensible URL that would then be queried for
// provenance results. The configuration needed for this downstream
// starts to get complicated when you consider downstream has to work
// with a query service (SPARQL endpoint) and a Wikibase base URL that
// might both need to be configured.
func validateIRI(iri string) bool {
	const statement string = "statement"
	if strings.Contains(iri, statement) {
		return false
	}
	return true
}

// testHTTP is a rudimentary test for a URI that we can query.
func testHTTP(iri string) bool {
	if !strings.HasPrefix(iri, "http") {
		return false
	}
	return true
}

// getQID will retrieve the QID from a Wikidata IRI. It can handle
// Properties which require a special suffice, and entities which are
// our standard QIDs e.g. Q12345.
func getQID(iri string) (string, error) {
	const prop string = "prop"
	const property string = "Property"
	if !testHTTP(iri) {
		return "", nil
	}
	parsedIRI, err := url.Parse(iri)
	if err != nil {
		return "", err
	}
	qid := path.Base(parsedIRI.Path)
	if strings.Contains(iri, prop) {
		return fmt.Sprintf("%s:%s", property, qid), nil
	}
	return qid, nil
}

// ErrProvAttach provides a method for the caller to quickly anticipate
// errors in the provenance results they may want to investigate. If
// there are no errors, then all is in ordnung.
var ErrProvAttach error = fmt.Errorf("warning: there were errors retrieving provenance from Wikibase API")

// AttachProvenance will attach WikiBase provenance to SPARQL results
// from Wikidata.
func (sparql *WikiProv) attachProvenance(
	sparqlParam string,
	lenHistory int,
	threads int,
) error {
	var qids map[string]bool
	qids = make(map[string]bool)
	for _, value := range sparql.Bindings {
		wikidataIRI := value[sparqlParam].Value
		if !validateIRI(wikidataIRI) {
			continue
		}
		qid, err := getQID(wikidataIRI)
		if err != nil {
			return err
		}
		qids[qid] = false
	}
	if len(qids) < 1 {
		return fmt.Errorf("No results returned from given sparqlParam: %s", sparqlParam)
	}

	var uniqueQIDs []string
	for sparqlParam := range qids {
		uniqueQIDs = append(uniqueQIDs, sparqlParam)
	}

	preProvCache := getProvThreaded(uniqueQIDs, lenHistory, threads)
	provCache := []wikiprov.Provenance{}

	for _, value := range preProvCache {
		if value.Error != nil {
			continue
		}
		if value.Title == "" && value.Revision == 0 && value.Permalink == "" {
			continue
		}
		provCache = append(provCache, value)
	}
	if len(provCache) == 0 && lenHistory > 0 {
		return fmt.Errorf(
			"history configured but unable to retrieve history from Wikibase",
		)
	}

	sparql.Provenance = provCache

	// Check for errors so that we can warn the caller.
	for _, prov := range provCache {
		if prov.Error != nil {
			return fmt.Errorf("%w: %s", ErrProvAttach, prov.Error)
		}
	}

	return nil
}

// getProvThreaded goes out to Wikibase and collects the provenance
// associated with a record. The function takes an argument that limits
// the number of channels to be used to do work to provide some level
// of throttling and to also increase performance of this. For ~5000
// records this can take 15 minutes without concurrency.
func getProvThreaded(qids []string, lenHistory int, maxChan int) []wikiprov.Provenance {
	ch := make(chan wikiprov.Provenance)
	var mutex sync.Mutex
	counter := 0
	for channels := 0; channels < maxChan; channels++ {
		go func(ch chan wikiprov.Provenance, mutex *sync.Mutex) {
			for {
				mutex.Lock()
				idx := counter
				counter++
				mutex.Unlock()
				if counter > len(qids) {
					// Finished processing, exit.
					return
				}
				qid := qids[idx]
				// Retrieve the provenance information from Wikibase.
				prov := getProvenance(qid, lenHistory)
				ch <- prov
			}
		}(ch, &mutex)
	}
	var provCache []wikiprov.Provenance
	provCache = make([]wikiprov.Provenance, len(qids))
	getData(ch, provCache)
	return provCache
}

// getData invokes the go routines and then adds the results to the
// provenance array.
func getData(ch <-chan wikiprov.Provenance, provCache []wikiprov.Provenance) {
	for idx := 0; idx < len(provCache); idx++ {
		provCache[idx] = <-ch
	}
}

// getProvenance is a helper which is used to call wikiprov's primary
// function collecting provenance for a record from the underlying
// Wikibase implementation.
func getProvenance(qid string, lenHistory int) wikiprov.Provenance {
	prov, err := wikiprov.GetWikidataProvenance(qid, lenHistory)
	if err != nil {
		// We'll handle the error upstream.
		prov.Error = err
	}
	return prov
}

// String will return a summary of a Wikiprov structure as JSON.
func (sparql WikiProv) String() string {
	str, err := json.MarshalIndent(sparql, "", "  ")
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
