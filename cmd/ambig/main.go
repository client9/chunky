package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky/tok"
)

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
	byTags := flag.Bool("tags", false, "aggregate by tag combination instead of by word")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: ambig [flags] <file> [file ...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	freq := make(map[string]int)
	tagSets := make(map[string]string) // word → "TAG1 TAG2 ..." (stable per word form)

	for _, path := range flag.Args() {
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
				sentence := stripTags(tagged)
				if sentence == "" {
					continue
				}
				for _, s := range tok.Parse(sentence) {
					for _, t := range s.Tokens {
						if len(t.Tags) <= 1 {
							continue
						}
						freq[t.Word]++
						if _, seen := tagSets[t.Word]; !seen {
							parts := make([]string, len(t.Tags))
							for i, tag := range t.Tags {
								parts[i] = tag.String()
							}
							tagSets[t.Word] = strings.Join(parts, " ")
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

	if !*byTags {
		total := 0
		for _, w := range wf {
			total += w.Count
		}
		cumulative := 0
		for i, w := range wf {
			cumulative += w.Count
			fmt.Printf("%6d. %6d  %5.1f%%  %-20s  %s\n", i+1, w.Count, 100*float64(cumulative)/float64(total), w.Word, tagSets[w.Word])
		}
		return
	}

	// Aggregate by tag combination.
	comboCount := make(map[string]int)
	comboWords := make(map[string][]wordFreq)
	for _, w := range wf {
		combo := tagSets[w.Word]
		comboCount[combo] += w.Count
		comboWords[combo] = append(comboWords[combo], w)
	}

	type comboFreq struct {
		Tags  string
		Count int
	}
	cf := make([]comboFreq, 0, len(comboCount))
	for tags, n := range comboCount {
		cf = append(cf, comboFreq{tags, n})
	}
	sort.Slice(cf, func(i, j int) bool {
		if cf[i].Count != cf[j].Count {
			return cf[i].Count > cf[j].Count
		}
		return cf[i].Tags < cf[j].Tags
	})

	total := 0
	for _, c := range cf {
		total += c.Count
	}
	cumulative := 0
	for i, c := range cf {
		cumulative += c.Count
		// Top 5 words for this combination.
		words := comboWords[c.Tags]
		top := words
		if len(top) > 5 {
			top = top[:5]
		}
		examples := make([]string, len(top))
		for j, w := range top {
			examples[j] = w.Word
		}
		fmt.Printf("%6d. %6d  %5.1f%%  %-30s  %s\n", i+1, c.Count, 100*float64(cumulative)/float64(total), c.Tags, strings.Join(examples, " "))
	}
}
