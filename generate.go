package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const (
	WORDNET_PATH   = "data/part-of-speech.txt"
	OFFENSIVE_PATH = "data/offensive.txt"
	// With this algorithm we get best results when we limit the number of consecutive words,
	// then string fragments together with conjunctions. Otherwise we get a really long
	// run-on word salad that is not convincingly grammatical.
	MAGIC_FRAGMENT_LENGTH = 4
)

type GenerateOptions struct {
	count      int
	length     int
	prudish    bool
	no_spaces  bool
	add_digit  bool
	add_symbol bool
}

// word_type -> "can be followed by..."
var GRAMMAR_RULES = map[string][]string{
	"snoun":        []string{"adverb", "verb", "pronoun", "conjunction"},
	"pnoun":        []string{"adverb", "verb", "pronoun", "conjunction"},
	"verb":         []string{"snoun", "pnoun", "preposition", "adjective", "conjunction", "sarticle", "particle"},
	"adjective":    []string{"snoun", "pnoun"},
	"adverb":       []string{"verb"},
	"preposition":  []string{"snoun", "pnoun", "adverb", "adjective", "verb"},
	"pronoun":      []string{"verb", "adverb", "conjunction"},
	"conjunction":  []string{"snoun", "pnoun", "pronoun", "verb", "sarticle", "particle"},
	"sarticle":     []string{"snoun", "adjective"},
	"particle":     []string{"pnoun", "adjective"},
	"interjection": []string{"snoun", "pnoun", "preposition", "adjective", "conjunction", "sarticle", "particle"},
}

var SYMBOLS = []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "-", "+", "_", "="}

var word_types = []string{"snoun", "pnoun", "verb", "adjective", "adverb", "preposition", "pronoun", "conjunction", "sarticle", "particle", "interjection"}

var offensive = map[string]uint{}

func random_word(word_map map[string][]string, word_type string, options GenerateOptions) string {
	grw := func(words []string) (string, bool) {
		word := random_choice(words)
		_, ok := offensive[word]
		return word, ok
	}
	if words, ok := word_map[word_type]; ok {
		word, off := grw(words)
		if options.prudish && off {
			log.Printf("Got offensive word: %v\n", word)
			i := 0
			for i = 0; off && i < 10; i++ {
				word, off = grw(words)
				if off {
					log.Printf("Got offensive word (retry): %v\n", word)
				}
			}
			if i >= 10 {
				log.Printf("Gave up trying to get non-offensive word!")
				word = ""
			}
		}
		return word
	} else {
		log.Printf("WARNING: random_word couldn't find word_type in word_map: %v\n", word_type)
		return "()"
	}
}

// A fragment is an autonomous run of words constructed using grammar rules
func generate_fragment(word_map map[string][]string, fragment_length int, options GenerateOptions) []string {
	fragment_slice := make([]string, fragment_length)
	prev_type_index := random_range(int64(len(word_types) - 1))                     // Random initial word type
	fragment_slice[0] = random_word(word_map, word_types[prev_type_index], options) // Random initial word
	this_word_type := ""
	for i := 1; i < fragment_length; i++ {
		// Get random allowed word type by type of the previous word
		next_word_type_count := int32(len(GRAMMAR_RULES[word_types[prev_type_index]]) - 1)
		if next_word_type_count > 0 { //rand.Int31n cannot take zero as a param
			this_word_type = GRAMMAR_RULES[word_types[prev_type_index]][random_range(int64(next_word_type_count))]
		} else {
			this_word_type = GRAMMAR_RULES[word_types[prev_type_index]][0]
		}
		fragment_slice[i] = random_word(word_map, this_word_type, options) //Random word of the allowed random type
		for j, v := range word_types {                                     // Update previous word type with current word type for next iteration
			if v == this_word_type {
				prev_type_index = int64(j)
			}
		}
	}
	return fragment_slice
}

func generate_passphrase(word_map map[string][]string, options GenerateOptions) []string {
	iterations := options.length / MAGIC_FRAGMENT_LENGTH
	phrase_slice := make([]string, 1)

	phrase_slice = append(phrase_slice, generate_fragment(word_map, MAGIC_FRAGMENT_LENGTH, options)...)
	if iterations >= 1 {
		for i := 1; i <= iterations; i++ {
			phrase_slice = append(phrase_slice, random_word(word_map, "conjunction", options))
			phrase_slice = append(phrase_slice, generate_fragment(word_map, MAGIC_FRAGMENT_LENGTH, options)...)
		}
	}
	return phrase_slice
}

