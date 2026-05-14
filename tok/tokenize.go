package tok

import (
	"unicode"
	"unicode/utf8"

	"github.com/client9/chunky"
)

// Token and Sentence are the canonical pipeline types defined in the root package.
type Token = chunky.Token
type Sentence = chunky.Sentence

// Tokenize splits s into tokens on Unicode whitespace, recording each token's
// byte offset into the original string. No normalization, filtering, or tagging
// is applied; each whitespace-delimited field becomes exactly one Token.
func Tokenize(s string) []Token {
	out := make([]Token, 0, 16)
	i := 0
	for i < len(s) {
		for i < len(s) {
			r, size := utf8.DecodeRuneInString(s[i:])
			if !unicode.IsSpace(r) {
				break
			}
			i += size
		}
		if i >= len(s) {
			break
		}
		start := i
		for i < len(s) {
			r, size := utf8.DecodeRuneInString(s[i:])
			if unicode.IsSpace(r) {
				break
			}
			i += size
		}
		out = append(out, Token{Word: s[start:i], Offset: start})
	}
	return out
}

// SurfaceTokenize returns the token strings produced by the pre-sentence
// pipeline (bracket stripping, normalization, punctuation splitting,
// contractions). Useful for callers that need surface forms without tagging.
func SurfaceTokenize(s string) []string {
	tokens := SplitPunctuation(NormalizeText(StripBrackets(Tokenize(s))))
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = t.Word
	}
	return out
}
