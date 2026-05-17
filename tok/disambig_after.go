package tok

import "strings"

// subjectPronouns are nominative pronouns that introduce finite clauses.
// "after he/she/they/we/I/who" → SCONJ; object pronouns → ADP.
var subjectPronouns = map[string]bool{
	"i": true, "he": true, "she": true, "we": true,
	"they": true, "who": true, "whoever": true,
}

// DisambiguateAfter resolves the ADP/SCONJ ambiguity on temporal conjunctions.
//
// "after", "before", "until" lexicon: ADP|SCONJ
// ("till" is excluded — it carries an additional VERB tag that conflicts with NOUN rules.)
//
// ADP (takes an NP or gerund phrase):    "after the war", "before 2000", "until completing"
// SCONJ (introduces a finite clause):    "after the war ended", "before she arrived"
//
// Resolved cases:
//   - next=subjectPronoun → SCONJ  ("after he left", "before they arrived")
//   - next=VERB|AUX → ADP  (gerund/participle: "after graduating", "after being named")
//   - next=NUM      → ADP  (date or quantity: "after 1945", "until 2025")
//   - next=NOUN|PROPN → ADP  ("after midnight", "before the war", "until victory")
//   - next=PRON (object, not subject) → ADP  ("after him", "before them")
func DisambiguateAfter(tokens []Token) []Token {
	for i := range tokens {
		disambiguateAfter(tokens, i)
	}
	return tokens
}

func disambiguateAfter(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADP) || !t.HasTag(TagSCONJ) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case subjectPronouns[strings.ToLower(next.Word)]:
		resolve = TagSCONJ // "after he/she/they/we/I/who" — subject introduces clause
	case next.HasTag(TagVERB | TagAUX):
		resolve = TagADP // gerund/participle: "after graduating", "after being named"
	case next.HasTag(TagNUM):
		resolve = TagADP // "after 1945", "until 2025", "before 3 pm"
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagADP // "after midnight", "before the war", "until victory"
	case next.HasTag(TagPRON) && !subjectPronouns[strings.ToLower(next.Word)]:
		resolve = TagADP // "after him", "before them", "until it"
	case next.HasTag(TagPUNCT):
		resolve = TagADP // sentence-final: "after.", rarely SCONJ
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+after"
	}
}
