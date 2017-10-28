# A simple search engine service

## Quick start

```bash
$ git clone https://github.com/Lupino/yuntan-search.git
$ cd yuntan-search
$ docker build -t yuntan-search .
$ docker run -d --name yuntan-search -p 8095:8095 yuntan-search
```


## Create index

```bash
$ export INDEX_NAME=articles
$ curl --data-binary @data/mapping.json -H 'Content-Type: application/json' http://127.0.0.1:8095/api/$INDEX_NAME
{"status":"ok"}
```

## Index document

```bash
$ export DOC_ID=test
$ curl -XPUT -H 'Content-Type: application/json' --data-binary @data/doc.json http://127.0.0.1:8095/api/$INDEX_NAME/$DOC_ID
{"status":"ok"}
```

## Search document

```bash
$ curl -XPOST -H 'Content-Type: application/json' --data-binary @data/search.json http://127.0.0.1:8095/api/$INDEX_NAME/_search | json_reformat
{
    "status": {
        "total": 1,
        "failed": 0,
        "successful": 1
    },
    "request": {
        "query": {
            "query": "你好"
        },
        "size": 10,
        "from": 0,
        "highlight": null,
        "fields": null,
        "facets": null,
        "explain": false,
        "sort": [
            "-_score"
        ],
        "includeLocations": false
    },
    "hits": [
        {
            "index": "articles",
            "id": "test",
            "score": 0.21697770945227396,
            "sort": [
                "_score"
            ]
        }
    ],
    "total_hits": 1,
    "max_score": 0.21697770945227396,
    "took": 1199353,
    "facets": {

    }
}
```
