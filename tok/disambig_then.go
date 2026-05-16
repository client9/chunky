package tok

import "github.com/client9/chunky"

// DisambiguateThen resolves the ADJ/ADV/SCONJ ambiguity on "then" and "Then".
//
// Three uses exist, but only two are detectable with local context:
//   - ADJ  "the then president"  — directly follows a DET
//   - SCONJ "if X, then Y"       — requires a prior clause; not detectable locally
//   - ADV  "He left, then came"  — everything else (default)
//
// SCONJ is left as ADV: for chunking purposes both produce O (outside all
// chunks), and SCONJ "then" in the linter's prose target is rare.
func DisambiguateThen(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "then" && t.Word != "Then" {
			continue
		}
		if !t.HasTag(chunky.TagADJ) || !t.HasTag(chunky.TagADV) || !t.HasTag(chunky.TagSCONJ) {
			continue
		}
		tag := chunky.Tag(chunky.TagADV)
		if i > 0 && tokens[i-1].IsResolved() && tokens[i-1].Tags == chunky.TagDET {
			tag = chunky.TagADJ
		}
		tokens[i].Tags = tag
		tokens[i].Rule = t.Rule + "+then"
	}
	return tokens
}
