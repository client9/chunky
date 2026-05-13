package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// Sentence is a sequence of tokens forming a single sentence, with the byte
// offset of the first token in the original source string.
type Sentence struct {
	Tokens []Token
	Offset int
}

// isAbbrevToken returns true if t is a known abbreviation that should not
// trigger a sentence boundary when followed by a period.
func isAbbrevToken(t Token) bool {
	// single uppercase letter: initial (C., A.)
	if len(t.Word) == 1 && t.Word[0] >= 'A' && t.Word[0] <= 'Z' {
		return true
	}
	// internal dot: U.S, A.M, e.g, i.e → abbreviation
	if strings.Contains(t.Word, ".") {
		return true
	}
	// ADV-tagged short abbreviations without dots: etc, vs
	for _, tag := range t.Canidates {
		if tag == chunky.TagADV {
			return true
		}
	}
	return false
}

// isBoundary returns true if the period-like PUNCT token at index i is a
// sentence boundary.
func isBoundary(tokens []Token, i int) bool {
	t := tokens[i]

	// only . ! ? end sentences
	if t.Word != "." && t.Word != "!" && t.Word != "?" {
		return false
	}

	// ellipsis: ".." is not a boundary
	if i+1 < len(tokens) && tokens[i+1].Word == "." {
		return false
	}
	if i > 0 && tokens[i-1].Word == "." {
		return false
	}

	// preceded by a known abbreviation → not a boundary
	if i > 0 && isAbbrevToken(tokens[i-1]) {
		return false
	}

	// PROPN . PROPN → middle initial, not a boundary
	if i > 0 && i+1 < len(tokens) &&
		len(tokens[i-1].Canidates) > 0 && tokens[i-1].Canidates[0] == chunky.TagPROPN &&
		len(tokens[i+1].Canidates) > 0 && tokens[i+1].Canidates[0] == chunky.TagPROPN {
		return false
	}

	return true
}

// Parse runs the full pipeline on s and returns tagged sentences.
// This is the primary entry point for callers; the pipeline order is an
// implementation detail that only this package needs to know.
func Parse(s string) []Sentence {
	return Segment(TagUnknowns(MergeCompounds(FilterBrackets(TagString(s)))))
}

// Segment splits a flat token slice into sentences. LexicalRetag is applied
// per-sentence so that sentence-initial capitalized words are handled correctly.
func Segment(tokens []Token) []Sentence {
	if len(tokens) == 0 {
		return nil
	}

	var sentences []Sentence
	start := 0

	flush := func(end int) {
		if end <= start {
			return
		}
		sent := tokens[start:end]
		// apply LexicalRetag per-sentence so i==0 is correctly sentence-initial
		sent = LexicalRetag(sent)
		sentences = append(sentences, Sentence{
			Tokens: sent,
			Offset: sent[0].Offset,
		})
		start = end
	}

	for i, t := range tokens {
		if len(t.Canidates) > 0 && t.Canidates[0] == chunky.TagPUNCT && isBoundary(tokens, i) {
			flush(i + 1)
		}
	}
	flush(len(tokens)) // trailing sentence with no terminal punctuation

	return sentences
}
