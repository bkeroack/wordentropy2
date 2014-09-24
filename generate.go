package generate

import (
	"bufio"
	"log"
	"os"
	"strings"
)

var WORDNET_PATH = "part-of-speech.txt"

//Load Wordnet into a mapping of word type to words of that type
func LoadWordMap() map[string][]string {

	word_map := make(map[string][]string)

	file, err := os.Open(WORDNET_PATH)
	if err != nil {
		log.Fatalf("Error opening wordnet: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		plural := false
		line := scanner.Text()
		line_array := strings.Split(line, '\t')
		if len(line_array) != 2{
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
		if val, ok := 
	}
}
