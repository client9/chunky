package tok

import "github.com/client9/chunky"

// LexicalRetag applies context-sensitive corrections based on capitalization.
// It runs per-sentence (inside Segment) so i==0 means sentence-initial.
func LexicalRetag(tokens []Token) []Token {
	for i, t := range tokens {
		if len(t.Word) == 0 || t.Word[0] < 'A' || t.Word[0] > 'Z' {
			continue
		}
		if i == 0 {
			// Sentence-initial: lexicon-tagged words ("The", "A") keep their tag.
			// Words resolved by morphology/inflection ("Ted" → VERB from -ed) → PROPN.
			if t.Rule != "lexicon" {
				tokens[i].Canidates = []chunky.Tag{chunky.TagPROPN}
				tokens[i].Rule = t.Rule + "+caps"
			}
			continue
		}
		// Non-sentence-initial: any known capitalized word → PROPN.
		if !t.IsUnknownTag() {
			tokens[i].Canidates = []chunky.Tag{chunky.TagPROPN}
			tokens[i].Rule = t.Rule + "+caps"
		}
	}
	return tokens
}
