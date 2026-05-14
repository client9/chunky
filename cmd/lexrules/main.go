package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky"
)

// Token holds a word and its POS tag parsed from "word/TAG" format.
type Token struct {
	Word string
	Tag  string
}

// FeatureFunc takes a sliding window of tokens and returns a feature key.
// Return "" to skip this position.
type FeatureFunc func(prev2, prev, curr, next, next2 Token) string

// features is the registry of all named feature templates.
var features = []struct {
	Name string
	Desc string
	Fn   FeatureFunc
}{
	{
		Name: "prevtag",
		Desc: "previous tag → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev.Tag
		},
	},
	{
		Name: "nexttag",
		Desc: "next tag → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return next.Tag
		},
	},
	{
		Name: "prevtag+nexttag",
		Desc: "prev+next tag sandwich → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev.Tag + "+" + next.Tag
		},
	},
	{
		Name: "prev2tag+prevtag",
		Desc: "trigram lookback → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev2.Tag + "+" + prev.Tag
		},
	},
	{
		Name: "prevtag+nexttag+next2tag",
		Desc: "prev + next two tags → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev.Tag + "+" + next.Tag + "+" + next2.Tag
		},
	},
	{
		Name: "prev2tag+prevtag+nexttag",
		Desc: "two prev + one next (3-token window) → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev2.Tag + "+" + prev.Tag + "+" + next.Tag
		},
	},
	{
		Name: "prev2tag+prevtag+nexttag+next2tag",
		Desc: "two prev + two next (4-token window) → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return prev2.Tag + "+" + prev.Tag + "+" + next.Tag + "+" + next2.Tag
		},
	},
	{
		Name: "word+prevtag",
		Desc: "word and previous tag → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return curr.Word + "+" + prev.Tag
		},
	},
	{
		Name: "word+nexttag",
		Desc: "word and next tag → current tag distribution",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return curr.Word + "+" + next.Tag
		},
	},
	{
		Name: "word",
		Desc: "word alone → tag distribution (ambiguity profile per word)",
		Fn: func(prev2, prev, curr, next, next2 Token) string {
			return curr.Word
		},
	},
}

