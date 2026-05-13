package main

import (
	"bufio"
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

	for _, path := range os.Args[1:] {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", path, err)
			os.Exit(1)
		}

		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

		for scanner.Scan() {
			sentence := stripTags(scanner.Text())
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
		fmt.Printf("%6d. %6d  %s\n", i+1, w.Count, w.Word)
	}
}
