package tok

import (
	"strings"

	"github.com/client9/chunky"
)

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
	// known abbreviations (etc, vs, am, pm, …) that don't end sentences
	if _, ok := chunky.AbbreviationTags[strings.ToLower(t.Word)]; ok {
		return true
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
		tokens[i-1].HasTag(chunky.TagPROPN) &&
		tokens[i+1].HasTag(chunky.TagPROPN) {
		return false
	}

	return true
}

// Parse runs the full pipeline on s and returns tagged sentences.
func Parse(s string) []Sentence {
	tokens := Tokenize(s)
	tokens = StripBrackets(tokens)
	tokens = NormalizeText(tokens)
	tokens = SplitPunctuation(tokens)
	tokens = SplitContractions(tokens)
	tokens = chunky.MergeLexical(tokens)
	tokens = LexicalTag(tokens)
	tokens = TagUnknowns(tokens)
	tokens = DisambiguateApostropheS(tokens)
	tokens = DisambiguateThere(tokens)
	return sentencePhase(tokens)
}

func sentencePhase(tokens []Token) []Sentence {
	sents := Segment(tokens)
	for i := range sents {
		sents[i].Tokens = DisambiguateContext(sents[i].Tokens)
	}
	return sents
}

// Segment splits a flat token slice into sentences. RetagCapitalized is applied
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
		sent = RetagCapitalized(sent)
		sentences = append(sentences, Sentence{
			Tokens: sent,
			Offset: sent[0].Offset,
		})
		start = end
	}

	for i, t := range tokens {
		if t.HasTag(chunky.TagPUNCT) && isBoundary(tokens, i) {
			flush(i + 1)
		}
	}
	flush(len(tokens))

	return sentences
}
