package main

import (
	"encoding/json"
	"fmt"
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
	LISTEN_PORT    = 8000
)

var word_map = map[string][]string{}

func main() {
	log.Printf("Loading word map")
	word_map = LoadWordMap()
	log.Printf("Starting and listening on %v\n", LISTEN_PORT)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Root)
	http.HandleFunc("/passphrases", Passphrases)
	//http.ListenAndServeTLS("0.0.0.0:4343", TLS_CERT, TLS_KEY, nil)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", LISTEN_PORT), nil)
}

func Root(w http.ResponseWriter, r *http.Request) {
	compiler := amber.New()
	err := compiler.ParseFile("views/main.amber")
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad template: main.amber: %v\n", err), 500)
	}
	tpl, err := compiler.Compile()
	// var amber_options amber.Options
	// var amber_diroptions amber.DirOptions
	// amber_options.LineNumbers = true
	// amber_diroptions.Recursive = false
	// tpl_map, err := amber.CompileDir("views/", amber_diroptions, amber_options)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error compiling templates: %v\n", err), 500)
	}
	//tpl_map["main"].Execute(w, nil)
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
