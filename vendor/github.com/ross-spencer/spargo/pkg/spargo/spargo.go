/* Golang SPARQL package

Package spargo enables the querying of a SPARQL data store using Golang.

	"...Too rich for some people's tastes..."
*/

package spargo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// DefaultAgent user-agent determined by Wikidata User-agent policy: https://meta.wikimedia.org/wiki/User-Agent_policy.
const DefaultAgent string = "spargo/0.4.1 (https://github.com/ross-spencer/spargo/; all.along.the.watchtower+github@gmail.com)"

// DefaultAccept is the default accept-content string to be used in the HTTP request header.
const DefaultAccept string = "application/sparql-results+json, application/json"

// SPARQLClient ...
type SPARQLClient struct {
	Client  *http.Client
	BaseURL string
	Agent   string
	Accept  string
	Query   string
}

// setupClient prepares a http client to talk to a SPARQL endpoint. If
// the struct hasn't already been configured for one, then it can be
// done using this setup method.
func setupClient(endpoint *SPARQLClient) {
	if endpoint.Client == nil {
		endpoint.Client = &http.Client{}
	}
}

// SPARQLGo takes our SparqlEndpoint structure and packages that as a request
// for our SPARQL endpoint of choice. For the given
func (endpoint *SPARQLClient) SPARQLGo() (SPARQLResult, error) {

	// Make sure there is a fresh http.Client{} associated with the
	// structure for our request.
	setupClient(endpoint)

	req, err := http.NewRequest("GET", endpoint.BaseURL, nil)

	if err != nil {
		return SPARQLResult{}, err
	}

	req.Header.Add("User-Agent", endpoint.Agent)
	req.Header.Add("Accept", endpoint.Accept)

	query := req.URL.Query()
	query.Add("query", endpoint.Query)
	req.URL.RawQuery = query.Encode()

	resp, err := endpoint.Client.Do(req)
	if err != nil {
		return SPARQLResult{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseErr := ResponseError{}
		return SPARQLResult{}, responseErr.makeError(200, resp.StatusCode)

	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return SPARQLResult{}, err
	}

	var sparqlResponse SPARQLResult
	err = json.Unmarshal(body, &sparqlResponse)
	if err != nil {
		return SPARQLResult{}, err
	}

	return sparqlResponse, nil
}

// SetUserAgent agent allows the user to set a custom user agent or use the
// library's default.
func (endpoint *SPARQLClient) SetUserAgent(agent string) {
	if agent == "" {
		agent = DefaultAgent
	}
	endpoint.Agent = agent
}

// SetAcceptHeader will allow us to request results in other data formats. Our
// default is SPARQL JSON.
func (endpoint *SPARQLClient) SetAcceptHeader(accept string) {
	if accept == "" {
		accept = DefaultAccept
	}
	endpoint.Accept = accept
}

// SetQuery enables us to set the SPARQL query.
func (endpoint *SPARQLClient) SetQuery(queryString string) {
	if queryString == "" {
		// Shall we perform some error handling here?
	}
	endpoint.Query = queryString
}

// SetURL lets us set the URL of the SPARQL endpoint to query.
func (endpoint *SPARQLClient) SetURL(url string) {
	endpoint.BaseURL = url
}

// ClientInit provides us with a helper function to set endpoint URL and
// query string in a single go. Default values are set for user-agent
// and accept-content strings.
func (endpoint *SPARQLClient) ClientInit(url string, queryString string) {
	endpoint.SetURL(url)
	endpoint.SetQuery(queryString)
	// There is scope to allow these to be set by the caller. For users
	// to override these values now they would call the setter functions
	// before calling SPARQLGo().
	endpoint.SetUserAgent("")
	endpoint.SetAcceptHeader("")
}