func featureByName(name string) (FeatureFunc, bool) {
	for _, f := range features {
		if f.Name == name {
			return f.Fn, true
		}
	}
	return nil, false
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

// mergeCompounds scans tokens left-to-right and replaces matching word
// sequences from chunky.CompoundTags with a single merged token.
// Longest match wins. The merged word is underscore-joined lowercase.
func mergeCompounds(tokens []Token) []Token {
	if len(tokens) == 0 {
		return tokens
	}
	out := make([]Token, 0, len(tokens))
	i := 0
	for i < len(tokens) {
		merged := false
		for length := chunky.CompoundMaxLen; length >= 2; length-- {
			if i+length > len(tokens) {
				continue
			}
			words := make([]string, length)
			for j := 0; j < length; j++ {
				words[j] = strings.ToLower(tokens[i+j].Word)
			}
			key := strings.Join(words, " ")
			if tag, ok := chunky.CompoundTags[key]; ok {
				out = append(out, Token{Word: strings.Join(words, "_"), Tag: tag.String()})
				i += length
				merged = true
				break
			}
		}
		if !merged {
			out = append(out, tokens[i])
			i++
		}
	}
	return out
}

func parseLine(line string) []Token {
	parts := strings.Fields(line)
	tokens := make([]Token, 0, len(parts))
	for _, p := range parts {
		if t, ok := parseToken(p); ok {
			tokens = append(tokens, t)
		}
	}
	return mergeCompounds(tokens)
}

func processFile(path string, selected []struct {
	Name string
	Fn   FeatureFunc
}, counts countMap) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var empty Token
	for scanner.Scan() {
		tokens := parseLine(scanner.Text())
		if len(tokens) < 3 {
			continue
		}
		for i := 1; i < len(tokens)-1; i++ {
			var prev2, next2 Token
			if i >= 2 {
				prev2 = tokens[i-2]
			} else {
				prev2 = empty
			}
			if i+2 < len(tokens) {
				next2 = tokens[i+2]
			} else {
				next2 = empty
			}
			prev, curr, next := tokens[i-1], tokens[i], tokens[i+1]
			for _, feat := range selected {
				key := feat.Fn(prev2, prev, curr, next, next2)
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

// udTags is the fixed ordered set of UD tags used as matrix columns.
// Tags not in this list are summed into "other".
var udTags = []string{
	"ADJ", "ADP", "ADV", "AUX", "CCONJ", "DET",
	"INTJ", "NOUN", "NUM", "PART", "PRON", "PROPN",
	"PUNCT", "SCONJ", "SYM", "VERB", "X",
}

func keyTotal(tagCounts map[string]int) int {
	n := 0
	for _, c := range tagCounts {
		n += c
	}
	return n
}

func printMatrix(counts countMap, minN int) {
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// header
	fmt.Fprintf(w, "feature,key")
	for _, tag := range udTags {
		fmt.Fprintf(w, ",%s", strings.ToLower(tag))
	}
	fmt.Fprintf(w, ",other,total\n")

	featNames := make([]string, 0, len(counts))
	for name := range counts {
		featNames = append(featNames, name)
	}
	sort.Strings(featNames)

	for _, name := range featNames {
		keyMap := counts[name]
		keys := make([]string, 0, len(keyMap))
		for k := range keyMap {
			if keyTotal(keyMap[k]) >= minN {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)

		for _, key := range keys {
			tagCounts := keyMap[key]
			total := keyTotal(tagCounts)
			other := total
			fmt.Fprintf(w, "%s,%s", name, key)
			for _, tag := range udTags {
				n := tagCounts[tag]
				other -= n
				fmt.Fprintf(w, ",%d", n)
			}
			fmt.Fprintf(w, ",%d,%d\n", other, total)
		}
	}
}

func printHistogram(counts countMap, sortByCount bool, minN int, coverage float64) {
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
			if keyTotal(keyMap[k]) >= minN {
				keys = append(keys, k)
			}
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
			tagCum := 0.0
			for _, t := range tf {
				pct := float64(t.Count) / float64(total) * 100
				parts = append(parts, fmt.Sprintf("%s: %.1f%%", t.Tag, pct))
				tagCum += pct
				if tagCum >= coverage {
					break
				}
			}
			cumPct := float64(cumulative) / float64(grandTotal) * 100
			fmt.Printf("  %4d. %-30s (n=%d, cum=%.1f%%):  %s\n", i+1, key, total, cumPct, strings.Join(parts, ", "))
		}
		fmt.Println()
	}
}

func main() {
	var featFlags []string
	flag.Func("feat", "feature template name; may be repeated (default: list available)", func(s string) error {
		for _, name := range strings.Split(s, ",") {
			name = strings.TrimSpace(name)
			if name != "" {
				featFlags = append(featFlags, name)
			}
		}
		return nil
	})
	fmt_ := flag.String("fmt", "text", "output format: text or matrix (CSV)")
	sortByCount := flag.Bool("sort-count", false, "sort feature keys by total count descending (default: alphabetical)")
	minN := flag.Int("min-n", 1, "minimum total count to include a feature key in output")
	coverage := flag.Float64("coverage", 90, "show only the top tags whose cumulative share reaches this percentage (0 = show all)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: lexrules [flags] <file> [file ...]\n")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\navailable feature templates:")
		for _, f := range features {
			fmt.Fprintf(os.Stderr, "  %-30s  %s\n", f.Name, f.Desc)
		}
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	if len(featFlags) == 0 {
		fmt.Fprintln(os.Stderr, "no -feat specified; available templates:")
		for _, f := range features {
			fmt.Fprintf(os.Stderr, "  %-30s  %s\n", f.Name, f.Desc)
		}
		os.Exit(1)
	}

	selected := make([]struct {
		Name string
		Fn   FeatureFunc
	}, 0, len(featFlags))
	for _, name := range featFlags {
		fn, ok := featureByName(name)
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown feature %q\n", name)
			os.Exit(1)
		}
		selected = append(selected, struct {
			Name string
			Fn   FeatureFunc
		}{name, fn})
	}

	counts := make(countMap)
	for _, path := range flag.Args() {
		if err := processFile(path, selected, counts); err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", path, err)
			os.Exit(1)
		}
	}

	switch *fmt_ {
	case "matrix":
		printMatrix(counts, *minN)
	default:
		cov := *coverage
		if cov <= 0 {
			cov = 100
		}
		printHistogram(counts, *sortByCount, *minN, cov)
	}
}
