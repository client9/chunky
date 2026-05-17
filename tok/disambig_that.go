package tok

// DisambiguateThat resolves the PRON/SCONJ/DET ambiguity on "that" and "That".
//
// Resolved as SCONJ (complementizer):
//   - prev=VERB   → "said that", "argued that", "knowing that"
//   - prev=ADJ    → "confident that", "sure that", "aware that"
//   - next=DET    → "that the/a/an ..." (includes still-ambiguous DET like "more")
//   - prev=NOUN, next=PRON → appositive clause: "the fact that he/she/it/they..."
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
		case next.HasTag(TagDET):
			// "that the/a/more ..." — next has DET anywhere in its candidate set;
			// covers still-ambiguous words like "more" ({ADV,DET}) before DisambiguateMore runs
			resolve = TagSCONJ
		case resolvedAs(prev, TagNOUN) && resolvedAs(next, TagPRON):
			// Appositive complement clause: "the fact that he/she/it/they ..."
			// NOUN + that + subject-pronoun is reliably a complementizer, not a relative.
			resolve = TagSCONJ
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+that"
		}
	}
	return tokens
}
