package tok

import "strings"

// DisambiguateAdjAdv resolves the ADJ/ADV ambiguity on common dual-category words.
// Also handles {ADJ,ADV,NOUN} and {ADJ,ADV,VERB} tokens (guard only requires both ADJ+ADV bits).
//
// Common rule: next=NOUN|ADJ|PROPN → ADJ (prenominal modifier)
// Prev=VERB rule: prev=VERB → ADV (post-verbal modifier)
func DisambiguateAdjAdv(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
			continue
		}
		lw := strings.ToLower(t.Word)
		switch lw {
		case "alone", "bad", "best", "better", "clean", "cleaner", "cleanest",
			"closer", "closest", "cold", "colder", "coldest",
			"dead", "deep", "deeper", "deepest",
			"direct", "early", "earlier", "earliest",
			"far", "farthest", "fast", "faster", "fastest", "flat",
			"further", "furthest", "forward",
			"good", "great", "greater", "greatest",
			"hard", "harder", "hardest", "high", "higher", "highest",
			"ill", "late", "later", "latest", "likely", "long", "longer",
			"longest", "loud", "louder", "loudest",
			"low", "lower", "lowest", "neat", "neater", "neatest",
			"overseas", "plain", "prior", "quick", "quicker", "quickest",
			"real", "right", "rough", "rougher", "roughest",
			"short", "shorter", "shortest", "slow", "slower", "slowest",
			"small", "smaller", "smallest", "soft", "softer", "softest",
			"sure", "steady", "thick", "thicker", "thickest",
			"tight", "tighter", "tightest", "underground", "wide", "wider", "widest",
			"wrong":
		default:
			continue
		}
		prev := tokenAt(tokens, i-1)
		next := tokenAt(tokens, i+1)
		var resolve Tag
		switch lw {
		case "prior":
			switch {
			case next.HasTag(TagADP):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			}
		case "likely":
			switch {
			case next.HasTag(TagPART):
				resolve = TagADJ
			case next.HasTag(TagVERB | TagAUX):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			}
		case "later", "earlier", "longer", "further":
			switch {
			case strings.ToLower(prev.Word) == "no":
				resolve = TagADV // "no longer", "no further"
			case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
				resolve = TagADJ // "later chapter", "earlier version"
			case prev.HasTag(TagPRON | TagNOUN | TagNUM):
				resolve = TagADV // "see you later", "three years later", "days earlier"
			case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
				resolve = TagADV
			case next.HasTag(TagPUNCT | TagCCONJ):
				resolve = TagADV // "later.", "Later, he..."
			}
		case "early", "late":
			switch {
			case next.HasTag(TagADP):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagADJ | TagPROPN | TagNUM):
				resolve = TagADJ
			}
		case "right":
			// "right decision" → ADJ; "right here/now/away", "turned right" → ADV; "the right (to vote)" → NOUN
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			case next.HasTag(TagADV | TagDET | TagPRON):
				resolve = TagADV
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV
			case t.HasTag(TagNOUN) && resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
				resolve = TagNOUN
			}
		case "dead":
			// "dead wrong/tired/certain" → ADV; "dead end/heat" → ADJ
			switch {
			case next.HasTag(TagADJ):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagPROPN):
				resolve = TagADJ
			}
		case "long", "best":
			// "long road", "best option" → ADJ; "waited long", "works best" → ADV
			switch {
			case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
				resolve = TagADJ // "best option", "long road" (only unambiguous nouns)
			case prev.HasTag(TagVERB | TagAUX):
				resolve = TagADV
			case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && next.HasTag(TagPUNCT):
				resolve = TagADV // "she knew best.", "he waited long."
			case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
				resolve = TagADV // "best served cold", "how long has"
			case next.HasTag(TagADP | TagSCONJ):
				resolve = TagADV // "long before", "best of all"
			}
		case "better":
			// "better option" → ADJ; "works better", "feel better" → ADV
			switch {
			case prev.HasTag(TagVERB | TagAUX):
				resolve = TagADV
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			}
		case "free", "clear":
			// "free access", "clear signal" → ADJ; "ran free", "stand clear" → ADV
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				resolve = TagADJ
			case next.HasTag(TagPUNCT | TagCCONJ | TagADP | TagPART):
				resolve = TagADJ // "free and open", "free of charge", "free to choose"
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				resolve = TagADV
			case t.HasTag(TagVERB) && prev.HasTag(TagAUX):
				resolve = TagVERB // "will free", "must clear"
			}
		default:
			// alone, bad, clean, closer, cold, direct, far, fast, faster, flat,
			// good, hard, high, higher, ill, low, lower, overseas, plain, short,
			// slow, sure, steady, thick, underground, wrong
			switch {
			case next.HasTag(TagNOUN | TagADJ | TagPROPN):
				// prenominal takes priority: "running high temperatures" → ADJ
				resolve = TagADJ
			case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
				// post-verbal modifier: "ran high", "went wrong", "fell flat"
				// Guard: copulas (is/are/was) carry AUX bit so won't fire here.
				resolve = TagADV
			case prev.HasTag(TagADV | TagPART):
				resolve = TagADV // "very far", "not hard", "too short", "almost alone"
			case prev.HasTag(TagAUX) && next.HasTag(TagPUNCT|TagCCONJ|TagSCONJ|TagADP):
				resolve = TagADJ // predicative after copula: "was high.", "is good,"
			case next.HasTag(TagADV) && !next.HasTag(TagNOUN|TagVERB|TagADJ):
				resolve = TagADV // "far enough", "hard enough", "high above"
			case next.HasTag(TagADP) && !next.HasTag(TagNOUN|TagVERB):
				resolve = TagADV // "far from", "short of", "hard on"
			case next.HasTag(TagPUNCT|TagCCONJ) && prev.HasTag(TagADP|TagADV):
				resolve = TagADV // "came in high,", "running high and"
			case t.HasTag(TagNOUN) && resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
				// "the high", "a flat", "the good" — nominal use after determiner
				resolve = TagNOUN
			case t.HasTag(TagNOUN) && prev.HasTag(TagNOUN|TagADJ) && !prev.HasTag(TagVERB) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
				resolve = TagNOUN // "record high", "all-time low"
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+adj-adv"
		}
	}
	return tokens
}
