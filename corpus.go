package chunky

import "strings"

// ParseTaggedLine parses a "word/TAG ..." corpus line into tokens.
// Tags[0] holds the corpus-assigned tag; Tags is empty if the tag string
// is unrecognized. Offset is not set (corpus lines carry no byte positions).
func ParseTaggedLine(line string) []Token {
	fields := strings.Fields(line)
	tokens := make([]Token, 0, len(fields))
	for _, f := range fields {
		i := strings.LastIndex(f, "/")
		if i <= 0 || i == len(f)-1 {
			continue
		}
		word := f[:i]
		t := Token{Word: word}
		if tag, err := ParseTag(f[i+1:]); err == nil {
			t.Tags = tag
		}
		tokens = append(tokens, t)
	}
	return tokens
}
