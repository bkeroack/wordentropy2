package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	WORDNET_PATH = "part-of-speech.txt"
	// With this algorithm we get best results when we limit the number of consecutive words,
	// then string fragments together with conjunctions. Otherwise we get a really long
	// run-on word salad that is not convincingly grammatical.
	MAGIC_FRAGMENT_LENGTH = 4
)

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

var word_types = []string{"snoun", "pnoun", "verb", "adjective", "adverb", "preposition", "pronoun", "conjunction", "sarticle", "particle", "interjection"}

func random_word(word_map map[string][]string, word_type string) string {
	if words, ok := word_map[word_type]; ok {
		return words[rand.Int31n(int32(len(words)-1))]
	} else {
		log.Printf("WARNING: random_word couldn't find word_type in word_map: %v\n", word_type)
		return "()"
	}
}

// A fragment is an autonomous run of words constructed using grammar rules
func generate_fragment(word_map map[string][]string, fragment_length int) []string {
	fragment_slice := make([]string, fragment_length)
	prev_type_index := rand.Int31n(int32(len(word_types) - 1))             // Random initial word type
	fragment_slice[0] = random_word(word_map, word_types[prev_type_index]) // Random initial word
	for i := 1; i <= fragment_length; i++ {
		// Get random allowed word type by type of the previous word
		this_word_type := GRAMMAR_RULES[word_types[prev_type_index]][rand.Int31n(int32(len(GRAMMAR_RULES[word_types[prev_type_index]])-1))]
		fragment_slice[i] = random_word(word_map, this_word_type) //Random word of the allowed random type
		for j, v := range word_types {                            // Update previous word type with current word type for next iteration
			if v == this_word_type {
				prev_type_index = int32(j)
			}
		}
	}
	return fragment_slice
}

func generate_passphrase(word_map map[string][]string, plen int) [][]string {
	iterations := plen / MAGIC_FRAGMENT_LENGTH
	phrase_slice := make([]string, iterations)

	phrase_slice[0] = stringsgenerate_fragment(word_map, MAGIC_FRAGMENT_LENGTH)
	if iterations >= 1 {
		for i := 1; i <= iterations; i++ {
			fragment_slice := append([]string{random_word(word_map, "conjunction")}, generate_fragment(word_map, MAGIC_FRAGMENT_LENGTH))
		}
	}
	return phrase_slice
}

//Generate count number of random passphrases of size length from word_map
func GeneratePassphrases(word_map map[string][]string, count int, length int) []string {
	// Generate count passphrase slices
	// Merge each passphrase slice into a single string
	// Split string by spaces (individual random "words" can actually be multiword phrases)
	// Truncate slice to length words
	// Merge truncated slice back into string
	// Return slice of strings (final random passphrases)
	passphrases := make([]string, count)
	for i := 0; i < count; i++ {
		ps := generate_passphrase(word_map, length)
		pj := strings.Join(ps, " ")
		ps = strings.Split(pj, " ")
		ps = ps[:length-1]
		passphrases[i] = strings.Join(ps, " ")
	}
	return passphrases
}

//Load Wordnet into a mapping of word type to words of that type
func LoadWordMap() map[string][]string {

	word_map := map[string][]string{
		"snoun":        make([]string, 1),
		"pnoun":        make([]string, 1),
		"verb":         make([]string, 1),
		"adjective":    make([]string, 1),
		"adverb":       make([]string, 1),
		"preposition":  make([]string, 1),
		"pronoun":      make([]string, 1),
		"conjunction":  make([]string, 1),
		"sarticle":     make([]string, 1),
		"particle":     make([]string, 1),
		"interjection": make([]string, 1),
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
		word_map[word_type] = append(word_map[word_type], word)
	}
	for k, v := range word_map {
		log.Printf("Word type: %v; count: %v", k, len(v))
	}
	return word_map
}
