package tok

import "strings"

// DisambiguateAs resolves the ADP/ADV/SCONJ ambiguity on "as" and "As".
//
// ADP ("as a role", "known as the leader", "as of January", "such as"):
//   - next=DET and prev≠PUNCT  ("as a/the/an ...")
//   - next.Word == "of"        ("as of [date]")
//   - prev.Word == "well"      ("well as the" — middle of "as well as")
//
// ADV ("as well", "as yet"):
//   - prev=PUNCT and next=ADV  (", as well", ", as yet")
//
// SCONJ ("as it turned out", ", as he said"):
//   - prev=PUNCT and next=PRON (", as it/he/she/they ...")
func DisambiguateAs(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "as" && t.Word != "As" {
			continue
		}
		if !t.HasTag(TagADP) || !t.HasTag(TagSCONJ) {
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case next.Word == "of":
			resolve = TagADP // "as of January", "as of now"
		case resolvedAs(prev, TagPUNCT) && next.HasTag(TagADV):
			resolve = TagADV // ", as well", ", as yet"
		case resolvedAs(prev, TagPUNCT) && resolvedAs(next, TagPRON):
			resolve = TagSCONJ // ", as it turned out", ", as he said"
		case next.HasTag(TagDET) && !resolvedAs(prev, TagPUNCT):
			resolve = TagADP // "served as the/a ...", "as a result", "such as the"
		case next.HasTag(TagNOUN|TagPROPN) && !resolvedAs(prev, TagPUNCT):
			resolve = TagADP // "served as chairman", "known as Smith"
		case strings.ToLower(prev.Word) == "well":
			resolve = TagADP // "well as the/a" — middle of "as well as"
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+as"
		}
	}
	return tokens
}
