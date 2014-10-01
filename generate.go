package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var WORDNET_PATH = "part-of-speech.txt"

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
