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
	return tokens
}

// isInLexicon reports whether word has a lexicon entry (before unknown tagging).
func isInLexicon(word string) bool {
	tokens := []tok.Token{{Word: word}}
	tokens = tok.LexicalTag(tokens)
	return !tokens[0].IsUnknownTag()
}

func main() {
	topN := flag.Int("top", 10, "number of per-tag error patterns to show")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "usage: eval-pos [-top N] <conll2000-file>")
		os.Exit(1)
	}

	sentences, err := readConll(args[0])
	if err != nil {
		log.Fatalf("read: %v", err)
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
		tagged := tagWords(words)

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
