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
			// We cannot distinguish a proper noun from a sentence-initial verb or
			// participle (e.g. "Walked the dog." vs "Ted is here."), so we leave
			// the pipeline tag unchanged.
			continue
		}
		// Non-sentence-initial: any known capitalized word → PROPN,
		// except PRON ("I" is always capitalized in English).
		if !t.IsUnknownTag() && t.Tags[0] != chunky.TagPRON {
			tokens[i].Tags = []chunky.Tag{chunky.TagPROPN}
			tokens[i].Rule = t.Rule + "+caps"
		}
	}
	return tokens
}
