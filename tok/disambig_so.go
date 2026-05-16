package tok

import "strings"

// DisambiguateSo resolves the ADV/SCONJ ambiguity on "so" and "once".
//
// "so" (ADV 84%, SCONJ 13%):
//
//	Intensifier (ADV):   "so good", "so quickly"
//	Subordinator (SCONJ): "so that everyone could see"
//
//	Resolved cases:
//	- next=ADJ|ADV → ADV   ("so good", "so carefully")
//	- next=SCONJ   → SCONJ ("so that")
//
// "once" (ADV 89%, SCONJ 10%):
//
//	Frequency/temporal (ADV): "happened once", "had once been", "once again"
//	Subordinator (SCONJ):     "once the treaty was signed"
//
//	Resolved cases:
//	- prev=AUX|ADV → ADV  ("had once been", "once again" when prev is adverb)
//	- next=ADV     → ADV  ("once more", "once again")
//
// Note: "so that" is handled upstream by MergeLexical as a compound SCONJ token;
// the next=SCONJ rule catches remaining cases like "so as", "so because".
func DisambiguateSo(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADV) || !t.HasTag(TagSCONJ) {
			continue
		}
		prev, next := tokenAt(tokens, i-1), tokenAt(tokens, i+1)
		var resolve Tag
		switch strings.ToLower(t.Word) {
		case "so":
			switch {
			case next.HasTag(TagADJ | TagADV):
				resolve = TagADV
			case resolvedAs(next, TagSCONJ):
				resolve = TagSCONJ
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "did so", "think so"
			case prev.HasTag(TagCCONJ):
				resolve = TagADV // "and so on", "or so", "but so"
			case next.HasTag(TagADP):
				resolve = TagADV // "so on", "so far", "so to speak"
			case next.HasTag(TagPRON | TagDET):
				resolve = TagSCONJ // "so he left", "so the war ended"
			case next.HasTag(TagAUX):
				resolve = TagSCONJ // "so was", "so had", "so does" (inverted/fronted)
			case resolvedAs(next, TagNOUN) || resolvedAs(next, TagPROPN):
				resolve = TagSCONJ // "so peace prevailed", "so victory was"
			case next.HasTag(TagPUNCT):
				resolve = TagADV // sentence-final "so."
			}
		case "once":
			switch {
			case prev.HasTag(TagAUX|TagADV) || next.HasTag(TagADV):
				resolve = TagADV
			case next.HasTag(TagSCONJ | TagADP):
				resolve = TagADV // "once upon a time", "once in a while", "once before"
			case strings.ToLower(next.Word) == "a" || strings.ToLower(next.Word) == "an":
				resolve = TagADV // "once a month", "once an hour" (frequentive)
			case next.HasTag(TagVERB | TagAUX):
				resolve = TagSCONJ // "once completed", "once approved"
			case next.HasTag(TagPRON):
				resolve = TagSCONJ // "once he arrived", "once they left"
			case next.Word == "the":
				resolve = TagSCONJ // "once the treaty was signed"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV // "it happened once.", "visited once more"
			case next.HasTag(TagPUNCT | TagCCONJ):
				resolve = TagADV // "once.", "once,", "once and for all"
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+so"
		}
	}
	return tokens
}
