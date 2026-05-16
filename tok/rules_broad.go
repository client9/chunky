package tok

import "github.com/client9/chunky"

// detPronBroadRules resolves DET/PRON ambiguity (this/these/some/any) using
// only the immediately following tag. These are 1-slot fallbacks that fire
// after the more-specific generated rules.
var detPronBroadRules = []ContextRule{
	// Before a nominal or pre-nominal modifier → determiner
	{Tag1: chunky.TagDET, Tag2: chunky.TagPRON, Next: chunky.TagNOUN | chunky.TagPROPN | chunky.TagADJ | chunky.TagNUM, Mask: maskNext, Resolve: chunky.TagDET},
	// Before a verbal head or boundary → pronoun
	{Tag1: chunky.TagDET, Tag2: chunky.TagPRON, Next: chunky.TagAUX | chunky.TagVERB | chunky.TagPUNCT | chunky.TagADP | chunky.TagCCONJ, Mask: maskNext, Resolve: chunky.TagPRON},
}

// advDetBroadRules resolves ADV/DET ambiguity (most/more/less/much) using
// only the immediately following tag.
var advDetBroadRules = []ContextRule{
	// Before an adjective or adverb → intensifier (ADV): "most important", "more quickly"
	{Tag1: chunky.TagADV, Tag2: chunky.TagDET, Next: chunky.TagADJ | chunky.TagADV, Mask: maskNext, Resolve: chunky.TagADV},
	// Before a noun head → quantifier (DET): "most people", "more money"
	{Tag1: chunky.TagADV, Tag2: chunky.TagDET, Next: chunky.TagNOUN | chunky.TagPROPN, Mask: maskNext, Resolve: chunky.TagDET},
}

// nounVerbBroadRules resolves NOUN/VERB ambiguity using only the immediately
// following tag. The generated 4-slot rules handle next=DET already (1-slot
// rule exists); these add the remaining high-confidence next-tag signals.
// Corpus precision on target verbs (says/feel/make/cut/…):
//
//	next=PRON  97% VERB   "company says it",  "firm says he"
//	next=ADJ   97% VERB   "remains unclear",  "feels right"
//	next=ADV   89% VERB   "rose sharply",     "fell further"
var nounVerbBroadRules = []ContextRule{
	{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Next: chunky.TagPRON | chunky.TagADJ | chunky.TagADV, Mask: maskNext, Resolve: chunky.TagVERB},
}

// adpSconjBroadRules resolves ADP/SCONJ ambiguity (after/before/until) using
// the immediately following tag. Only resolves to ADP — SCONJ has no clean
// single-token signal and is left for the generated 4-slot rules.
//
// ADP signals: next is non-clausal (no subject NP follows the preposition).
// next=NUM/AUX/VERB are 98–100% ADP in corpus; next=ADJ/ADV are also clean.
// next=DET/PRON/NOUN/PROPN are too mixed (SCONJ also takes NP subjects).
var adpSconjBroadRules = []ContextRule{
	{Tag1: chunky.TagADP, Tag2: chunky.TagSCONJ, Next: chunky.TagVERB | chunky.TagAUX | chunky.TagNUM | chunky.TagADJ | chunky.TagADV, Mask: maskNext, Resolve: chunky.TagADP},
}
