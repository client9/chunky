package tok

// DisambiguateUp resolves "up" tokens that carry a NOUN bit (Brown corpus noise).
// For ADP|ADV tokens without the NOUN bit, see disambig_particles.go.
// Context distinguishes ADP from ADV.
//
//   - next=DET|NOUN|PROPN|NUM|PRON → ADP  ("up the hill", "up one floor")
//   - next=PART                    → ADP  ("up to", "up for")
//   - prev=VERB (not AUX)          → ADP  ("picked up", "grew up", "scaled up")
//   - prev=NOUN|PRON               → ADP  ("cost up", "he up")
//   - prev=AUX                     → ADV  ("can't keep up", "will end up")
func DisambiguateUp(tokens []Token) []Token {
	for i := range tokens {
		disambiguateUp(tokens, i)
	}
	return tokens
}

func disambiguateUp(tokens []Token, i int) {
	t := tokens[i]
	if t.Word != "up" && t.Word != "Up" {
		return
	}
	if !t.HasTag(TagADP) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagDET | TagNOUN | TagPROPN | TagNUM | TagPRON):
		resolve = TagADP
	case next.HasTag(TagPART):
		resolve = TagADP // "up to", "up for"
	case resolvedAs(prev, TagAUX):
		resolve = TagADV
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADP // phrasal particle: "pick up", "scale up"
	case resolvedAs(prev, TagNOUN) || resolvedAs(prev, TagPRON):
		resolve = TagADP
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+up"
	}
}
