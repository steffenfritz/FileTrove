package wikiprov

// testJSON provides us with an example output from Wikibase to work
// with. There's some nice data here as well as some fields that we're
// not using yet but might be interesting to look at in time.
//
// From: https://www.wikidata.org/wiki/Q12345 (Count von Count)
//
const testJSON string = `{
    "continue": {
        "continue": "||",
        "rvcontinue": "20210104032931|1334519319"
    },
    "query": {
        "normalized": [
            {
                "from": "q12345",
                "to": "Q12345"
            }
        ],
        "pages": {
            "13925": {
                "ns": 0,
                "pageid": 13925,
                "revisions": [
                    {
                        "comment": "/* wbsetsitelink-add-both:2|dewiki */ Graf Zahl, [[Q70894304]]",
                        "parentid": 1419073806,
                        "revid": 1419131078,
                        "timestamp": "2021-05-11T20:17:31Z",
                        "user": "user1"
                    },
                    {
                        "comment": "/* wbsetclaim-create:2||1 */ [[Property:P97]]: [[Q3519259]]",
                        "parentid": 1419073622,
                        "revid": 1419073806,
                        "timestamp": "2021-05-11T16:53:18Z",
                        "user": "user2"
                    },
                    {
                        "comment": "/* wbsetclaim-create:2||1 */ [[Property:P5247]]: 3005-34149",
                        "parentid": 1419064895,
                        "revid": 1419073622,
                        "timestamp": "2021-05-11T16:52:45Z",
                        "user": "user3"
                    },
                    {
                        "comment": "/* wbcreateclaim-create:1| */ [[Property:P8345]]: [[Q106804572]], #quickstatements; #temporary_batch_1620750589351",
                        "parentid": 1393551702,
                        "revid": 1419064895,
                        "timestamp": "2021-05-11T16:16:09Z",
                        "user": "user4"
                    },
                    {
                        "comment": "/* wbsetaliases-add:3|ru */ \u0413\u0440\u0430\u0444 \u0417\u043d\u0430\u043a, Count von Count, The Count",
                        "parentid": 1393551693,
                        "revid": 1393551702,
                        "timestamp": "2021-03-31T10:27:19Z",
                        "user": "user5"
                    }
                ],
                "title": "Q12345"
            }
        }
    },
    "warnings": {
        "main": {
            "*": "Unrecognized parameter: rvprops."
        }
    }
}`
