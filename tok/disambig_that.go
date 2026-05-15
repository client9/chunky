package tok

import "github.com/client9/chunky"

// DisambiguateThat resolves the PRON/SCONJ/DET ambiguity on "that" and "That".
//
// Only the most reliable case is handled: "that" directly before a DET
// article is the complementizer SCONJ ("He said that the car...").
// Other uses (DET "that car", PRON "after that") require wider context
// and are left for downstream rules.
func DisambiguateThat(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "that" && t.Word != "That" {
			continue
		}
		if !t.HasTag(chunky.TagPRON) || !t.HasTag(chunky.TagSCONJ) || !t.HasTag(chunky.TagDET) {
			continue
		}
		if i+1 >= len(tokens) || len(tokens[i+1].Tags) == 0 {
			continue
		}
		if tokens[i+1].Tags[0] != chunky.TagDET {
			continue
		}
		tokens[i].Tags = []chunky.Tag{chunky.TagSCONJ}
		tokens[i].Rule = t.Rule + "+that"
	}
	return tokens
}
