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
			// Sentence-initial capitalization is grammatical, not semantic.
			// Only promote to PROPN if the word got a non-noun tag from a
			// morphological rule — e.g. "Ted" tagged VERB by the -ed inflection.
			// Words already tagged NOUN (unk:word fallback, morph suffix) stay NOUN:
			// "Hantaviruses" starts a sentence but is not a proper noun.
			if t.Rule != "lexicon" && !t.IsUnknownTag() && t.Canidates[0] != chunky.TagNOUN {
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
