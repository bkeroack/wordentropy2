package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/eknkc/amber"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	TLS_CERT       = "tls/cert-unified.pem"
	TLS_KEY        = "tls/cert.key"
	COUNT_DEFAULT  = 5
	LENGTH_DEFAULT = 5
	LISTEN_PORT    = 443
	STATS_PATH     = "data/stats/"
)

var word_map = map[string][]string{}
var templates = map[string]*template.Template{}
var template_names = [...]string{
	"main.amber",
	"about.amber",
	"random.amber",
}

func write_distribution_csv(stats map[string]word_stats) {
	err := os.Mkdir(STATS_PATH, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating stats path: %v\n", err)
	}

	for k, v := range stats {
		f, err := os.Create(fmt.Sprintf("%v/%v.csv", STATS_PATH, k))
		if err != nil {
			log.Fatalf("Error creating stats csv for %v: %v\n", k, err)
		}
		w := csv.NewWriter(f)
		dist := v.Distribution_map
		for l, c := range dist {
			w.Write([]string{strconv.Itoa(l), strconv.Itoa(c)})
		}
		w.Flush()
		f.Close()
	}
}

func generate_plots(csv string) bool {
	return false
}

func main() {
	log.Printf("Loading word map")

	word_map = LoadWordMap()
	wordlist_stats := GenerateStatistics()

	for k, v := range wordlist_stats {
		log.Printf("Word type: %v; total_count: %v; largest_word: %v\n",
			k, v.Total_count, v.Max_char_count)
	}

	write_distribution_csv(wordlist_stats)

	compile_templates()

	config := &tls.Config{MinVersion: tls.VersionTLS10}
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%v", LISTEN_PORT), Handler: nil, TLSConfig: config}

	log.Printf("Starting and listening on %v\n", LISTEN_PORT)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Root)
	http.HandleFunc("/about", About)
	http.HandleFunc("/how-random", Random)
	http.HandleFunc("/passphrases", Passphrases)
	server.ListenAndServeTLS(TLS_CERT, TLS_KEY)
}

func compile_templates() {
	log.Printf("Compiling templates...\n")
	templates = make(map[string]*template.Template)
	compiler := amber.New()
	for t := range template_names {
		template := template_names[t]
		err := compiler.ParseFile(fmt.Sprintf("views/%v", template))
		if err != nil {
			log.Fatalf("Bad template: %v: %v\n", template, err)
		}
		tpl, err := compiler.Compile()
		if err != nil {
			log.Fatalf("Error compiling template: %v: %v\n", template, err)
		}
		templates[strings.Split(template, ".")[0]] = tpl
	}
	log.Printf("Compilation complete\n")
}

func Root(w http.ResponseWriter, r *http.Request) {
	log.Printf("root\t%v\t%v\n", r.RemoteAddr, r.UserAgent())
	templates["main"].Execute(w, nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	log.Printf("about\t%v\t%v\n", r.RemoteAddr, r.UserAgent())
	templates["about"].Execute(w, nil)
}

func Random(w http.ResponseWriter, r *http.Request) {
	log.Printf("random\t%v\t%v\n", r.RemoteAddr, r.UserAgent())
	templates["random"].Execute(w, nil)
}

func Passphrases(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	count := COUNT_DEFAULT
	length := LENGTH_DEFAULT
	w.Header().Set("Content-Type", "application/json")

	type passphrase_output struct {
		Count       int
		Length      int
		Passphrases []string
	}

	type error_msg struct {
		Error string
	}

	var output passphrase_output
	var err_json error_msg

	var err error
	if val, ok := query_values["count"]; ok {
		count, err = strconv.Atoi(val[0])
		if err != nil {
			err_json.Error = fmt.Sprintf("Bad count parameter passed: %v; %v\n", err, val[0])
			err_str, _ := json.Marshal(err_json)
			log.Printf(err_json.Error)
			http.Error(w, string(err_str), 401)
		}
	}

	if val, ok := query_values["length"]; ok {
		length, err = strconv.Atoi(val[0])
		if err != nil {
			err_json.Error = fmt.Sprintf("Bad length parameter passed: %v; %v\n", err, val[0])
			err_str, _ := json.Marshal(err_json)
			log.Printf(err_json.Error)
			http.Error(w, string(err_str), 401)
		}
	}

	log.Printf("passphrases\tcount=%v\tlength=%v\t%v\t%v\n", count, length, r.RemoteAddr, r.UserAgent())

	passphrases := GeneratePassphrases(word_map, count, length)

	output.Count = count
	output.Length = length
	output.Passphrases = passphrases

	//emit json
	j, err := json.Marshal(output)
	if err != nil {
		log.Printf("ERROR marshalling json: %v\n", err)
		http.Error(w, "{ \"error\": \"Internal Server Error\" }", 500)
	}
	w.Write(j)
}
