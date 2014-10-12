package main

import (
	"encoding/json"
	"github.com/eknkc/amber"
	"log"
	"net/http"
	"strconv"
)

const (
	TLS_CERT       = "tls/cert-unified.pem"
	TLS_KEY        = "tls/cert.key"
	COUNT_DEFAULT  = 5
	LENGTH_DEFAULT = 5
)

var word_map = map[string][]string{}

func main() {
	log.Printf("Loading word map")
	word_map = LoadWordMap()
	log.Printf("Starting and listening on 4343")
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", Root)
	http.HandleFunc("/passphrases", Passphrases)
	//http.ListenAndServeTLS("0.0.0.0:4343", TLS_CERT, TLS_KEY, nil)
	http.ListenAndServe("0.0.0.0:8000", nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
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

func Passphrases(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	count := COUNT_DEFAULT
	length := LENGTH_DEFAULT

	type passphrase_output struct {
		Count       int
		Length      int
		Passphrases []string
	}
	var output passphrase_output

	log.Printf("got Query: %v\n", query_values)

	var err error
	if val, ok := query_values["count"]; ok {
		count, err = strconv.Atoi(val[0])
		if err != nil {
			log.Printf("WARNING: bad count parameter passed: %v; %v\n", err, val[0])
		}
	}

	if val, ok := query_values["length"]; ok {
		length, err = strconv.Atoi(val[0])
		if err != nil {
			log.Printf("WARNING: bad lenth parameter passed: %v; %v\n", err, val[0])
		}
	}

	log.Printf("COUNT: %v\n", count)
	log.Printf("LENGTH: %v\n", length)

	passphrases := GeneratePassphrases(word_map, count, length)

	for i, val := range passphrases {
		log.Printf("phrase %v: %v\n", i, val)
	}

	output.Count = count
	output.Length = length
	output.Passphrases = passphrases

	//emit json
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(output)
	if err != nil {
		log.Printf("ERROR marshalling json: %v\n", err)
		http.Error(w, "Internal Server Error", 500)
	}
	w.Write(j)
}
