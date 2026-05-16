package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// contractionSuffixes is the set of apostrophe-led suffixes that trigger a split.
var contractionSuffixes = map[string]bool{
	"'ll": true, "'re": true, "'ve": true, "'m": true, "'d": true, "'s": true,
	"'t": true,
}

// SplitContractions expands contraction tokens into (stem, suffix) pairs.
// Irregular forms in chunky.ContractionNorm are handled first. Words in
// chunky.AbbreviationTags that are not in ContractionNorm stay whole (ain't).
// The n't split consults the lexicon to decide whether n belongs to the stem
// (can't → can+'t) or the suffix (don't → do+n't).
func SplitContractions(tokens []Token) []Token {
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
			if wordtagmap[stemNoN] != 0 || chunky.AbbreviationTags[stemNoN] != 0 {
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
