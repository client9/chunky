package tok

import "strings"

// DisambiguateDown resolves "down" ({ADP,ADV,PROPN}) and "near" ({ADJ,ADP,ADV,PROPN}).
// The PROPN bit is Brown corpus noise; context distinguishes ADP from ADV.
//
// down:
//   - next=DET|NOUN|PROPN|NUM → ADP  ("down the stairs", "down Main St.")
//   - prev=VERB (pure)        → ADV  ("went down", "broke down")
//
// near:
//   - next=DET|NOUN|PROPN     → ADP  ("near the station", "near London")
//   - next=ADJ|ADV            → ADV  ("near perfect", "near enough")
//   - prev=VERB (pure)        → ADV  ("standing near")
func DisambiguateDown(tokens []Token) []Token {
	for i := range tokens {
		disambiguateDown(tokens, i)
	}
	return tokens
}

func disambiguateDown(tokens []Token, i int) {
	t := tokens[i]
	lw := strings.ToLower(t.Word)
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "down":
		if !t.HasTag(TagADP) {
			return
		}
		switch {
		case next.HasTag(TagDET | TagNOUN | TagPROPN | TagPRON):
			resolve = TagADP // "down the stairs", "down her street"
		case next.HasTag(TagPART):
			resolve = TagADP // "down to", "down for"
		case resolvedAs(prev, TagAUX):
			resolve = TagADV // "is down", "will be down"
		case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
			resolve = TagADV
		case next.HasTag(TagPUNCT | TagADV | TagADJ | TagSCONJ):
			resolve = TagADV // "fell down.", "down low", "down as"
		case resolvedAs(prev, TagPUNCT):
			resolve = TagADV // ", down she went"
		}
	case "near":
		if !t.HasTag(TagADP) {
			return
		}
		switch {
		case next.HasTag(TagDET | TagNOUN | TagPROPN):
			resolve = TagADP
		case next.HasTag(TagADJ | TagADV):
			resolve = TagADV
		case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
			resolve = TagADV
		}
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+down"
	}
}
