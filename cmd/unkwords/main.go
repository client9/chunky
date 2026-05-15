package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky/tok"
)

// stripTags reconstructs the original sentence from a "word/TAG ..." line.
func stripTags(line string) string {
	parts := strings.Fields(line)
	words := make([]string, 0, len(parts))
	for _, p := range parts {
		i := strings.LastIndex(p, "/")
		if i > 0 {
			words = append(words, p[:i])
		}
	}
	return strings.Join(words, " ")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: unkwords <file> [file ...]\n")
		os.Exit(1)
	}

	freq := make(map[string]int)
	spacyTags := make(map[string]map[string]int) // word → spaCy tag → count

	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", path, err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

		for scanner.Scan() {
			var rec struct {
				Type      string   `json:"type"`
				Sentences []string `json:"sentences"`
			}
			if err := json.Unmarshal([]byte(scanner.Text()), &rec); err != nil || rec.Type == "header" {
				continue
			}
			for _, tagged := range rec.Sentences {
				// Build word→tag lookup from spaCy's tokenization for this sentence.
				spacyForSent := make(map[string]string)
				for _, part := range strings.Fields(tagged) {
					i := strings.LastIndex(part, "/")
					if i > 0 {
						spacyForSent[part[:i]] = part[i+1:]
					}
				}

				sentence := stripTags(tagged)
				if sentence == "" {
					continue
				}
				var tokens []tok.Token
				for _, s := range tok.Parse(sentence) {
					tokens = append(tokens, s.Tokens...)
				}
				for _, t := range tokens {
					if t.IsUnknownTag() {
						freq[t.Word]++
						if tag, ok := spacyForSent[t.Word]; ok {
							if spacyTags[t.Word] == nil {
								spacyTags[t.Word] = make(map[string]int)
							}
							spacyTags[t.Word][tag]++
						}
					}
				}
			}
		}
		f.Close()

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	// Sort by descending frequency, then alphabetically.
	type wordFreq struct {
		Word  string
		Count int
	}
	wf := make([]wordFreq, 0, len(freq))
	for w, n := range freq {
		wf = append(wf, wordFreq{w, n})
	}
	sort.Slice(wf, func(i, j int) bool {
		if wf[i].Count != wf[j].Count {
			return wf[i].Count > wf[j].Count
		}
		return wf[i].Word < wf[j].Word
	})

	for i, w := range wf {
		// Format spaCy tags sorted by descending count.
		tagCounts := spacyTags[w.Word]
		type tagFreq struct {
			Tag   string
			Count int
		}
		tf := make([]tagFreq, 0, len(tagCounts))
		for tag, n := range tagCounts {
			tf = append(tf, tagFreq{tag, n})
		}
		sort.Slice(tf, func(a, b int) bool {
			if tf[a].Count != tf[b].Count {
				return tf[a].Count > tf[b].Count
			}
			return tf[a].Tag < tf[b].Tag
		})
		parts := make([]string, len(tf))
		for j, t := range tf {
			parts[j] = fmt.Sprintf("%s:%d", t.Tag, t.Count)
		}
		spacyStr := strings.Join(parts, " ")
		fmt.Printf("%6d. %6d  %-30s  %s\n", i+1, w.Count, w.Word, spacyStr)
	}
}
