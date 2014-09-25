package main

import (
	"github.com/eknkc/amber"
	"log"
	"net/http"
)

var TLS_CERT = "tls/cert-unified.pem"
var TLS_KEY = "tls/cert.key"

var word_map = map[string][]string{}

func main() {
	log.Printf("Loading word map")
	word_map = LoadWordMap()
	log.Printf("Starting and listening on 4343")
	http.HandleFunc("/", Root)
	http.ListenAndServeTLS("0.0.0.0:4343", TLS_CERT, TLS_KEY, nil)
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
