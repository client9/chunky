package tok

import "github.com/client9/chunky"

// LexicalRetag applies context-sensitive corrections to tokens that already
// have lexicon tags. Unlike TagUnknowns, this operates on known words.
func LexicalRetag(tokens []Token) []Token {
	for i, t := range tokens {
		if t.IsUnknownTag() {
			continue
		}
		// Non-sentence-initial capitalized word → PROPN should be primary candidate.
		// The lexicon tags the base form (lowercase) so capitalization is ignored
		// at lookup time. A capitalized known word mid-sentence is almost certainly
		// a proper noun, even if the lowercase form is a common word.
		if i > 0 && len(t.Word) > 0 && t.Word[0] >= 'A' && t.Word[0] <= 'Z' {
			tokens[i].Canidates = []chunky.Tag{chunky.TagPROPN}
			tokens[i].Rule = t.Rule + "+caps"
		}
	}
	return tokens
}
