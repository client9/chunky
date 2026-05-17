package tok

import "strings"

// DisambiguateLike resolves the ADJ/ADP ambiguity on "like", "due",
// "pending", and "pursuant".
//
// like:
//   - prev=DET → ADJ  ("a like manner")
//   - next=DET|PRON|PROPN|NUM → ADP  ("like the idea", "like him")
//   - next=NOUN: left ambiguous (55/45 split in corpus)
//
// due: "due to" merged as compound; standalone due is almost always ADJ.
//   - next=NOUN|ADJ → ADJ  ("due date", "due diligence")
//   - prev=DET → ADJ  ("the due amount")
//
// pending:
//   - prev=DET|ADJ → ADJ  ("the pending case")
//   - prev=PRON|NOUN|PROPN (subject position) + next=NOUN → ADP  ("pending approval")
//
// pursuant: only prepositional in "pursuant to X"
//   - next.Word="to" → ADP
func DisambiguateLike(tokens []Token) []Token {
	for i := range tokens {
		disambiguateLike(tokens, i)
	}
	return tokens
}

func disambiguateLike(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADP) {
		return
	}
	lw := strings.ToLower(t.Word)
	prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
	var resolve Tag
	switch lw {
	case "like":
		switch {
		case resolvedAs(prev, TagDET):
			resolve = TagADJ
		case next.HasTag(TagDET | TagPRON | TagPROPN | TagNUM):
			resolve = TagADP
		}
	case "due":
		switch {
		case resolvedAs(prev, TagDET) || next.HasTag(TagNOUN|TagADJ|TagPROPN):
			resolve = TagADJ
		}
	case "pending":
		switch {
		case resolvedAs(prev, TagDET) || prev.HasTag(TagADJ):
			resolve = TagADJ
		case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && next.HasTag(TagNOUN|TagPROPN):
			resolve = TagADP // "she is pending approval"
		}
	case "pursuant":
		if strings.ToLower(next.Word) == "to" {
			resolve = TagADP
		}
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+like"
	}
}
