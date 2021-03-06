package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/dustin/go-humanize"
	"log"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type word_stats struct {
	Total_count      int
	Max_char_count   int
	Distribution_map map[int]int
}

var plot_map = map[string]string{} // file basename => proper name
var wordlist_stats = map[string]word_stats{}

type stat_ui struct {
	Name      string
	Count     string
	Count_int int
}

type Stat_ui []stat_ui

var sanitized_stats = Stat_ui{}
var combined_plot_url = ""
var name_map = map[string]string{
	"particle": "Plural Article",
	"sarticle": "Singular Article",
	"pnoun":    "Plural Noun",
	"snoun":    "Singular Noun",
	"ALL":      "All Words (total)",
}

func (slice Stat_ui) Len() int {
	return len(slice)
}

func (slice Stat_ui) Less(i, j int) bool {
	return slice[i].Count_int < slice[j].Count_int
}

func (slice Stat_ui) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

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

func get_plots(word_map map[string][]string, wordlist_stats map[string]word_stats) {
	// used by view layer for plot image names and titles
	plot_map = make(map[string]string)
	for k, _ := range word_map {
		if val, ok := name_map[k]; ok {
			plot_map[k] = val
		} else {
			plot_map[k] = k
		}
	}

	// sanitize stats listing with proper names
	if len(wordlist_stats) > 0 {
		for k, v := range wordlist_stats {
			hum_v := humanize.Comma(int64(v.Total_count))
			if x, ok := name_map[k]; ok {
				sanitized_stats = append(sanitized_stats, stat_ui{x, hum_v, v.Total_count})
			} else {
				first_letter := string([]rune(k)[0])
				remainder := string([]rune(k)[1:len(k)])
				capitalized := fmt.Sprintf("%v%v", strings.ToUpper(first_letter), remainder)
				sanitized_stats = append(sanitized_stats, stat_ui{capitalized, hum_v, v.Total_count})
			}
		}
	} else {
		log.Fatalf("get_plots: wordlist_stats not initialized!\n")
	}

	sort.Sort(sort.Reverse(sanitized_stats))

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

func GenerateStatistics(word_map map[string][]string) map[string]word_stats {
	statistics := make(map[string]word_stats)
	distribution_map := make(map[string]map[int]int)
	distribution_map["ALL"] = make(map[int]int)

	global_max_len := 0
	global_word_count := 0
	for k, val := range word_map {
		if _, ok := distribution_map[k]; !ok {
			distribution_map[k] = make(map[int]int)
		}
		stat := word_stats{len(val), 0, map[int]int{}}
		global_word_count += len(val)
		max_len := 0
		for w := range val {
			word := val[w]
			word_len := len(word)
			if word_len == 0 {
				log.Printf("WARNING: found zero length word (type: %v, index: %v)\n", k, w)
			}
			if word_len > max_len {
				max_len = word_len
			}
			if max_len > global_max_len {
				global_max_len = max_len
			}
			if _, ok := distribution_map[k][word_len]; ok {
				distribution_map[k][word_len]++
			} else {
				distribution_map[k][word_len] = 1
			}
			if _, ok := distribution_map["ALL"][word_len]; ok {
				distribution_map["ALL"][word_len]++
			} else {
				distribution_map["ALL"][word_len] = 1
			}
		}
		stat.Max_char_count = max_len
		stat.Distribution_map = distribution_map[k]
		statistics[k] = stat
	}
	statistics["ALL"] = word_stats{global_word_count, global_max_len, distribution_map["ALL"]}
	return statistics
}
