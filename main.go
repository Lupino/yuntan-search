package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/blevesearch/bleve"
	bleveHttp "github.com/blevesearch/bleve/http"

	// import general purpose configuration
	// _ "github.com/blevesearch/bleve/config"
	"github.com/blevesearch/bleve/index/store/goleveldb"

	// token maps
	_ "github.com/blevesearch/bleve/analysis/tokenmap"

	// fragment formatters
	_ "github.com/blevesearch/bleve/search/highlight/format/ansi"
	_ "github.com/blevesearch/bleve/search/highlight/format/html"

	// fragmenters
	_ "github.com/blevesearch/bleve/search/highlight/fragmenter/simple"

	// highlighters
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/ansi"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/html"
	_ "github.com/blevesearch/bleve/search/highlight/highlighter/simple"

	// char filters
	_ "github.com/blevesearch/bleve/analysis/char/html"
	_ "github.com/blevesearch/bleve/analysis/char/regexp"
	_ "github.com/blevesearch/bleve/analysis/char/zerowidthnonjoiner"

	// analyzers
	_ "github.com/blevesearch/bleve/analysis/analyzer/custom"
	_ "github.com/blevesearch/bleve/analysis/analyzer/keyword"
	_ "github.com/blevesearch/bleve/analysis/analyzer/simple"
	_ "github.com/blevesearch/bleve/analysis/analyzer/standard"
	_ "github.com/blevesearch/bleve/analysis/analyzer/web"

	// token filters
	_ "github.com/blevesearch/bleve/analysis/token/apostrophe"
	_ "github.com/blevesearch/bleve/analysis/token/compound"
	_ "github.com/blevesearch/bleve/analysis/token/edgengram"
	_ "github.com/blevesearch/bleve/analysis/token/elision"
	_ "github.com/blevesearch/bleve/analysis/token/keyword"
	_ "github.com/blevesearch/bleve/analysis/token/length"
	_ "github.com/blevesearch/bleve/analysis/token/lowercase"
	_ "github.com/blevesearch/bleve/analysis/token/ngram"
	_ "github.com/blevesearch/bleve/analysis/token/shingle"
	_ "github.com/blevesearch/bleve/analysis/token/stop"
	_ "github.com/blevesearch/bleve/analysis/token/truncate"
	_ "github.com/blevesearch/bleve/analysis/token/unicodenorm"

	// tokenizers
	_ "github.com/blevesearch/bleve/analysis/tokenizer/exception"
	_ "github.com/blevesearch/bleve/analysis/tokenizer/regexp"
	_ "github.com/blevesearch/bleve/analysis/tokenizer/single"
	_ "github.com/blevesearch/bleve/analysis/tokenizer/unicode"
	_ "github.com/blevesearch/bleve/analysis/tokenizer/web"
	_ "github.com/blevesearch/bleve/analysis/tokenizer/whitespace"

    "github.com/Lupino/tokenizer"

	// date time parsers
	_ "github.com/blevesearch/bleve/analysis/datetime/flexible"
	_ "github.com/blevesearch/bleve/analysis/datetime/optional"

	// languages
	// _ "github.com/blevesearch/bleve/analysis/lang/ar"
	// _ "github.com/blevesearch/bleve/analysis/lang/bg"
	// _ "github.com/blevesearch/bleve/analysis/lang/ca"
	_ "github.com/blevesearch/bleve/analysis/lang/cjk"
	// _ "github.com/blevesearch/bleve/analysis/lang/ckb"
	// _ "github.com/blevesearch/bleve/analysis/lang/cs"
	// _ "github.com/blevesearch/bleve/analysis/lang/el"
	_ "github.com/blevesearch/bleve/analysis/lang/en"
	// _ "github.com/blevesearch/bleve/analysis/lang/eu"
	// _ "github.com/blevesearch/bleve/analysis/lang/fa"
	// _ "github.com/blevesearch/bleve/analysis/lang/fr"
	// _ "github.com/blevesearch/bleve/analysis/lang/ga"
	// _ "github.com/blevesearch/bleve/analysis/lang/gl"
	// _ "github.com/blevesearch/bleve/analysis/lang/hi"
	// _ "github.com/blevesearch/bleve/analysis/lang/hy"
	// _ "github.com/blevesearch/bleve/analysis/lang/id"
	// _ "github.com/blevesearch/bleve/analysis/lang/in"
	// _ "github.com/blevesearch/bleve/analysis/lang/it"
	// _ "github.com/blevesearch/bleve/analysis/lang/pt"

	// kv stores
	// _ "github.com/blevesearch/bleve/index/store/boltdb"
	// _ "github.com/blevesearch/bleve/index/store/goleveldb"
	// _ "github.com/blevesearch/bleve/index/store/gtreap"
	// _ "github.com/blevesearch/bleve/index/store/moss"

	// index types
	_ "github.com/blevesearch/bleve/index/upsidedown"
)

