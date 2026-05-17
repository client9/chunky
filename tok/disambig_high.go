package tok

import "strings"

// DisambiguateHigh resolves "high" ({ADJ,ADV,NOUN}).
// ADJ dominates in every context (90–100%); no reliable NOUN or ADV trigger exists.
//
//   - next=NOUN|PROPN|ADJ|NUM|CCONJ|PUNCT → ADJ  (94–100%)
//   - prev=DET|ADP|ADV|AUX               → ADJ  (96–99%)
//   - next=ADP is left unresolved         (ADJ 65%, too mixed)
func DisambiguateHigh(tokens []Token) []Token {
	for i, t := range tokens {
		if strings.ToLower(t.Word) != "high" {
			continue
		}
		if !t.HasTag(TagADJ) {
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.HasTag(TagNOUN | TagPROPN | TagADJ | TagNUM | TagCCONJ | TagPUNCT):
			resolve = TagADJ
		case resolvedAs(prev, TagDET) || resolvedAs(prev, TagADP) || resolvedAs(prev, TagADV) || resolvedAs(prev, TagAUX):
			resolve = TagADJ
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+high"
		}
	}
	return tokens
}
