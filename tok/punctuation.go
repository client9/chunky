package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// SplitPunctuation splits each token's leading and trailing punctuation into
// separate tokens and expands contractions. It handles internal spaces that
// may have been introduced by StripBrackets replacing embedded citations.
//
// For a token whose Word contains internal spaces, the offset of any character
// at index i within Word equals Token.Offset+i, because StripBrackets replaces
// bracket spans with spaces of equal byte length.
func SplitPunctuation(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens)+8)
	for _, t := range tokens {
		out = splitOneToken(out, t)
	}
	return splitContractions(out)
}

func splitOneToken(out []Token, t Token) []Token {
	p := t.Word
	pos := t.Offset

	// Trim leading spaces (produced by StripBrackets replacing an embedded
	// citation that started at the beginning of the field, e.g. "[1]word").
	for len(p) > 0 && p[0] == ' ' {
		p = p[1:]
		pos++
	}
	if len(p) == 0 {
		return out
	}
	if len(p) == 1 {
		return append(out, Token{Word: p, Offset: pos})
	}

	// Split leading '('.
	if p[0] == '(' {
		out = append(out, Token{Word: "(", Offset: pos})
		p = p[1:]
		pos++
	}
	if len(p) == 0 {
		return out
	}

	// Find the last non-space character to locate trailing punctuation.
	// Spaces may appear before the punctuation when an embedded citation
	// was replaced with spaces (e.g. "word   .").
	lastNonSpace := len(p) - 1
	for lastNonSpace >= 0 && p[lastNonSpace] == ' ' {
		lastNonSpace--
	}

	last, lastPos := "", 0
	if lastNonSpace >= 0 {
		ch := p[lastNonSpace]
		if ch == ',' || ch == '.' || ch == ':' || ch == ';' || ch == '!' || ch == '?' {
			// Keep the dot when the whole word (trimmed) is a dotted abbreviation.
			candidate := strings.TrimRight(p[:lastNonSpace+1], " ")
			if ch != '.' || !chunky.DottedAbbreviations[strings.ToLower(candidate)] {
				last = string(ch)
				lastPos = pos + lastNonSpace
				p = p[:lastNonSpace]
			}
		}
	}

	// Trim trailing spaces left after removing the punctuation character.
	p = strings.TrimRight(p, " ")

	if len(p) == 0 {
		if last != "" {
			out = append(out, Token{Word: last, Offset: lastPos})
		}
		return out
	}

	// Split trailing ')'.
	if p[len(p)-1] == ')' {
		out = append(out, Token{Word: p[:len(p)-1], Offset: pos})
		out = append(out, Token{Word: ")", Offset: pos + len(p) - 1})
		if last != "" {
			out = append(out, Token{Word: last, Offset: lastPos})
		}
		return out
	}

	out = append(out, Token{Word: p, Offset: pos})
	if last != "" {
		out = append(out, Token{Word: last, Offset: lastPos})
	}
	return out
}

// contractionSuffixes is the set of apostrophe-led suffixes that trigger a split.
var contractionSuffixes = map[string]bool{
	"'ll": true, "'re": true, "'ve": true, "'m": true, "'d": true, "'s": true,
	"'t": true,
}

// splitContractions expands contraction tokens into (stem, suffix) pairs.
// Irregular forms in chunky.ContractionNorm are handled first. Words in
// chunky.AbbreviationTags that are not in ContractionNorm stay whole (ain't).
func splitContractions(tokens []Token) []Token {
	out := make([]Token, 0, len(tokens)+4)
	for _, t := range tokens {
		lower := strings.ToLower(t.Word)

		// Irregular forms: won't → will + n't, shan't → shall + n't.
		if parts, ok := chunky.ContractionNorm[lower]; ok {
			out = append(out, Token{Word: parts[0], Offset: t.Offset})
			out = append(out, Token{Word: parts[1], Offset: t.Offset + len(t.Word) - len(parts[1])})
			continue
		}

		// Words that stay whole (ain't, o'clock, etc.).
		if _, ok := chunky.AbbreviationTags[lower]; ok {
			out = append(out, t)
			continue
		}

		ap := strings.IndexByte(t.Word, '\'')
		if ap <= 0 {
			out = append(out, t)
			continue
		}

		suffix := lower[ap:]

		// n't: move the n to the suffix when stem-without-n is in the lexicon
		// (don't→do+n't, shouldn't→should+n't); keep n in stem otherwise
		// (can't→can+'t, where "ca" is not a word).
		if suffix == "'t" && ap >= 2 && (t.Word[ap-1] == 'n' || t.Word[ap-1] == 'N') {
			stemNoN := strings.ToLower(t.Word[:ap-1])
			if len(wordtagmap[stemNoN]) > 0 || len(chunky.AbbreviationTags[stemNoN]) > 0 {
				out = append(out, Token{Word: t.Word[:ap-1], Offset: t.Offset})
				out = append(out, Token{Word: "n't", Offset: t.Offset + ap - 1})
			} else {
				out = append(out, Token{Word: t.Word[:ap], Offset: t.Offset})
				out = append(out, Token{Word: "'t", Offset: t.Offset + ap})
			}
			continue
		}

		if contractionSuffixes[suffix] {
			out = append(out, Token{Word: t.Word[:ap], Offset: t.Offset})
			out = append(out, Token{Word: t.Word[ap:], Offset: t.Offset + ap})
			continue
		}

		out = append(out, t)
	}
	return out
}
