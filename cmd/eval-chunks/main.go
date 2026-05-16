// eval-chunks scores the chunker against a CoNLL-2000 formatted file.
// Each line is "word PennPOS IOB-chunk"; blank lines separate sentences.
// Only NP, VP, and PP are scored; other chunk types are treated as O.
//
// Usage:
//
//	go run ./cmd/eval-chunks data/conll2000-test.txt
//	go run ./cmd/eval-chunks -errors NP data/conll2000-test.txt   # show FP/FN spans
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky"
	"github.com/client9/chunky/tok"
)

type conllToken struct {
	word    string
	pennPos string
	gold    chunky.ChunkTag
}

// readConll reads a CoNLL-2000 file into sentences. Each sentence is a slice
// of conllTokens. Unknown chunk types (SBAR, PRT, etc.) are silently mapped to O.
func readConll(path string) ([][]conllToken, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var sentences [][]conllToken
	var current []conllToken

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			if len(current) > 0 {
				sentences = append(sentences, current)
				current = nil
			}
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		word := fields[0]
		pennPos := fields[1]
		rawChunk := fields[2]
		ct, err := chunky.ParseChunkTag(rawChunk)
		if err != nil {
			ct = chunky.ChunkTag{IOB: 'O'}
		}
		current = append(current, conllToken{word: word, pennPos: pennPos, gold: ct})
	}
	if len(current) > 0 {
		sentences = append(sentences, current)
	}
	return sentences, sc.Err()
}

// tagWords runs the POS pipeline on a pre-tokenized word list, bypassing
// the tokenizer. Suitable for CoNLL-style input where tokenization is given.
func tagWords(words []string) []tok.Token {
	tokens := make([]tok.Token, len(words))
	for i, w := range words {
		tokens[i] = tok.Token{Word: w}
	}
	tokens = tok.LexicalTag(tokens)
	tokens = tok.TagUnknowns(tokens)
	tokens = tok.RetagCapitalized(tokens)
	tokens = tok.DisambiguateWords(tokens)
	for {
		prev := tok.CopyTags(tokens)
		tokens = tok.DisambiguateContext(tokens)
		tokens = tok.Chunk(tokens)
		tokens = tok.DisambiguateByChunk(tokens)
		if tok.TagsEqual(tokens, prev) {
			break
		}
	}
	tokens = tok.Chunk(tokens)
	return tokens
}

// span is a half-open [start, end) token range with a chunk kind.
type span struct {
	kind       chunky.ChunkKind
	start, end int
}

// goldSpans extracts chunk spans from CoNLL gold labels.
func goldSpans(tokens []conllToken) []span {
	var spans []span
	for i := 0; i < len(tokens); {
		ct := tokens[i].gold
		if ct.IOB != 'B' {
			i++
			continue
		}
		kind := ct.Kind
		j := i + 1
		for j < len(tokens) && tokens[j].gold.IOB == 'I' && tokens[j].gold.Kind == kind {
			j++
		}
		spans = append(spans, span{kind: kind, start: i, end: j})
		i = j
	}
	return spans
}

// predSpans extracts chunk spans from pipeline-tagged tokens.
func predSpans(tokens []tok.Token) []span {
	var spans []span
	for i := 0; i < len(tokens); {
		ct := tokens[i].Chunk
		if ct.IOB != 'B' {
			i++
			continue
		}
		kind := ct.Kind
		j := i + 1
		for j < len(tokens) && tokens[j].Chunk.IOB == 'I' && tokens[j].Chunk.Kind == kind {
			j++
		}
		spans = append(spans, span{kind: kind, start: i, end: j})
		i = j
	}
	return spans
}

type counts struct{ tp, fp, fn int }

func (c counts) precision() float64 {
	if c.tp+c.fp == 0 {
		return 0
	}
	return float64(c.tp) / float64(c.tp+c.fp)
}

func (c counts) recall() float64 {
	if c.tp+c.fn == 0 {
		return 0
	}
	return float64(c.tp) / float64(c.tp+c.fn)
}

func (c counts) f1() float64 {
	p, r := c.precision(), c.recall()
	if p+r == 0 {
		return 0
	}
	return 2 * p * r / (p + r)
}

// spanWords returns the space-joined words of a span from the sentence.
func spanWords(sent []conllToken, s span) string {
	parts := make([]string, s.end-s.start)
	for i, ct := range sent[s.start:s.end] {
		parts[i] = ct.word
	}
	return strings.Join(parts, " ")
}

// spanTags returns the space-joined POS tags of a span.
func spanTags(sent []conllToken, s span) string {
	parts := make([]string, s.end-s.start)
	for i, ct := range sent[s.start:s.end] {
		parts[i] = ct.pennPos
	}
	return strings.Join(parts, " ")
}

// spanPredTags returns the space-joined predicted tags for a span.
func spanPredTags(tagged []tok.Token, s span) string {
	parts := make([]string, s.end-s.start)
	for i, t := range tagged[s.start:s.end] {
		parts[i] = t.String()[len(t.Word)+1:] // strip "word/" prefix
	}
	return strings.Join(parts, " ")
}

type errorEntry struct {
	kind    string // "FP" or "FN"
	words   string
	goldPos string // Penn POS tags from gold data
	predPos string // our predicted tags
	count   int
}

