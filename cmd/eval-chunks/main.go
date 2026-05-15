// eval-chunks scores the chunker against a CoNLL-2000 formatted file.
// Each line is "word PennPOS IOB-chunk"; blank lines separate sentences.
// Only NP, VP, and PP are scored; other chunk types are treated as O.
//
// Usage:
//
//	go run ./cmd/eval-chunks data/conll2000-test.txt
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/client9/chunky"
	"github.com/client9/chunky/tok"
)

type conllToken struct {
	word string
	gold chunky.ChunkTag
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
		rawChunk := fields[2]
		ct, err := chunky.ParseChunkTag(rawChunk)
		if err != nil {
			// Unknown type (SBAR, PRT, CONJP, …) → treat as O.
			ct = chunky.ChunkTag{IOB: 'O'}
		}
		current = append(current, conllToken{word: word, gold: ct})
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
	tokens = tok.DisambiguateContext(tokens)
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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: eval-chunks <conll2000-file>")
		os.Exit(1)
	}

	sentences, err := readConll(os.Args[1])
	if err != nil {
		log.Fatalf("read: %v", err)
	}

	totals := map[chunky.ChunkKind]*counts{
		chunky.ChunkNP: {},
		chunky.ChunkVP: {},
		chunky.ChunkPP: {},
	}
	overall := &counts{}

	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, ct := range sent {
			words[i] = ct.word
		}
		tagged := tagWords(words)

		gold := goldSpans(sent)
		pred := predSpans(tagged)

		// Index gold spans for O(1) lookup.
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
				}
			}
		}
		for _, s := range gold {
			if c, ok := totals[s.kind]; ok {
				if !predSet[s] {
					c.fn++
					overall.fn++
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
}
