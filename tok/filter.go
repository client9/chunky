package tok

import "strings"

// stripInlineCitations removes Wikipedia-style numeric citation markers
// embedded within a field, e.g. "Planeteers.[8]" → "Planeteers."
// Only strips markers whose content is all digits. Stripping from the end
// preserves byte offsets for all characters that remain.
func stripInlineCitations(s string) string {
	for {
		i := strings.LastIndexByte(s, '[')
		if i <= 0 {
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
		s = s[:i] + s[i+j+1:]
	}
	return s
}

// FilterBrackets removes bracketed spans from the token stream.
// Both single-token forms ([1], [sic]) and multi-token spans
// ([critical section]) are deleted. Curly braces are handled the same way.
// Unclosed brackets pass through rather than consuming the rest of the stream.
func FilterBrackets(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens))
	var buf []Token

	for _, t := range tokens {
		if buf != nil {
			if strings.HasSuffix(t.Word, "]") || strings.HasSuffix(t.Word, "}") {
				buf = nil // found close — discard span
			} else {
				buf = append(buf, t)
			}
			continue
		}
		if strings.HasPrefix(t.Word, "[") || strings.HasPrefix(t.Word, "{") {
			if strings.HasSuffix(t.Word, "]") || strings.HasSuffix(t.Word, "}") {
				continue // single-token [1], [sic], {x}
			}
			buf = []Token{t} // start buffering multi-token span
			continue
		}
		out = append(out, t)
	}

	// unclosed bracket — pass buffered tokens through rather than losing them
	return append(out, buf...)
}