//Generate count number of random passphrases of size length from word_map
func GeneratePassphrases(word_map map[string][]string, options GenerateOptions) []string {
	// Generate count passphrase slices
	// Merge each passphrase slice into a single string
	// Split string by spaces (individual random "words" can actually be multiword phrases)
	// Truncate slice to length words
	// Merge truncated slice back into string
	// Return slice of strings (final random passphrases)
	passphrases := make([]string, options.count)

	var sep string
	if options.no_spaces {
		sep = ""
	} else {
		sep = " "
	}
	for i := 0; i < options.count; i++ {
		ps := generate_passphrase(word_map, options)
		//log.Printf("ps: %v\n", ps)
		pj := strings.Join(ps, " ")
		ps = strings.Split(pj, " ")
		ps = ps[:options.length+1]
		pp := strings.TrimSpace(strings.Join(ps, sep))
		if options.add_digit {
			pp += random_digit()
		}
		if options.add_symbol {
			pp += random_choice(SYMBOLS)
		}
		passphrases[i] = pp
	}
	return passphrases
}

func GenerateStatistics() map[string]word_stats {
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

func LoadOffensiveWords() {
	offensive = make(map[string]uint)

	log.Printf("Loading offensive word list")
	f, err := os.Open(OFFENSIVE_PATH)
	if err != nil {
		log.Fatalf("Error opening offensive word list: %v\n", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		l := scanner.Text()
		offensive[strings.TrimSpace(l)] = 1
	}
}

//Load Wordnet into a mapping of word type to words of that type
func LoadWordMap() map[string][]string {

	word_map := map[string][]string{
		"snoun":        []string{},
		"pnoun":        []string{},
		"verb":         []string{},
		"adjective":    []string{},
		"adverb":       []string{},
		"preposition":  []string{},
		"pronoun":      []string{},
		"conjunction":  []string{},
		"sarticle":     []string{},
		"particle":     []string{},
		"interjection": []string{},
	}

	file, err := os.Open(WORDNET_PATH)
	if err != nil {
		log.Fatalf("Error opening wordnet: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word_type := ""
		plural := false
		line := scanner.Text()
		line_array := strings.Split(line, "\t")
		if len(line_array) != 2 {
			log.Printf("Bad string array length: %v, string: %v", len(line_array), line)
			continue
		}
		word := line_array[0]
		pos_tag := line_array[1]
		if strings.Contains(pos_tag, "N") || strings.Contains(pos_tag, "D") || strings.Contains(pos_tag, "I") {
			if strings.Contains(pos_tag, "P") {
				plural = true
			}
		}
		if strings.Contains(pos_tag, "D") || strings.Contains(pos_tag, "I") {
			if plural {
				word_type = "particle"
			} else {
				word_type = "sarticle"
			}
		} else if strings.Contains(pos_tag, "N") || strings.Contains(pos_tag, "h") || strings.Contains(pos_tag, "o") {
			if plural {
				word_type = "pnoun"
			} else {
				word_type = "snoun"
			}
		} else if strings.Contains(pos_tag, "V") || strings.Contains(pos_tag, "t") || strings.Contains(pos_tag, "i") {
			word_type = "verb"
		} else if strings.Contains(pos_tag, "A") {
			word_type = "adjective"
		} else if strings.Contains(pos_tag, "v") {
			word_type = "adverb"
		} else if strings.Contains(pos_tag, "C") {
			word_type = "conjunction"
		} else if strings.Contains(pos_tag, "p") || strings.Contains(pos_tag, "P") {
			word_type = "preposition"
		} else if strings.Contains(pos_tag, "r") {
			word_type = "pronoun"
		} else if strings.Contains(pos_tag, "!") {
			word_type = "interjection"
		} else {
			log.Printf("Unknown word type! word: %v; pos: %v\n", word, pos_tag)
			continue
		}
		if len(word) > 0 {
			word_map[word_type] = append(word_map[word_type], word)
		} else {
			log.Printf("WARNING: got zero length word: line: %v (interpreted type: %v)", line, word_type)
		}

	}

	return word_map
}
