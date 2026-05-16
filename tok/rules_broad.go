package tok

// detPronBroadRules resolves DET/PRON ambiguity (this/these/some/any) using
// only the immediately following tag. These are 1-slot fallbacks that fire
// after the more-specific generated rules.
var detPronBroadRules = []ContextRule{
	// Before a nominal or pre-nominal modifier → determiner
	{Tags: TagDET | TagPRON, Next: TagNOUN | TagPROPN | TagADJ | TagNUM, Mask: maskNext, Resolve: TagDET},
	// Before a verbal head or boundary → pronoun
	{Tags: TagDET | TagPRON, Next: TagAUX | TagVERB | TagPUNCT | TagADP | TagCCONJ, Mask: maskNext, Resolve: TagPRON},
}

// advDetBroadRules resolves ADV/DET ambiguity (most/more/less/much) using
// only the immediately following tag.
var advDetBroadRules = []ContextRule{
	// Before an adjective or adverb → intensifier (ADV): "most important", "more quickly"
	{Tags: TagADV | TagDET, Next: TagADJ | TagADV, Mask: maskNext, Resolve: TagADV},
	// Before a noun head → quantifier (DET): "most people", "more money"
	{Tags: TagADV | TagDET, Next: TagNOUN | TagPROPN, Mask: maskNext, Resolve: TagDET},
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
	{Tags: TagNOUN | TagVERB, Next: TagPRON | TagADJ | TagADV, Mask: maskNext, Resolve: TagVERB},
}

// adjNounBroadRules resolves ADJ/NOUN ambiguity using the immediately following
// tag. When an ADJ|NOUN token precedes a resolved NOUN or PROPN it is acting
// as a prenominal modifier → ADJ ("national team", "military force").
// Corpus precision for next=NOUN/PROPN is >99%.
var adjNounBroadRules = []ContextRule{
	{Tags: TagADJ | TagNOUN, Next: TagNOUN | TagPROPN, Mask: maskNext, Resolve: TagADJ},
}

// adpSconjBroadRules resolves ADP/SCONJ ambiguity (after/before/until) using
// the immediately following tag. Only resolves to ADP — SCONJ has no clean
// single-token signal and is left for the generated 4-slot rules.
//
// ADP signals: next is non-clausal (no subject NP follows the preposition).
// next=NUM/AUX/VERB are 98–100% ADP in corpus; next=ADJ/ADV are also clean.
// next=DET/PRON/NOUN/PROPN are too mixed (SCONJ also takes NP subjects).
var adpSconjBroadRules = []ContextRule{
	{Tags: TagADP | TagSCONJ, Next: TagVERB | TagAUX | TagNUM | TagADJ | TagADV, Mask: maskNext, Resolve: TagADP},
}