func main() {
	errKind := flag.String("errors", "", "dump FP/FN spans for this chunk type (NP, VP, PP)")
	topN := flag.Int("top", 30, "number of error patterns to show")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: eval-chunks [-errors NP|VP|PP] [-top N] <conll2000-file>")
		os.Exit(1)
	}

	sentences, err := readConll(args[0])
	if err != nil {
		log.Fatalf("read: %v", err)
	}

	var filterKind chunky.ChunkKind
	doErrors := *errKind != ""
	if doErrors {
		switch strings.ToUpper(*errKind) {
		case "NP":
			filterKind = chunky.ChunkNP
		case "VP":
			filterKind = chunky.ChunkVP
		case "PP":
			filterKind = chunky.ChunkPP
		default:
			log.Fatalf("unknown chunk type %q; use NP, VP, or PP", *errKind)
		}
	}

	totals := map[chunky.ChunkKind]*counts{
		chunky.ChunkNP: {},
		chunky.ChunkVP: {},
		chunky.ChunkPP: {},
	}
	overall := &counts{}

	// error tallies: "FP|words" or "FN|words" → count
	fpByWords := map[string]*errorEntry{}
	fnByWords := map[string]*errorEntry{}

	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, ct := range sent {
			words[i] = ct.word
		}
		tagged := tagWords(words)

		gold := goldSpans(sent)
		pred := predSpans(tagged)

		goldSet := make(map[span]bool, len(gold))
		for _, s := range gold {
			goldSet[s] = true
		}
		predSet := make(map[span]bool, len(pred))
		for _, s := range pred {
			predSet[s] = true
		}

		for _, s := range pred {
			if c, ok := totals[s.kind]; ok {
				if goldSet[s] {
					c.tp++
					overall.tp++
				} else {
					c.fp++
					overall.fp++
					if doErrors && s.kind == filterKind {
						w := spanWords(sent, s)
						key := w
						if e, ok := fpByWords[key]; ok {
							e.count++
						} else {
							fpByWords[key] = &errorEntry{
								kind:    "FP",
								words:   w,
								goldPos: spanTags(sent, s),
								predPos: spanPredTags(tagged, s),
								count:   1,
							}
						}
					}
				}
			}
		}
		for _, s := range gold {
			if c, ok := totals[s.kind]; ok {
				if !predSet[s] {
					c.fn++
					overall.fn++
					if doErrors && s.kind == filterKind {
						w := spanWords(sent, s)
						key := w
						if e, ok := fnByWords[key]; ok {
							e.count++
						} else {
							fnByWords[key] = &errorEntry{
								kind:    "FN",
								words:   w,
								goldPos: spanTags(sent, s),
								predPos: spanPredTags(tagged, s),
								count:   1,
							}
						}
					}
				}
			}
		}
	}

	fmt.Printf("%-6s  %7s  %7s  %7s  %6s  %6s  %6s\n",
		"type", "TP", "FP", "FN", "P", "R", "F1")
	fmt.Println(strings.Repeat("-", 54))
	for _, kind := range []chunky.ChunkKind{chunky.ChunkNP, chunky.ChunkVP, chunky.ChunkPP} {
		c := totals[kind]
		fmt.Printf("%-6s  %7d  %7d  %7d  %5.1f%%  %5.1f%%  %5.1f%%\n",
			kind, c.tp, c.fp, c.fn,
			c.precision()*100, c.recall()*100, c.f1()*100)
	}
	fmt.Println(strings.Repeat("-", 54))
	fmt.Printf("%-6s  %7d  %7d  %7d  %5.1f%%  %5.1f%%  %5.1f%%\n",
		"TOTAL", overall.tp, overall.fp, overall.fn,
		overall.precision()*100, overall.recall()*100, overall.f1()*100)
	fmt.Printf("\n%d sentences evaluated\n", len(sentences))

	if !doErrors {
		return
	}

	// Sort and print FP patterns.
	fps := make([]*errorEntry, 0, len(fpByWords))
	for _, e := range fpByWords {
		fps = append(fps, e)
	}
	sort.Slice(fps, func(i, j int) bool { return fps[i].count > fps[j].count })

	fmt.Printf("\n=== %s FP (predicted but not gold) — top %d patterns ===\n", *errKind, *topN)
	fmt.Printf("%-6s  %-40s  %-30s  %s\n", "count", "words", "gold-pos", "pred-pos")
	fmt.Println(strings.Repeat("-", 110))
	for i, e := range fps {
		if i >= *topN {
			break
		}
		fmt.Printf("%6d  %-40s  %-30s  %s\n", e.count, trunc(e.words, 40), trunc(e.goldPos, 30), e.predPos)
	}

	// Sort and print FN patterns.
	fns := make([]*errorEntry, 0, len(fnByWords))
	for _, e := range fnByWords {
		fns = append(fns, e)
	}
	sort.Slice(fns, func(i, j int) bool { return fns[i].count > fns[j].count })

	fmt.Printf("\n=== %s FN (gold but not predicted) — top %d patterns ===\n", *errKind, *topN)
	fmt.Printf("%-6s  %-40s  %-30s  %s\n", "count", "words", "gold-pos", "pred-pos")
	fmt.Println(strings.Repeat("-", 110))
	for i, e := range fns {
		if i >= *topN {
			break
		}
		fmt.Printf("%6d  %-40s  %-30s  %s\n", e.count, trunc(e.words, 40), trunc(e.goldPos, 30), e.predPos)
	}
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}
