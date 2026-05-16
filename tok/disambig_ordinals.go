package tok

import "strings"

// DisambiguateOrdinals resolves the ADV/NUM ambiguity on ordinal words.
//
// "first" lexicon: ADV|NUM
// "second", "third" lexicon: ADV|NOUN|NUM
//
// Ordinal (NUM): prev=DET|PRON and next=NOUN|ADJ|PROPN — prenominal use
//
//	"the first chapter", "my second attempt", "a third party"
//
// NOUN (second/third only): standalone or time-unit use
//
//	"a second", "in a second", "a close third"
//
// Sentential adverb (ADV): next=VERB — discourse/sequential use
//
//	"first, consider the options", "we must first decide"
func DisambiguateOrdinals(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADV) || !t.HasTag(TagNUM) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "first", "second", "third":
		default:
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch {
		case prev.HasTag(TagDET|TagPRON) && next.HasTag(TagNOUN|TagADJ|TagPROPN):
			resolve = TagNUM
		case prev.HasTag(TagNOUN | TagPROPN):
			resolve = TagNUM // "June first/second/third", "a split second"
		case resolvedAs(prev, TagADJ):
			resolve = TagNUM // "a close second", "a strong third"
		case next.HasTag(TagVERB | TagAUX):
			resolve = TagADV
		case next.HasTag(TagCCONJ):
			resolve = TagADV // "first and foremost", "first or last"
		case next.HasTag(TagADP):
			resolve = TagADV // "first of all", "second of the month"
		case next.HasTag(TagPUNCT) && !prev.HasTag(TagDET|TagPRON|TagNOUN|TagPROPN):
			resolve = TagADV // "First,", "second." — discourse marker at sentence start or after clause
		case t.HasTag(TagNOUN) && prev.HasTag(TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
			resolve = TagNOUN // "a second", "the third" (time unit, no following noun)
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+ord"
		}
	}
	return tokens
}
