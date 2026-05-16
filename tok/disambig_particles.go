package tok

import "strings"

// DisambiguateParticles resolves the ADP/ADV/NOUN ambiguity on "up" and "off".
//
// up:
//   - prev=VERB           → ADP  ("blow up", "pick up")
//   - next=DET|NOUN|PROPN|PRON|NUM → ADP  ("up the hill", "up to ten")
//   - next=PUNCT|CCONJ    → ADV  ("prices went up.", "up and running")
//   - next=ADV|ADJ        → ADV  ("up next", "up close")
//
// off:
//   - prev=VERB           → ADP  ("turn off", "cut off")
//   - next=DET|NOUN|PROPN|PRON → ADP  ("off the record", "off Broadway")
//   - next=PUNCT|CCONJ    → ADV  ("ran off.", "off and running")
func DisambiguateParticles(tokens []Token) []Token {
	for i, t := range tokens {
		lw := strings.ToLower(t.Word)
		if lw != "up" && lw != "off" {
			continue
		}
		if !t.HasTag(TagADP) || !t.HasTag(TagADV) {
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "up":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON | TagNUM):
				resolve = TagADP // "up the hill", "up to ten", "pick up the phone"
			case next.HasTag(TagPUNCT | TagCCONJ):
				resolve = TagADV // "went up.", "up and running" (no object → adverb)
			case next.HasTag(TagADV | TagADJ):
				resolve = TagADV // "up next", "up close"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADP // "blow up X", "pick up X" — transitive phrasal verb
			case next.HasTag(TagVERB | TagAUX | TagSCONJ):
				resolve = TagADV // "up comes", "up when", "what's up"
			}
		case "off":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON):
				resolve = TagADP // "off the record", "off Broadway"
			case next.HasTag(TagPUNCT | TagCCONJ | TagADV | TagADJ):
				resolve = TagADV // "went off.", "off and running", "off again"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADP // "turn off X", "cut off X" — transitive phrasal verb
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+particle"
		}
	}
	return tokens
}
