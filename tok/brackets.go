package tok

import "strings"

// StripBrackets removes bracketed noise from the token stream. Three cases:
//
//  1. Whole-token bracket ([1], [sic], {x}): token is removed entirely.
//  2. Embedded numeric citation ([digits] inside a longer word like "word[8]."):
//     the bracket span is replaced with spaces of equal byte length, preserving
//     the byte offsets of all surrounding characters.
//  3. Multi-token span ([critical section]): tokens are buffered until the
//     closing bracket is found and then discarded. Unclosed spans pass through.
func StripBrackets(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens))
	var buf []Token

	for _, t := range tokens {
		// Buffering a multi-token span.
		if buf != nil {
			if strings.HasSuffix(t.Word, "]") || strings.HasSuffix(t.Word, "}") {
				buf = nil // found close — discard span
			} else {
				buf = append(buf, t)
			}
			continue
		}

		open := strings.HasPrefix(t.Word, "[") || strings.HasPrefix(t.Word, "{")
		close := strings.HasSuffix(t.Word, "]") || strings.HasSuffix(t.Word, "}")

		// Whole-token bracket.
		if open && close {
			continue
		}

		// Embedded numeric citation(s): replace [digits] with spaces.
		// Checked before multi-token span so that "[1]." is handled here
		// rather than starting a span that never closes.
		if w := replaceEmbeddedCitations(t.Word); w != t.Word {
			out = append(out, Token{Word: w, Offset: t.Offset})
			continue
		}

		// Multi-token span start: starts with bracket, no embedded citation resolved.
		if open {
			buf = []Token{t}
			continue
		}

		out = append(out, t)
	}

	// Unclosed bracket — pass buffered tokens through.
	return append(out, buf...)
}

// replaceEmbeddedCitations replaces [digits] spans within s with spaces of the
// same byte length. This preserves byte offsets of all surrounding characters,
// allowing SplitPunctuation to correctly locate trailing punctuation.
func replaceEmbeddedCitations(s string) string {
	changed := false
	for {
		i := strings.LastIndexByte(s, '[')
		if i < 0 {
			break
		}
		j := strings.IndexByte(s[i:], ']')
		if j < 0 {
			break
		}
		inner := s[i+1 : i+j]
		if len(inner) == 0 {
			break
		}
		allDigits := true
		for _, c := range inner {
			if c < '0' || c > '9' {
				allDigits = false
				break
			}
		}
		if !allDigits {
			break
		}
		// Replace [digits] (j+1 bytes) with spaces of the same length.
		s = s[:i] + strings.Repeat(" ", j+1) + s[i+j+1:]
		changed = true
	}
	if !changed {
		return s
	}
	return s
}
