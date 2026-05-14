package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// MergeLexical scans the token stream left-to-right and replaces sequences
// matching chunky.CompoundTags with a single token carrying the compound tag.
// Longest match wins. The merged token's Word is the original surface form
// (space-joined), and its Offset is taken from the first token in the sequence.
func MergeLexical(tokens []Token) []Token {
	if len(tokens) == 0 {
		return tokens
	}
	out := make([]Token, 0, len(tokens))
	i := 0
	for i < len(tokens) {
		merged := false
		for length := chunky.CompoundMaxLen; length >= 2; length-- {
			if i+length > len(tokens) {
				continue
			}
			lower := make([]string, length)
			surface := make([]string, length)
			for j := 0; j < length; j++ {
				lower[j] = strings.ToLower(tokens[i+j].Word)
				surface[j] = tokens[i+j].Word
			}
			key := strings.Join(lower, " ")
			if tag, ok := chunky.CompoundTags[key]; ok {
				out = append(out, Token{
					Word:       strings.Join(surface, " "),
					Offset:     tokens[i].Offset,
					Candidates: []chunky.Tag{tag},
					Rule:       "compound",
				})
				i += length
				merged = true
				break
			}
		}
		if !merged {
			out = append(out, tokens[i])
			i++
		}
	}
	return out
}
