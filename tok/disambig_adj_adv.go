package tok

import "strings"

// All per-word helpers below are registered in wordHandlers (disambig_words.go).
// Each function has a tag guard but no word guard — dispatch guarantees the word.

func disambiguateAdjAdvDefault(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV
	case prev.HasTag(TagADV | TagPART):
		resolve = TagADV // "very far", "not hard", "too short", "almost alone"
	case prev.HasTag(TagAUX) && next.HasTag(TagPUNCT|TagCCONJ|TagSCONJ|TagADP):
		resolve = TagADJ // predicative after copula: "was good.", "is sure,"
	case next.HasTag(TagADV) && !next.HasTag(TagNOUN|TagVERB|TagADJ):
		resolve = TagADV // "far enough", "hard enough"
	case next.HasTag(TagADP) && !next.HasTag(TagNOUN|TagVERB):
		resolve = TagADV // "far from", "short of", "hard on"
	case next.HasTag(TagPUNCT|TagCCONJ) && prev.HasTag(TagADP|TagADV):
		resolve = TagADV // "came in low,", "running low and"
	case t.HasTag(TagNOUN) && resolvedAs(prev, TagDET) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagNOUN // "a flat", "the good"
	case t.HasTag(TagNOUN) && prev.HasTag(TagNOUN|TagADJ) && !prev.HasTag(TagVERB) && !next.HasTag(TagNOUN|TagADJ|TagPROPN):
		resolve = TagNOUN // "all-time low"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguatePrior(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagADP):
		resolve = TagADV // "prior to"
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ // "prior experience"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguateLikely(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagPART):
		resolve = TagADJ // "likely to win"
	case next.HasTag(TagVERB | TagAUX):
		resolve = TagADV
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

// disambiguateLaterGroup handles later/earlier/further — comparatives whose
// ADV vs ADJ distinction depends on whether they precede a noun or follow one.
func disambiguateLaterGroup(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case strings.ToLower(prev.Word) == "no":
		resolve = TagADV // "no further"
	case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
		resolve = TagADJ // "later chapter", "earlier version", "further evidence"
	case prev.HasTag(TagPRON | TagNOUN | TagNUM):
		resolve = TagADV // "see you later", "three years later", "days earlier"
	case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
		resolve = TagADV
	case next.HasTag(TagPUNCT | TagCCONJ):
		resolve = TagADV // "later.", "Later, he..."
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguateEarlyLate(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagADP):
		resolve = TagADV // "early in the morning", "late in the term"
	case next.HasTag(TagNOUN | TagADJ | TagPROPN | TagNUM):
		resolve = TagADJ // "early results", "late stage"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguateDead(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagADJ):
		resolve = TagADV // "dead wrong", "dead tired"
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagADJ // "dead end", "dead heat"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguateBest(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN|TagADJ|TagPROPN) && !next.HasTag(TagVERB|TagAUX):
		resolve = TagADJ
	case prev.HasTag(TagVERB | TagAUX):
		resolve = TagADV
	case prev.HasTag(TagPRON|TagNOUN|TagPROPN) && next.HasTag(TagPUNCT):
		resolve = TagADV // "she knew best."
	case resolvedAs(next, TagVERB) || resolvedAs(next, TagAUX):
		resolve = TagADV // "best served cold"
	case next.HasTag(TagADP | TagSCONJ):
		resolve = TagADV // "best of all"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

func disambiguateBetter(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagVERB | TagAUX):
		resolve = TagADV
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

// disambiguateClear resolves "clear" ({ADJ,ADV}).
// "clear signal" → ADJ; "stand clear" → ADV; "will clear" → VERB
func disambiguateClear(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagADV) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ
	case next.HasTag(TagPUNCT | TagCCONJ | TagADP | TagPART):
		resolve = TagADJ // "clear and present", "clear of charge", "clear to proceed"
	case prev.HasTag(TagVERB) && !prev.HasTag(TagAUX):
		resolve = TagADV
	case t.HasTag(TagVERB) && prev.HasTag(TagAUX):
		resolve = TagVERB // "will clear", "must clear"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-adv"
	}
}

// DisambiguateAdjAdv applies all ADJ/ADV helpers (kept for external use).
func DisambiguateAdjAdv(tokens []Token) []Token {
	for i := range tokens {
		lw := strings.ToLower(tokens[i].Word)
		switch lw {
		case "prior":
			disambiguatePrior(tokens, i)
		case "likely":
			disambiguateLikely(tokens, i)
		case "later", "earlier", "further":
			disambiguateLaterGroup(tokens, i)
		case "early", "late":
			disambiguateEarlyLate(tokens, i)
		case "dead":
			disambiguateDead(tokens, i)
		case "best":
			disambiguateBest(tokens, i)
		case "better":
			disambiguateBetter(tokens, i)
		case "clear":
			disambiguateClear(tokens, i)
		default:
			disambiguateAdjAdvDefault(tokens, i)
		}
	}
	return tokens
}
