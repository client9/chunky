// eval-pos scores POS tagging accuracy against a CoNLL-2000 formatted file.
// Each line is "word PennPOS IOB-chunk"; blank lines separate sentences.
//
// Penn tags are mapped to UD via TagFromPennTag before comparison. Tags with a
// unique UD equivalent (NN→NOUN, VB→VERB, …) are scored as "unambiguous". Tags
// that map to more than one UD tag (IN→ADP|SCONJ, TO→PART|ADP, WDT→DET|PRON)
// are scored separately as "compatible" — correct if our prediction is any
// valid rendering.
//
// Usage:
//
//	go run ./cmd/eval-pos data/conll2000-test.txt
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/bits"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky"
	"github.com/client9/chunky/tok"
)

type tagStats struct {
	total   int
	correct int
	// top errors: our predicted tag → count
	errors map[string]int
}

func newTagStats() *tagStats { return &tagStats{errors: map[string]int{}} }

func (s *tagStats) acc() float64 {
	if s.total == 0 {
		return 0
	}
	return float64(s.correct) / float64(s.total) * 100
}

type conllToken struct{ word, pennPos string }

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
		if len(fields) < 2 {
			continue
		}
		current = append(current, conllToken{word: fields[0], pennPos: fields[1]})
	}
	if len(current) > 0 {
		sentences = append(sentences, current)
	}
	return sentences, sc.Err()
}

func tagWords(words []string, useChunk bool) []tok.Token {
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
		if useChunk {
			tokens = tok.Chunk(tokens)
			tokens = tok.DisambiguateByChunk(tokens)
		}
		if tok.TagsEqual(tokens, prev) {
			break
		}
	}
	return tokens
}

// isInLexicon reports whether word has a lexicon entry (before unknown tagging).
func isInLexicon(word string) bool {
	tokens := []tok.Token{{Word: word}}
	tokens = tok.LexicalTag(tokens)
	return !tokens[0].IsUnknownTag()
}

// chunkDiff runs both pipelines and reports what DisambiguateByChunk does to
// ambiguous Penn tokens — whether each chunk-driven resolution is a gain or loss.
func chunkDiff(sentences [][]conllToken) {
	type key struct{ penn, withChunk, noChunk string }
	gains := map[key]int{}
	losses := map[key]int{}

	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, t := range sent {
			words[i] = t.word
		}
		withChunk := tagWords(words, true)
		withoutChunk := tagWords(words, false)

		for i, gt := range sent {
			goldBits := chunky.TagFromPennTag(gt.pennPos)
			if bits.OnesCount32(uint32(goldBits)) < 2 {
				continue // only care about ambiguous Penn tags
			}
			wc := withChunk[i]
			nc := withoutChunk[i]
			if wc.Tags == nc.Tags {
				continue // chunk changed nothing here
			}
			// chunk changed the prediction — classify as gain or loss
			wcOK := wc.IsResolved() && goldBits&wc.Tags != 0
			ncOK := (nc.IsResolved() && goldBits&nc.Tags != 0) || (!nc.IsResolved() && goldBits&nc.Tags != 0)

			var wcStr, ncStr string
			if wc.IsResolved() {
				wcStr = wc.Tags.String()
			} else {
				wcStr = "{ambig}"
			}
			if nc.IsResolved() {
				ncStr = nc.Tags.String()
			} else {
				ncStr = "{ambig}"
			}
			k := key{gt.pennPos, wcStr, ncStr}
			if wcOK && !ncOK {
				gains[k]++
			} else if !wcOK && ncOK {
				losses[k]++
			}
		}
	}

	type row struct {
		k key
		n int
	}
	printSection := func(title string, m map[key]int) {
		var rows []row
		for k, n := range m {
			rows = append(rows, row{k, n})
		}
		sort.Slice(rows, func(i, j int) bool { return rows[i].n > rows[j].n })
		fmt.Printf("\n%s:\n", title)
		fmt.Printf("  %-6s  %-12s  %-12s  %s\n", "Penn", "with-chunk", "no-chunk", "count")
		fmt.Println("  " + strings.Repeat("-", 50))
		for _, r := range rows {
			fmt.Printf("  %-6s  %-12s  %-12s  %d\n", r.k.penn, r.k.withChunk, r.k.noChunk, r.n)
		}
	}
	printSection("GAINS (chunk correct, no-chunk wrong)", gains)
	printSection("LOSSES (chunk wrong, no-chunk correct)", losses)
}

