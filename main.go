package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bkeroack/libwordentropy"
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
	DATA_PATH      = "data/"
	URL_FILE       = "data/plot_urls.txt"
)

var g *wordentropy.Generator
var templates = map[string]*template.Template{}
var template_name_map = map[string]string{}
var template_names = [...]string{
	"main.amber",
	"about.amber",
	"random.amber",
}

var localFlag *bool

func init() {
	localFlag = flag.Bool("local", false, "local testing mode (do not generate plots, bind to high port)")
	flag.Parse()

	gomaxprocs := os.ExpandEnv("${GOMAXPROCS}")
	if gomaxprocs == "" {
		log.Printf("GOMAXPROCS not set; default used (1)\n")
	} else {
		log.Printf("GOMAXPROCS: %v\n", gomaxprocs)
	}
}

func main() {

	log.Printf("Loading word map")

	var err error
	g, err = wordentropy.LoadGenerator(&wordentropy.WordListOptions{
		Wordlist:  "data/part-of-speech.txt",
		Offensive: "data/offensive.txt",
	})
	if err != nil {
		log.Fatalf("Error loading wordlist: %v\n", err)
	}

	word_map := g.GetWordMap()

	wordlist_stats := GenerateStatistics(word_map)

	for k, v := range wordlist_stats {
		log.Printf("Word type: %v; total_count: %v; largest_word: %v\n",
			k, v.Total_count, v.Max_char_count)
	}

	write_distribution_csv()

	port := LISTEN_PORT
	if *localFlag {
		port = 4343
	} else {
		generate_plots()
	}

	get_plots(word_map, wordlist_stats)
	compile_templates()

	config := &tls.Config{MinVersion: tls.VersionTLS10}
	server := &http.Server{Addr: fmt.Sprintf("0.0.0.0:%v", port), Handler: nil, TLSConfig: config}

	log.Printf("Starting and listening on %v\n", port)
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", Root)
	http.HandleFunc("/about", About)
	http.HandleFunc("/how-random", Random)
	http.HandleFunc("/passphrases", Passphrases)
	err = server.ListenAndServeTLS(TLS_CERT, TLS_KEY)
	if err != nil {
		log.Fatalf("Error starting HTTP listener: %v\n", err)
	}
}

func compile_template(n string) *template.Template {
	compiler := amber.New()
	err := compiler.ParseFile(fmt.Sprintf("views/%v", n))
	if err != nil {
		log.Fatalf("Bad template: %v: %v\n", n, err)
	}
	tpl, err := compiler.Compile()
	if err != nil {
		log.Fatalf("Error compiling template: %v: %v\n", n, err)
	}
	return tpl
}

func compile_templates() {
	log.Printf("Compiling templates...\n")
	templates = make(map[string]*template.Template)
	template_name_map = make(map[string]string)
	for t := range template_names {
		template := template_names[t]
		n := strings.Split(template, ".")[0]
		templates[n] = compile_template(template)
		template_name_map[n] = template
	}
	log.Printf("Compilation complete\n")
}

func execute_template(n string, w http.ResponseWriter, r *http.Request, data interface{}) {
	log.Printf("%v\t%v\t%v\n", n, r.RemoteAddr, r.UserAgent())
	var err error
	if !*localFlag {
		err = templates[n].Execute(w, data)
	} else {
		tmpl := compile_template(template_name_map[n])
		err = tmpl.Execute(w, data)
	}
	if err != nil {
		log.Printf("Error executing template: %v: %v\n", n, err)
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	execute_template("main", w, r, nil)
}

func About(w http.ResponseWriter, r *http.Request) {
	execute_template("about", w, r, nil)
}

func Random(w http.ResponseWriter, r *http.Request) {
	log.Printf("random\t%v\t%v\n", r.RemoteAddr, r.UserAgent())
	data := struct {
		Word_stats        []stat_ui
		Plots             map[string]string
		Combined_plot_url string
	}{
		sanitized_stats,
		plot_map,
		combined_plot_url,
	}
	execute_template("random", w, r, &data)
}

func process_passphrases_options(o *wordentropy.GenerateOptions, qv map[string][]string) (bool, string) {

	if val, ok := qv["count"]; ok {
		count, err := strconv.Atoi(val[0])
		if err != nil || count < 0 || count > 99 {
			return false, fmt.Sprintf("Bad count parameter passed: %v; %v\n", err, val[0])
		}
		o.Count = uint(count)
	}

	if val, ok := qv["length"]; ok {
		length, err := strconv.Atoi(val[0])
		if err != nil || length < 0 || length > 99 {
			return false, fmt.Sprintf("Bad length parameter passed: %v; %v\n", err, val[0])
		}
		o.Length = uint(length)
	}

	if val, ok := qv["prudish"]; ok {
		if val[0] == "true" {
			o.Prudish = true
		}
	}
	if val, ok := qv["no_spaces"]; ok {
		if val[0] == "true" {
			o.No_spaces = true
		}
	}
	if val, ok := qv["add_digit"]; ok {
		if val[0] == "true" {
			o.Add_digit = true
		}
	}
	if val, ok := qv["add_symbol"]; ok {
		if val[0] == "true" {
			o.Add_symbol = true
		}
	}
	return true, ""
}

func Passphrases(w http.ResponseWriter, r *http.Request) {
	query_values := r.URL.Query()
	options := wordentropy.GenerateOptions{}
	w.Header().Set("Content-Type", "application/json")

	type passphrase_output struct {
		Count       uint
		Length      uint
		Passphrases []string
	}

	var output passphrase_output

	ok, msg := process_passphrases_options(&options, query_values)
	if !ok {
		emit_json_error(w, msg, 401)
	}

	log.Printf("passphrases\toptions:%v\t%v\t%v\n", options, r.RemoteAddr, r.UserAgent())

	passphrases, err := g.GeneratePassphrases(&options)
	if err != nil {
		emit_json_error(w, "Error generating passphrases", 500)
	}

	output.Count = options.Count
	output.Length = options.Length
	output.Passphrases = passphrases

	//emit json
	j, err := json.Marshal(output)
	if err != nil {
		log.Printf("ERROR marshalling json: %v\n", err)
		http.Error(w, "{ \"error\": \"Internal Server Error\" }", 500)
	}
	w.Write(j)
}
