package tok

// disambiguateAdjNounDefault resolves {ADJ,NOUN} words that act as adjective
// modifiers when immediately preceding a nominal head.
//
// Registered for words where next+NOUN → ADJ is ≥97% in the corpus:
// social, public, native, red, white, unknown, general, equivalent, evil, etc.
func disambiguateAdjNounDefault(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ // "social media", "red car", "unknown quantity"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-noun"
	}
}

// disambiguateChief resolves "chief" ({ADJ,NOUN}).
// "chief executive" → ADJ; "editor in chief" → NOUN (next=ADP).
func disambiguateChief(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ // "chief executive", "chief engineer"
	case resolvedAs(next, TagADP):
		resolve = TagNOUN // "editor in chief", "commander in chief"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-noun"
	}
}

// disambiguateCapital resolves "capital" ({ADJ,NOUN}).
// "capital" is almost always NOUN — as a standalone head, in compound nouns
// ("capital city", "capital gains"), and after prepositions.
func disambiguateCapital(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagDET | TagADJ | TagNUM):
		resolve = TagNOUN // "the capital", "federal capital"
	case next.HasTag(TagADP | TagPUNCT | TagCCONJ | TagVERB | TagAUX):
		resolve = TagNOUN // "capital of", "capital.", "capital and"
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagNOUN // "capital city", "capital gains" (compound)
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-noun"
	}
}

// disambiguateFront resolves "front" ({ADJ,NOUN}).
// Before a nominal head it's ADJ ("front door"); after ADP or before ADP it's NOUN.
func disambiguateFront(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case next.HasTag(TagNOUN | TagADJ | TagPROPN):
		resolve = TagADJ // "front door", "front page", "front row"
	case next.HasTag(TagADP):
		resolve = TagNOUN // "front of the building"
	case resolvedAs(prev, TagADP):
		resolve = TagNOUN // "in front", "at the front"
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-noun"
	}
}

// disambiguateSideNoun resolves "side" ({ADJ,NOUN}).
// "side" is almost universally NOUN — as compound modifier ("side road"),
// as head ("the right side"), and after prepositions.
func disambiguateSideNoun(tokens []Token, i int) {
	t := tokens[i]
	if !t.HasTag(TagADJ) || !t.HasTag(TagNOUN) {
		return
	}
	prev := tokenAt(tokens, i-1)
	next := tokenAt(tokens, i+1)
	var resolve Tag
	switch {
	case prev.HasTag(TagADJ | TagDET | TagNUM):
		resolve = TagNOUN // "right side", "the side", "one side"
	case next.HasTag(TagADP | TagPUNCT | TagCCONJ | TagVERB | TagAUX):
		resolve = TagNOUN // "side of", "side.", "side and", "side was"
	case next.HasTag(TagNOUN | TagPROPN):
		resolve = TagNOUN // "side road", "side effect" (compound)
	}
	if resolve != 0 {
		tokens[i].Tags = resolve
		tokens[i].Rule = t.Rule + "+adj-noun"
	}
}
