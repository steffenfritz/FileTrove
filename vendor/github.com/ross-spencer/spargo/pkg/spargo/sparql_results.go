package spargo

import (
	"encoding/json"
	"fmt"
)

/*
Basic SPARQL results will be returned as follows:

	{
	   "head":{
	      "vars":[
	         "item",
	         "itemLabel"
	      ]
	   },
	   "results":{
	      "bindings":[
	         {
	            "predicate_one":{
	               "type":"uri",
	               "value":"http://www.wikidata.org/entity/Q28114535"
	            },
	            "predicate_two":{
	               "xml:lang":"en",
	               "type":"literal",
	               "value":"Mr. White"
	            }
	         },
	         {
	            "predicate_one":{
	               "type":"uri",
	               "value":"http://www.wikidata.org/entity/Q28665865"
	            },
	            "predicate_two":{
	               "xml:lang":"en",
	               "type":"literal",
	               "value":"Мyka"
	            }
	         }
	      ]
	   }
	}
*/

// Item describes the verbose output of a SPARQL query needed to contextualize
// it fully.
type Item struct {
	Lang     string `json:"xml:lang,omitempty"` // Populated if requested in query.
	Type     string `json:"type"`
	Value    string `json:"value"`
	DataType string `json:"datatype,omitempty"`
}

// Binding is made up of multiple Items we can access those here.
type Binding struct {
	Bindings []map[string]Item `json:"bindings"`
}

// SPARQLResult packages a SPARQL response from an endpoint.
type SPARQLResult struct {
	Head    map[string]interface{} `json:"head"`
	Results Binding                `json:"results"`
}

// String will return a string representation of SPARQLResult.
func (sparql SPARQLResult) String() string {
	str, err := json.MarshalIndent(sparql, "", "  ")
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s", str)
}
