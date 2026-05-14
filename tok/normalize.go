package tok

import "github.com/client9/typewriter"

var tw = typewriter.New(typewriter.Config{
	Categories: typewriter.Quotes | typewriter.Spaces,
})

// NormalizeText applies typographic normalization to each token's Word:
// curly quotes → straight, em/en dashes → ASCII, Unicode spaces → space.
// Runs after StripBrackets and before SplitPunctuation so that multi-character
// sequences (-- → —, curly quote pairs) are recognized on whole fields before
// punctuation is split off.
func NormalizeText(tokens []Token) []Token {
	for i, t := range tokens {
		if n := tw.Replace(t.Word); n != t.Word {
			tokens[i].Word = n
		}
	}
	return tokens
}
