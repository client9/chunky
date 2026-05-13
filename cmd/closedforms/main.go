package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

var closedTags = map[string]bool{
	"ADP":   true,
	"AUX":   true,
	"CCONJ": true,
	"DET":   true,
	"PART":  true,
	"PRON":  true,
	"SCONJ": true,
}

func main() {
	goOutput := flag.Bool("go", false, "output a Go source file instead of a text report")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: closedforms [-go] file1 file2 ...\n")
		os.Exit(1)
	}

	// word -> tag -> count (all occurrences of any word that appears with a closed tag)
	counts := make(map[string]map[string]int)
	closedWords := make(map[string]bool)

	for _, path := range flag.Args() {
		f, err := os.Open(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", path, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 1<<20), 1<<20)
		for scanner.Scan() {
			for token := range strings.FieldsSeq(scanner.Text()) {
				idx := strings.LastIndex(token, "/")
				if idx <= 0 {
					continue
				}
				word := strings.ToLower(token[:idx])
				tag := token[idx+1:]
				if tag == "" || strings.ContainsAny(word, ".'\u2018\u2019") || (len(word) == 1 && word != "i" && word != "a") {
					continue
				}
				if closedTags[tag] {
					closedWords[word] = true
				}
				if counts[word] == nil {
					counts[word] = make(map[string]int)
				}
				counts[word][tag]++
			}
		}
		f.Close()
	}

	words := make([]string, 0, len(closedWords))
	for w := range closedWords {
		words = append(words, w)
	}
	sort.Strings(words)

	type tagFreq struct {
		tag  string
		freq float64
	}

	// topTags returns the tags covering up to 95% cumulative frequency, after
	// applying the min-count and min-closed-fraction filters.
	topTags := func(word string) []tagFreq {
		tagCounts := counts[word]
		total := 0
		for _, c := range tagCounts {
			total += c
		}
		if total < 100 {
			return nil
		}
		closedCount := 0
		for tag, c := range tagCounts {
			if closedTags[tag] {
				closedCount += c
			}
		}
		if float64(closedCount)/float64(total) < 0.10 {
			return nil
		}

		freqs := make([]tagFreq, 0, len(tagCounts))
		for tag, c := range tagCounts {
			freqs = append(freqs, tagFreq{tag, float64(c) / float64(total)})
		}
		sort.Slice(freqs, func(i, j int) bool {
			if freqs[i].freq != freqs[j].freq {
				return freqs[i].freq > freqs[j].freq
			}
			return freqs[i].tag < freqs[j].tag
		})

		var result []tagFreq
		cumulative := 0.0
		for _, tf := range freqs {
			result = append(result, tf)
			cumulative += tf.freq
			if cumulative >= 0.95 {
				break
			}
		}
		return result
	}

	if *goOutput {
		fmt.Println("package pos")
		fmt.Println()
		fmt.Println("// ClosedFormTags maps closed-class word forms to their most likely POS tags,")
		fmt.Println("// covering at least 95% of observed usage.")
		fmt.Println("var ClosedFormTags = map[string][]Tag{")
		for _, word := range words {
			freqs := topTags(word)
			if freqs == nil {
				continue
			}
			tags := make([]string, len(freqs))
			for i, tf := range freqs {
				tags[i] = "Tag" + tf.tag
			}
			fmt.Printf("\t%q: {%s},\n", word, strings.Join(tags, ", "))
		}
		fmt.Println("}")
		return
	}

	for _, word := range words {
		freqs := topTags(word)
		if freqs == nil {
			continue
		}
		parts := make([]string, len(freqs))
		for i, tf := range freqs {
			parts[i] = fmt.Sprintf("%.2f %s", tf.freq, tf.tag)
		}
		fmt.Printf("%s: %s\n", word, strings.Join(parts, ", "))
	}
}
