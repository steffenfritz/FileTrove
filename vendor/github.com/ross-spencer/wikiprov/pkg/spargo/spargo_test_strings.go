package spargo

// errorTest and errorTests provide a simple test layout for when there
// is no connection to the server at any point in time.
//
// Threading provides the greatest opportunity for error in this code?
// we try out different permutations here.
//
type errorTest struct {
	qids    []string
	threads int
}

var errorTests = []errorTest{
	{[]string{"1"}, 1},
	{[]string{"1"}, 5},
	{[]string{"1"}, 7},
	{[]string{"1"}, 10},
	{[]string{"1"}, 100},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 1},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 5},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 7},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 10},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 100},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 1},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 5},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 7},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 10},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 100},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 1},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 2},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 3},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 4},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 5},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 6},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 7},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 8},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 9},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 10},
}

// threadTest and threadTests provide a simple test layout to ensure
// that without error and with multiple combinations of concurrency
// provenance can be retrieved from a test server reliably.
type threadTest struct {
	qids    []string
	threads int
	result  string
}

var threadTests = []threadTest{
	{[]string{"1"}, 1, threadedProvenance},
	{[]string{"1"}, 5, threadedProvenance},
	{[]string{"1"}, 7, threadedProvenance},
	{[]string{"1"}, 10, threadedProvenance},
	{[]string{"1"}, 100, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 1, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 5, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 7, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 10, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}, 100, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 1, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 5, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 7, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 10, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}, 100, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 1, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 2, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 3, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 4, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 5, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 6, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 7, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 8, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 9, threadedProvenance},
	{[]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F"}, 10, threadedProvenance},
}

// threadedProvenance is an example JSON string that will be used to
// measure the correct return of the provenance retrieval functions over
// many different configurations of concurrency.
var threadedProvenance string = `{
        "continue": {
            "rvcontinue": "20200221144033|1120067133",
            "continue": "||"
        },
        "query": {
            "pages": {
                "5147078": {
                    "pageid": 5147078,
                    "ns": 0,
                    "title": "Q12345",
                    "revisions": [{
                            "revid": 2600,
                            "parentid": 1247208427,
                            "user": "Emmanuel Goldstein",
                            "timestamp": "2020-08-31T23:13:00Z",
                            "sha1": "4fa4f3344e2db600c11273028e63ba21976ede80",
                            "comment": "edit comment #1"
                        },
                        {
                            "revid": 1000,
                            "parentid": 1120067133,
                            "user": "Robert Smith",
                            "timestamp": "2020-08-01T23:13:00Z",
                            "sha1": "88a134dc3b112584e143003cadf0fdf3a4503dfe",
                            "comment": "edit comment #2"
                        }
                    ]
                }
            }
        }
    }`

// wikidataResultsJSON is a small extract of results from the Wikidata
// query service that we would like to attach provenance to.
var wikidataResultsJSON = `{
    "head": {
        "vars": [
            "uri",
            "uriLabel",
            "puid",
            "extension",
            "mimetype",
            "encodingLabel",
            "referenceLabel",
            "date",
            "relativityLabel",
            "offset",
            "sig"
        ]
    },
    "results": {
        "bindings": [
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1377"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100135637"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.1.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1378"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100136218"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.2.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1379"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100136955"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.3.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1380"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100136960"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.4.0",
                    "xml:lang": "en"
                }
            },
            {
                "encodingLabel": {
                    "type": "literal",
                    "value": "hexadecimal",
                    "xml:lang": "en"
                },
                "extension": {
                    "type": "literal",
                    "value": "dwb"
                },
                "mimetype": {
                    "type": "literal",
                    "value": "application/octet-stream"
                },
                "offset": {
                    "datatype": "http://www.w3.org/2001/XMLSchema#decimal",
                    "type": "literal",
                    "value": "1"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1381"
                },
                "relativityLabel": {
                    "type": "literal",
                    "value": "beginning of file",
                    "xml:lang": "en"
                },
                "sig": {
                    "type": "literal",
                    "value": "870100"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100137240"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "VariCAD Drawing",
                    "xml:lang": "en"
                }
            },
            {
                "encodingLabel": {
                    "type": "literal",
                    "value": "hexadecimal",
                    "xml:lang": "en"
                },
                "extension": {
                    "type": "literal",
                    "value": "dwb"
                },
                "mimetype": {
                    "type": "literal",
                    "value": "application/octet-stream"
                },
                "offset": {
                    "datatype": "http://www.w3.org/2001/XMLSchema#decimal",
                    "type": "literal",
                    "value": "0"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1381"
                },
                "relativityLabel": {
                    "type": "literal",
                    "value": "beginning of file",
                    "xml:lang": "en"
                },
                "sig": {
                    "type": "literal",
                    "value": "870100"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100137240"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "VariCAD Drawing",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xpdz"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1385"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100151671"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "Bruker PDZ",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "pdz"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1385"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://www.wikidata.org/entity/Q100151671"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "Bruker PDZ",
                    "xml:lang": "en"
                }
            }
        ]
    }
}
`

