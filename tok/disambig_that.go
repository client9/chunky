package tok

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
		if !t.HasTag(TagPRON) || !t.HasTag(TagSCONJ) || !t.HasTag(TagDET) {
			continue
		}
		if resolvedAs(tokenAt(tokens, i+1), TagDET) {
			tokens[i].Tags = TagSCONJ
			tokens[i].Rule = t.Rule + "+that"
		}
	}
	return tokens
}
