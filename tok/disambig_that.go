package tok

// DisambiguateThat resolves the PRON/SCONJ/DET ambiguity on "that" and "That".
//
// Resolved as SCONJ (complementizer):
//   - prev=VERB   → "said that", "argued that", "knowing that"
//   - prev=ADJ    → "confident that", "sure that", "aware that"
//   - next=DET    → "that the/a/an ..." (existing rule)
//
// DET ("that car") and PRON ("after that") require wider context and are
// left for downstream rules.
func DisambiguateThat(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "that" && t.Word != "That" {
			continue
		}
		if !t.HasTag(TagPRON) || !t.HasTag(TagSCONJ) || !t.HasTag(TagDET) {
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case prev.HasTag(TagVERB) && next.HasTag(TagDET|TagPRON|TagNOUN|TagPROPN|TagADV|TagAUX):
			// "said that he/the/it/there..." — complementizer, not object pronoun
			resolve = TagSCONJ
		case resolvedAs(prev, TagADJ) && next.HasTag(TagDET|TagPRON|TagNOUN|TagPROPN|TagADV|TagAUX):
			// "confident that it/the..." — complementizer
			resolve = TagSCONJ
		case resolvedAs(next, TagDET):
			resolve = TagSCONJ
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+that"
		}
	}
	return tokens
}
