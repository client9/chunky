package tok

// DisambiguateLike resolves the ADJ/ADP ambiguity on "like" and "Like".
//
// Prepositional use (ADP): "like him", "like the idea", "like London"
// Adjectival use (ADJ):    "a like situation", "in like fashion" (rare)
//
// Resolved cases (corpus precision ~100%):
//   - prev=DET → ADJ  ("a like manner")
//   - next=DET|PRON|PROPN|NUM → ADP
//
// next=NOUN is left ambiguous: "like music" (ADP) and "like situation" (ADJ)
// split roughly 55/45 in corpus.
func DisambiguateLike(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "like" && t.Word != "Like" {
			continue
		}
		if !t.HasTag(TagADJ) || !t.HasTag(TagADP) {
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case resolvedAs(prev, TagDET):
			resolve = TagADJ
		case next.HasTag(TagDET | TagPRON | TagPROPN | TagNUM):
			resolve = TagADP
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+like"
		}
	}
	return tokens
}
