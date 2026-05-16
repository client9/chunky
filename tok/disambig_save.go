package tok

import "strings"

// DisambiguateSave resolves {ADP,VERB} for "save" and "respecting".
//
// save:
//   - next=DET|NOUN|PRON|PROPN → VERB  ("save the file", "save lives")
//   - next=ADP(="for")         → left ambiguous  ("save for one exception")
//
// respecting:
//   - next=DET|NOUN|PRON|PROPN → ADP  ("respecting the law", formal preposition)
//   - prev=VERB|AUX             → VERB ("is respecting the rules")
func DisambiguateSave(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADP) || !t.HasTag(TagVERB) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "save", "respecting":
		default:
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "save":
			switch {
			case next.HasTag(TagDET | TagNOUN | TagPRON | TagPROPN):
				resolve = TagVERB
			case next.HasTag(TagPUNCT):
				resolve = TagVERB // imperative: "Save!"
			}
		case "respecting":
			switch {
			case resolvedAs(prev, TagVERB) || resolvedAs(prev, TagAUX):
				resolve = TagVERB // "is respecting the rules" (prev unambiguously VERB)
			case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON):
				resolve = TagADP
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+save"
		}
	}
	return tokens
}
