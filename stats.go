package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

var plot_map = map[string]string{} // file basename => proper name
var wordlist_stats = map[string]word_stats{}
var combined_plot_url = ""

func write_distribution_csv() {
	err := os.Mkdir(STATS_PATH, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Error creating stats path: %v\n", err)
	}

	for k, v := range wordlist_stats {
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

func get_plots() {
	plot_map = make(map[string]string)
	for k, _ := range word_map {
		if val, ok := name_map[k]; ok {
			plot_map[k] = val
		} else {
			plot_map[k] = k
		}
	}

	f, err := os.Open(URL_FILE)
	if err != nil {
		log.Fatalf("Error opening plot URL file: %v\n", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		combined_plot_url = scanner.Text() //should only be one line
	}

	if len(combined_plot_url) < 5 {
		log.Fatalf("Malformed plot URL file: %v\n", combined_plot_url)
	}
}

func generate_plots() {
	log.Printf("Generating plots")

	cmd := exec.Command("python", "gen_plots.py")
	cmd.Dir = DATA_PATH
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error generating plots: %v\n", err)
	}
}