// sampleErrors prints up to n examples where gold==wantGold and pred==wantPred,
// with surrounding context. wantPred may be "{ambig}" to match unresolved tokens.
func sampleErrors(sentences [][]conllToken, wantGold, wantPred string, n int) {
	goldTag, err := chunky.ParseTag(wantGold)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad gold tag %q: %v\n", wantGold, err)
		return
	}
	shown := 0
	type wordEntry struct{ word, penn, pred string }
	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, t := range sent {
			words[i] = t.word
		}
		tagged := tagWords(words, true)
		for i, gt := range sent {
			goldBits := chunky.TagFromPennTag(gt.pennPos)
			if goldBits != goldTag {
				continue
			}
			pred := tagged[i]
			var predStr string
			if pred.IsResolved() {
				predStr = pred.Tags.String()
			} else {
				predStr = "{ambig}"
			}
			if predStr != wantPred {
				continue
			}
			// print context window
			start := i - 3
			if start < 0 {
				start = 0
			}
			end := i + 4
			if end > len(sent) {
				end = len(sent)
			}
			var ctx []string
			for j := start; j < end; j++ {
				w := sent[j].word + "/" + sent[j].pennPos
				if j == i {
					w = "[" + w + "→" + predStr + "]"
				}
				ctx = append(ctx, w)
			}
			fmt.Printf("  %s\n", strings.Join(ctx, " "))
			shown++
			if shown >= n {
				return
			}
		}
	}
}

// wordErrors prints a frequency table of words involved in a specific gold→pred error.
func wordErrors(sentences [][]conllToken, wantGold, wantPred string) {
	goldTag, err := chunky.ParseTag(wantGold)
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad gold tag %q: %v\n", wantGold, err)
		return
	}
	counts := map[string]int{}
	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, t := range sent {
			words[i] = t.word
		}
		tagged := tagWords(words, true)
		for i, gt := range sent {
			goldBits := chunky.TagFromPennTag(gt.pennPos)
			if goldBits != goldTag {
				continue
			}
			pred := tagged[i]
			var predStr string
			if pred.IsResolved() {
				predStr = pred.Tags.String()
			} else {
				predStr = "{ambig}"
			}
			if predStr != wantPred {
				continue
			}
			counts[strings.ToLower(gt.word)]++
		}
	}
	type row struct {
		word string
		n    int
	}
	var rows []row
	for w, n := range counts {
		rows = append(rows, row{w, n})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].n > rows[j].n })
	fmt.Printf("  %-20s  %s\n", "word", "count")
	fmt.Println("  " + strings.Repeat("-", 30))
	for k, r := range rows {
		if k >= 30 {
			break
		}
		fmt.Printf("  %-20s  %d\n", r.word, r.n)
	}
}