var bindAddr = flag.String("addr", ":8095", "http listen address")
var dataDir = flag.String("dataDir", "data", "data directory")
var segoAddr = flag.String("segoAddr", "localhost:3000", "SegoTokenizer address")

func main() {
	flag.Parse()

    bleve.Config.DefaultKVStore = goleveldb.Name
    tokenizer.DefaultSegoTokenizerHost = *segoAddr

	// walk the data dir and register index names
	dirEntries, err := ioutil.ReadDir(*dataDir)
	if err != nil {
		log.Fatalf("error reading data dir: %v", err)
	}

	for _, dirInfo := range dirEntries {
		indexPath := *dataDir + string(os.PathSeparator) + dirInfo.Name()

		// skip single files in data dir since a valid index is a directory that
		// contains multiple files
		if !dirInfo.IsDir() {
			log.Printf("not registering %s, skipping", indexPath)
			continue
		}

		i, err := bleve.Open(indexPath)
		if err != nil {
			log.Printf("error opening index %s: %v", indexPath, err)
		} else {
			log.Printf("registered index: %s", dirInfo.Name())
			bleveHttp.RegisterIndexName(dirInfo.Name(), i)
			// set correct name in stats
			i.SetName(dirInfo.Name())
		}
	}

	router := mux.NewRouter()
	router.StrictSlash(true)

	// add the API
	createIndexHandler := bleveHttp.NewCreateIndexHandler(*dataDir)
	createIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", createIndexHandler).Methods("PUT")

	getIndexHandler := bleveHttp.NewGetIndexHandler()
	getIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", getIndexHandler).Methods("GET")

	deleteIndexHandler := bleveHttp.NewDeleteIndexHandler(*dataDir)
	deleteIndexHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}", deleteIndexHandler).Methods("DELETE")

	listIndexesHandler := bleveHttp.NewListIndexesHandler()
	router.Handle("/api", listIndexesHandler).Methods("GET")

	docIndexHandler := bleveHttp.NewDocIndexHandler("")
	docIndexHandler.IndexNameLookup = indexNameLookup
	docIndexHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docIndexHandler).Methods("PUT")

	docCountHandler := bleveHttp.NewDocCountHandler("")
	docCountHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_count", docCountHandler).Methods("GET")

	docGetHandler := bleveHttp.NewDocGetHandler("")
	docGetHandler.IndexNameLookup = indexNameLookup
	docGetHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docGetHandler).Methods("GET")

	docDeleteHandler := bleveHttp.NewDocDeleteHandler("")
	docDeleteHandler.IndexNameLookup = indexNameLookup
	docDeleteHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}", docDeleteHandler).Methods("DELETE")

	searchHandler := bleveHttp.NewSearchHandler("")
	searchHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_search", searchHandler).Methods("POST")

	listFieldsHandler := bleveHttp.NewListFieldsHandler("")
	listFieldsHandler.IndexNameLookup = indexNameLookup
	router.Handle("/api/{indexName}/_fields", listFieldsHandler).Methods("GET")

	debugHandler := bleveHttp.NewDebugDocumentHandler("")
	debugHandler.IndexNameLookup = indexNameLookup
	debugHandler.DocIDLookup = docIDLookup
	router.Handle("/api/{indexName}/{docID}/_debug", debugHandler).Methods("GET")

	aliasHandler := bleveHttp.NewAliasHandler()
	router.Handle("/api/_aliases", aliasHandler).Methods("POST")

	// start the HTTP server
	http.Handle("/", router)
	log.Printf("Listening on %v", *bindAddr)
	log.Fatal(http.ListenAndServe(*bindAddr, nil))
}
