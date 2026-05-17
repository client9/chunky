package tok

import "strings"

// DisambiguateLong resolves "long" ({ADJ,ADV,VERB}).
//
//   - next=NOUN|PROPN|NUM → ADJ  ("a long time", "long road")
//   - prev=DET             → ADJ  ("the long", "a long")
//   - next=AUX|ADV|SCONJ  → ADV  ("how long will", "so long ago", "as long as")
//   - next=VERB            → ADV  ("how long waited", "not long after")
func DisambiguateLong(tokens []Token) []Token {
	for i, t := range tokens {
		if strings.ToLower(t.Word) != "long" {
			continue
		}
		if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.HasTag(TagNOUN | TagPROPN | TagNUM):
			resolve = TagADJ
		case resolvedAs(prev, TagDET):
			resolve = TagADJ
		case next.HasTag(TagAUX | TagADV | TagSCONJ):
			resolve = TagADV
		case next.HasTag(TagVERB) && !next.HasTag(TagAUX):
			resolve = TagADV
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+long"
		}
	}
	return tokens
}
