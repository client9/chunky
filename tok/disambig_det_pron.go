package tok

import "strings"

// DisambiguateDetPron resolves the DET/PRON ambiguity on floating quantifiers
// and demonstratives.
//
// Quantifiers (each, some, any):
//   DET: "each team", "some water", "any questions"
//   PRON: "they each have", "some of us", "each of them"
//
// Demonstratives (this, these, those):
//   DET:  "this decision", "these issues", "those teams"
//   PRON: "this is clear", "these are done", "of those"
//
// Resolved cases:
//   quantifiers:
//   - next=ADP                              → PRON ("each of them")
//   - prev=PRON|NOUN|PROPN + next=VERB|AUX → PRON (floating: "they each have")
//   - next=NOUN|PROPN (unambiguous)         → DET (prenominal)
//   - next=ADJ (unambiguous)                → DET ("each individual")
//   - resolvedAs(next, VERB)                → PRON ("some argue")
//
//   demonstratives:
//   - next=NOUN|ADJ|PROPN|NUM (unambiguous) → DET ("this decision")
//   - next=VERB|AUX                         → PRON ("this is clear")
//   - next=PUNCT|CCONJ                      → PRON ("this.", "these,")
//   - resolvedAs(prev, ADP) + no following noun → PRON ("of those")
func DisambiguateDetPron(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagDET) || !t.HasTag(TagPRON) {
			continue
		}
		lw := strings.ToLower(t.Word)
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "each", "some", "any":
			switch {
			case next.HasTag(TagADP):
				resolve = TagPRON // "each of them", "some of us", "any of these"
			case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && next.HasTag(TagVERB|TagAUX):
				resolve = TagPRON // floating quantifier: "they each have", "we each took"
			case next.HasTag(TagNOUN|TagPROPN) && !next.HasTag(TagVERB):
				resolve = TagDET // "each team", "some water", "any questions"
			case next.HasTag(TagADJ) && !next.HasTag(TagVERB|TagNOUN):
				resolve = TagDET // "each individual", "any further"
			case resolvedAs(next, TagVERB):
				resolve = TagPRON // "some argue", "some believe", "each knows"
			}
		case "this", "these", "those":
			switch {
			case next.HasTag(TagNOUN|TagPROPN|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
				resolve = TagDET // "this decision", "these three", "those old"
			case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
				resolve = TagPRON // "this is clear", "these are done", "those were"
			case next.HasTag(TagPUNCT|TagCCONJ):
				resolve = TagPRON // "this.", "these,", "those and"
			case resolvedAs(prev, TagADP) && !next.HasTag(TagNOUN|TagADJ|TagPROPN|TagNUM):
				resolve = TagPRON // "of those", "in this" (not followed by noun)
			}
		case "another":
			switch {
			case next.HasTag(TagNOUN|TagPROPN|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
				resolve = TagDET // "another day", "another three"
			case next.HasTag(TagADP):
				resolve = TagPRON // "another of them"
			case next.HasTag(TagVERB|TagAUX|TagPUNCT):
				resolve = TagPRON // "another said", "another."
			}
		case "what":
			switch {
			case next.HasTag(TagNOUN|TagPROPN|TagADJ|TagNUM) && !next.HasTag(TagVERB|TagAUX):
				resolve = TagDET // "what year", "what time", "what kind"
			case next.HasTag(TagVERB|TagAUX|TagPUNCT|TagCCONJ):
				resolve = TagPRON // "what happened", "what?", "what he said"
			case resolvedAs(prev, TagVERB) || resolvedAs(prev, TagADP):
				resolve = TagPRON // "know what", "for what"
			}
		default:
			continue
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+detpron"
		}
	}
	return tokens
}
