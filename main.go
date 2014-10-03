package main

import (
	"github.com/eknkc/amber"
	"log"
	"net/http"
	"strconv"
)

const (
	TLS_CERT       = "tls/cert-unified.pem"
	TLS_KEY        = "tls/cert.key"
	COUNT_DEFAULT  = 5
	LENGTH_DEFAULT = 3
)

var word_map = map[string][]string{}

func main() {
	log.Printf("Loading word map")
	word_map = LoadWordMap()
	log.Printf("Starting and listening on 4343")
	http.Handle("/static", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/", Root)
	http.ListenAndServeTLS("0.0.0.0:4343", TLS_CERT, TLS_KEY, nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	count := COUNT_DEFAULT
	length := LENGTH_DEFAULT

	if val := query_values.Get("count"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Bad value for count parameter", 400)
		} else {
			count = i
		}
	}

	if val := query_values.Get("length"); val != "" {
		i, err := strconv.Atoi(val)
		if err != nil {
			http.Error(w, "Bad value for length parameter", 400)
		} else {
			length = i
		}
	}

	passphrases := GeneratePassphrases(word_map, count, length)

	for i, val := range passphrases {
		log.Printf("phrase: %v\n", val)
	}

	compiler := amber.New()
	err := compiler.ParseFile("views/main.amber")
	if err != nil {
		http.Error(w, "Bad template: main.amber", 500)
	}
	tpl, err := compiler.Compile()
	if err != nil {
		http.Error(w, "Error compiling template: main.amber", 500)
	}
	tpl.Execute(w, nil)
}
