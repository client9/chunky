package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// Token holds a word and its POS tag parsed from "word/TAG" format.
type Token struct {
	Word string
	Tag  string
}

// FeatureFunc takes the previous, current, and next token and returns a feature
// key string. Return "" to skip this position.
type FeatureFunc func(prev, curr, next Token) string

// features is the registry of all named feature templates.
// Add new entries here to extend the system.
var features = []struct {
	Name string
	Fn   FeatureFunc
}{
	{
		Name: "prevtag+nexttag",
		Fn: func(prev, curr, next Token) string {
			return prev.Tag + "-" + next.Tag
		},
	},
}

// counts[featureName][featureKey][tag] = count
type countMap map[string]map[string]map[string]int

func parseToken(s string) (Token, bool) {
	i := strings.LastIndex(s, "/")
	if i <= 0 || i == len(s)-1 {
		return Token{}, false
	}
	return Token{Word: s[:i], Tag: s[i+1:]}, true
}

func parseLine(line string) []Token {
	parts := strings.Fields(line)
	tokens := make([]Token, 0, len(parts))
	for _, p := range parts {
		if t, ok := parseToken(p); ok {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func processFile(path string, counts countMap) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		tokens := parseLine(scanner.Text())
		if len(tokens) < 3 {
			continue
		}
		for i := 1; i < len(tokens)-1; i++ {
			prev, curr, next := tokens[i-1], tokens[i], tokens[i+1]
			for _, feat := range features {
				key := feat.Fn(prev, curr, next)
				if key == "" {
					continue
				}
				if counts[feat.Name] == nil {
					counts[feat.Name] = make(map[string]map[string]int)
				}
				if counts[feat.Name][key] == nil {
					counts[feat.Name][key] = make(map[string]int)
				}
				counts[feat.Name][key][curr.Tag]++
			}
		}
	}
	return scanner.Err()
}

func keyTotal(tagCounts map[string]int) int {
	n := 0
	for _, c := range tagCounts {
		n += c
	}
	return n
}

func printHistogram(counts countMap, sortByCount bool) {
	featNames := make([]string, 0, len(counts))
	for name := range counts {
		featNames = append(featNames, name)
	}
	sort.Strings(featNames)

	for _, name := range featNames {
		fmt.Printf("=== %s ===\n", name)
		keyMap := counts[name]

		keys := make([]string, 0, len(keyMap))
		for k := range keyMap {
			keys = append(keys, k)
		}

		if sortByCount {
			sort.Slice(keys, func(i, j int) bool {
				ti := keyTotal(keyMap[keys[i]])
				tj := keyTotal(keyMap[keys[j]])
				if ti != tj {
					return ti > tj
				}
				return keys[i] < keys[j]
			})
		} else {
			sort.Strings(keys)
		}

		grandTotal := 0
		for _, tc := range keyMap {
			grandTotal += keyTotal(tc)
		}

		cumulative := 0
		for i, key := range keys {
			tagCounts := keyMap[key]
			total := keyTotal(tagCounts)
			cumulative += total

			// Sort tags by descending count.
			type tagFreq struct {
				Tag   string
				Count int
			}
			tf := make([]tagFreq, 0, len(tagCounts))
			for tag, n := range tagCounts {
				tf = append(tf, tagFreq{tag, n})
			}
			sort.Slice(tf, func(i, j int) bool {
				if tf[i].Count != tf[j].Count {
					return tf[i].Count > tf[j].Count
				}
				return tf[i].Tag < tf[j].Tag
			})

			parts := make([]string, 0, len(tf))
			for _, t := range tf {
				pct := float64(t.Count) / float64(total) * 100
				parts = append(parts, fmt.Sprintf("%s: %.1f%%", t.Tag, pct))
			}
			cumPct := float64(cumulative) / float64(grandTotal) * 100
			fmt.Printf("  %4d. %s (n=%d, cum=%.1f%%):  %s\n", i+1, key, total, cumPct, strings.Join(parts, ", "))
		}
		fmt.Println()
	}
}

func main() {
	sortByCount := flag.Bool("sort-count", false, "sort feature keys by total count descending (default: alphabetical)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: lexrules [flags] <file> [file ...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	counts := make(countMap)

	for _, path := range flag.Args() {
		if err := processFile(path, counts); err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	printHistogram(counts, *sortByCount)
}