// wikidataResultsJSONExampleDotCom provides a set of results that mimic
// a change in Base URL so that the library can work on the widest
// possible number of implementations.
var wikidataResultsJSONExampleDotCom = `{
    "head": {
        "vars": [
            "uri",
            "uriLabel",
            "puid",
            "extension",
            "mimetype",
            "encodingLabel",
            "referenceLabel",
            "date",
            "relativityLabel",
            "offset",
            "sig"
        ]
    },
    "results": {
        "bindings": [
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1377"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100135637"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.1.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1378"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100136218"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.2.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1379"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100136955"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.3.0",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xml"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1380"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100136960"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "XDOMEA 2.4.0",
                    "xml:lang": "en"
                }
            },
            {
                "encodingLabel": {
                    "type": "literal",
                    "value": "hexadecimal",
                    "xml:lang": "en"
                },
                "extension": {
                    "type": "literal",
                    "value": "dwb"
                },
                "mimetype": {
                    "type": "literal",
                    "value": "application/octet-stream"
                },
                "offset": {
                    "datatype": "http://www.w3.org/2001/XMLSchema#decimal",
                    "type": "literal",
                    "value": "1"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1381"
                },
                "relativityLabel": {
                    "type": "literal",
                    "value": "beginning of file",
                    "xml:lang": "en"
                },
                "sig": {
                    "type": "literal",
                    "value": "870100"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100137240"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "VariCAD Drawing",
                    "xml:lang": "en"
                }
            },
            {
                "encodingLabel": {
                    "type": "literal",
                    "value": "hexadecimal",
                    "xml:lang": "en"
                },
                "extension": {
                    "type": "literal",
                    "value": "dwb"
                },
                "mimetype": {
                    "type": "literal",
                    "value": "application/octet-stream"
                },
                "offset": {
                    "datatype": "http://www.w3.org/2001/XMLSchema#decimal",
                    "type": "literal",
                    "value": "0"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1381"
                },
                "relativityLabel": {
                    "type": "literal",
                    "value": "beginning of file",
                    "xml:lang": "en"
                },
                "sig": {
                    "type": "literal",
                    "value": "870100"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100137240"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "VariCAD Drawing",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "xpdz"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1385"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100151671"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "Bruker PDZ",
                    "xml:lang": "en"
                }
            },
            {
                "extension": {
                    "type": "literal",
                    "value": "pdz"
                },
                "puid": {
                    "type": "literal",
                    "value": "fmt/1385"
                },
                "uri": {
                    "type": "uri",
                    "value": "http://example.com/entity/Q100151671"
                },
                "uriLabel": {
                    "type": "literal",
                    "value": "Bruker PDZ",
                    "xml:lang": "en"
                }
            }
        ]
    }
}
`

// attachedProvenance is a test JSON struct that will be converted to
// wikiprov.Provenance from its native structure below.
var attachedProvenance string = `{
        "continue": {
            "rvcontinue": "20200221144033|1120067133",
            "continue": "||"
        },
        "query": {
            "pages": {
                "5147078": {
                    "pageid": 5147078,
                    "ns": 0,
                    "title": "Q12345",
                    "revisions": [{
                            "revid": 2600,
                            "parentid": 1247208427,
                            "user": "Emmanuel Goldstein",
                            "timestamp": "2020-08-31T23:13:00Z",
                            "sha1": "4fa4f3344e2db600c11273028e63ba21976ede80",
                            "comment": "edit comment #1"
                        },
                        {
                            "revid": 1000,
                            "parentid": 1120067133,
                            "user": "Robert Smith",
                            "timestamp": "2020-08-01T23:13:00Z",
                            "sha1": "88a134dc3b112584e143003cadf0fdf3a4503dfe",
                            "comment": "edit comment #2"
                        }
                    ]
                }
            }
        }
    }`