func main() {
	topN := flag.Int("top", 10, "number of per-tag error patterns to show")
	noChunk := flag.Bool("nochunk", false, "disable Chunk+DisambiguateByChunk from the fixed-point loop")
	chunkDiffFlag := flag.Bool("chunkdiff", false, "show per-token impact of chunker on ambiguous Penn tags")
	sampleGold := flag.String("gold", "", "sample errors: expected gold tag (e.g. PROPN)")
	samplePred := flag.String("pred", "", "sample errors: predicted tag (e.g. NOUN or {ambig})")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: eval-pos [-top N] [-nochunk] [-chunkdiff] [-gold TAG -pred TAG] <conll2000-file>")
		os.Exit(1)
	}

	sentences, err := readConll(args[0])
	if err != nil {
		log.Fatalf("read: %v", err)
	}

	if *chunkDiffFlag {
		chunkDiff(sentences)
		return
	}

	if *sampleGold != "" && *samplePred != "" {
		fmt.Printf("Word frequency for gold=%s pred=%s:\n", *sampleGold, *samplePred)
		wordErrors(sentences, *sampleGold, *samplePred)
		fmt.Printf("\nExamples (gold=%s pred=%s):\n", *sampleGold, *samplePred)
		sampleErrors(sentences, *sampleGold, *samplePred, 20)
		return
	}

	type bucket struct {
		total, correct, ambigPred int
	}
	var (
		unamb   bucket
		amb     bucket
		inLex   bucket
		oov     bucket
		skipped int
	)

	// perTag tracks stats for each UD tag that appears as an unambiguous gold target.
	perTag := map[chunky.Tag]*tagStats{}

	for _, sent := range sentences {
		words := make([]string, len(sent))
		for i, t := range sent {
			words[i] = t.word
		}
		tagged := tagWords(words, !*noChunk)

		for i, gt := range sent {
			goldBits := chunky.TagFromPennTag(gt.pennPos)
			if goldBits == 0 {
				skipped++
				continue
			}

			pred := tagged[i]
			predResolved := pred.IsResolved()
			predTag := pred.Tags

			isAmbig := bits.OnesCount32(uint32(goldBits)) > 1

			// compatible: predTag is one of the valid UD renderings
			compatible := predResolved && goldBits&predTag != 0
			// also accept if pred is still ambiguous but contains a valid gold tag
			if !compatible && !predResolved && goldBits&predTag != 0 {
				compatible = true
			}

			if isAmbig {
				amb.total++
				if compatible {
					amb.correct++
				}
				if !predResolved {
					amb.ambigPred++
				}
			} else {
				// Unambiguous: goldBits has exactly one bit set.
				goldTag := goldBits
				correct := predResolved && predTag == goldTag

				unamb.total++
				if correct {
					unamb.correct++
				}
				if !predResolved {
					unamb.ambigPred++
				}

				lex := isInLexicon(gt.word)
				if lex {
					inLex.total++
					if correct {
						inLex.correct++
					}
				} else {
					oov.total++
					if correct {
						oov.correct++
					}
				}

				if perTag[goldTag] == nil {
					perTag[goldTag] = newTagStats()
				}
				ps := perTag[goldTag]
				ps.total++
				if correct {
					ps.correct++
				} else {
					var predStr string
					if predResolved {
						predStr = predTag.String()
					} else {
						predStr = "{ambig}"
					}
					ps.errors[predStr]++
				}
			}
		}
	}

	total := unamb.total + amb.total
	totalCorrect := unamb.correct + amb.correct

	pct := func(num, den int) string {
		if den == 0 {
			return "  n/a"
		}
		return fmt.Sprintf("%5.1f%%", float64(num)/float64(den)*100)
	}

	fmt.Printf("POS accuracy — CoNLL-2000 test (%d tokens, %d sentences)\n\n",
		total+skipped, len(sentences))
	fmt.Printf("%-30s  %7s  %7s  %7s  %7s\n", "", "tokens", "correct", "acc", "ambig")
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("%-30s  %7d  %7d  %7s  %6d\n", "Overall (unambig+compat)",
		total, totalCorrect, pct(totalCorrect, total), unamb.ambigPred+amb.ambigPred)
	fmt.Println(strings.Repeat("-", 65))
	fmt.Printf("%-30s  %7d  %7d  %7s  %6d\n", "Unambiguous Penn tags",
		unamb.total, unamb.correct, pct(unamb.correct, unamb.total), unamb.ambigPred)
	fmt.Printf("%-30s  %7d  %7d  %7s\n", "  In-lexicon",
		inLex.total, inLex.correct, pct(inLex.correct, inLex.total))
	fmt.Printf("%-30s  %7d  %7d  %7s\n", "  OOV",
		oov.total, oov.correct, pct(oov.correct, oov.total))
	fmt.Printf("%-30s  %7d  %7d  %7s  %6d\n", "Ambig Penn (IN/TO/WDT/VB*/JJ)",
		amb.total, amb.correct, pct(amb.correct, amb.total), amb.ambigPred)
	if skipped > 0 {
		fmt.Printf("Skipped (unmapped Penn tags): %d\n", skipped)
	}

	// Per-tag breakdown, sorted by token count descending.
	type tagRow struct {
		tag  chunky.Tag
		name string
		s    *tagStats
	}
	var rows []tagRow
	for tag, s := range perTag {
		if s.total > 0 {
			rows = append(rows, tagRow{tag, tag.String(), s})
		}
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].s.total > rows[j].s.total })

	fmt.Printf("\nPer-tag accuracy (unambiguous Penn tags only):\n")
	fmt.Printf("  %-8s  %7s  %7s  %7s  top errors\n", "tag", "tokens", "correct", "acc")
	fmt.Println("  " + strings.Repeat("-", 80))
	for _, r := range rows {
		type errEntry struct {
			pred  string
			count int
		}
		var errs []errEntry
		for pred, count := range r.s.errors {
			errs = append(errs, errEntry{pred, count})
		}
		sort.Slice(errs, func(i, j int) bool { return errs[i].count > errs[j].count })

		var errParts []string
		for k, e := range errs {
			if k >= *topN {
				break
			}
			errParts = append(errParts, fmt.Sprintf("%s×%d", e.pred, e.count))
		}
		errStr := strings.Join(errParts, "  ")
		fmt.Printf("  %-8s  %7d  %7d  %7s  %s\n",
			r.name, r.s.total, r.s.correct, pct(r.s.correct, r.s.total), errStr)
	}
}
