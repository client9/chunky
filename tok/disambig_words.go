package tok

import "strings"

// wordHandlers maps lowercase word forms to their per-token disambiguation handler.
// DisambiguateWords uses a single map lookup instead of calling every handler for every token.
// Every word form that has special disambiguation rules must appear here explicitly.
var wordHandlers = map[string]func([]Token, int){
	// existential vs locative
	"there": disambiguateThere,
	// modal vs month name
	"may": disambiguateMay,
	// complementizer vs determiner vs pronoun
	"that": disambiguateThat,
	// discourse adverb vs attributive adjective
	"then": disambiguateThen,
	// modal vs noun
	"will": disambiguateWill,
	// ADJ/ADP group
	"like": disambiguateLike, "due": disambiguateLike,
	"pending": disambiguateLike, "pursuant": disambiguateLike,
	// particles: "off" is single-handler; "up" needs both particle and NOUN-bit handler
	"off": disambiguateParticles,
	"up":  disambiguateUpAll,
	// cardinal directions
	"south":     disambiguateDirectionals,
	"north":     disambiguateDirectionals,
	"east":      disambiguateDirectionals,
	"west":      disambiguateDirectionals,
	"northwest": disambiguateDirectionals,
	"northeast": disambiguateDirectionals,
	"southeast": disambiguateDirectionals,
	"southwest": disambiguateDirectionals,
	// intensifier vs subordinator
	"so":   disambiguateSo,
	"once": disambiguateOnce,
	// ordinal words
	"first":  disambiguateOrdinals,
	"second": disambiguateOrdinals,
	"third":  disambiguateOrdinals,
	// ADJ/DET/NOUN group
	"final":      disambiguateHalf,
	"half":       disambiguateHalf,
	"individual": disambiguateHalf,
	// ADJ/VERB
	"free":  disambiguateFree,
	"clear": disambiguateClear,
	// high group
	"high":    disambiguateHigh,
	"higher":  disambiguateHigh,
	"highest": disambiguateHigh,
	// right
	"right": disambiguateRight,
	// long group
	"long":    disambiguateLong,
	"longer":  disambiguateLong,
	"longest": disambiguateLong,
	// adj/adv default branch
	"alone":       disambiguateAdjAdvDefault,
	"bad":         disambiguateAdjAdvDefault,
	"clean":       disambiguateAdjAdvDefault,
	"cleaner":     disambiguateAdjAdvDefault,
	"cleanest":    disambiguateAdjAdvDefault,
	"closer":      disambiguateAdjAdvDefault,
	"closest":     disambiguateAdjAdvDefault,
	"cold":        disambiguateAdjAdvDefault,
	"colder":      disambiguateAdjAdvDefault,
	"coldest":     disambiguateAdjAdvDefault,
	"deep":        disambiguateAdjAdvDefault,
	"deeper":      disambiguateAdjAdvDefault,
	"deepest":     disambiguateAdjAdvDefault,
	"direct":      disambiguateAdjAdvDefault,
	"earliest":    disambiguateAdjAdvDefault,
	"far":         disambiguateAdjAdvDefault,
	"farthest":    disambiguateAdjAdvDefault,
	"fast":        disambiguateAdjAdvDefault,
	"faster":      disambiguateAdjAdvDefault,
	"fastest":     disambiguateAdjAdvDefault,
	"flat":        disambiguateAdjAdvDefault,
	"forward":     disambiguateAdjAdvDefault,
	"furthest":    disambiguateAdjAdvDefault,
	"good":        disambiguateAdjAdvDefault,
	"great":       disambiguateAdjAdvDefault,
	"greater":     disambiguateAdjAdvDefault,
	"greatest":    disambiguateAdjAdvDefault,
	"hard":        disambiguateAdjAdvDefault,
	"harder":      disambiguateAdjAdvDefault,
	"hardest":     disambiguateAdjAdvDefault,
	"ill":         disambiguateAdjAdvDefault,
	"latest":      disambiguateAdjAdvDefault,
	"loud":        disambiguateAdjAdvDefault,
	"louder":      disambiguateAdjAdvDefault,
	"loudest":     disambiguateAdjAdvDefault,
	"low":         disambiguateAdjAdvDefault,
	"lower":       disambiguateAdjAdvDefault,
	"lowest":      disambiguateAdjAdvDefault,
	"neat":        disambiguateAdjAdvDefault,
	"neater":      disambiguateAdjAdvDefault,
	"neatest":     disambiguateAdjAdvDefault,
	"overseas":    disambiguateAdjAdvDefault,
	"plain":       disambiguateAdjAdvDefault,
	"quick":       disambiguateAdjAdvDefault,
	"quicker":     disambiguateAdjAdvDefault,
	"quickest":    disambiguateAdjAdvDefault,
	"real":        disambiguateAdjAdvDefault,
	"rough":       disambiguateAdjAdvDefault,
	"rougher":     disambiguateAdjAdvDefault,
	"roughest":    disambiguateAdjAdvDefault,
	"short":       disambiguateAdjAdvDefault,
	"shorter":     disambiguateAdjAdvDefault,
	"shortest":    disambiguateAdjAdvDefault,
	"slow":        disambiguateAdjAdvDefault,
	"slower":      disambiguateAdjAdvDefault,
	"slowest":     disambiguateAdjAdvDefault,
	"small":       disambiguateAdjAdvDefault,
	"smaller":     disambiguateAdjAdvDefault,
	"smallest":    disambiguateAdjAdvDefault,
	"soft":        disambiguateAdjAdvDefault,
	"softer":      disambiguateAdjAdvDefault,
	"softest":     disambiguateAdjAdvDefault,
	"sure":        disambiguateAdjAdvDefault,
	"steady":      disambiguateAdjAdvDefault,
	"thick":       disambiguateAdjAdvDefault,
	"thicker":     disambiguateAdjAdvDefault,
	"thickest":    disambiguateAdjAdvDefault,
	"tight":       disambiguateAdjAdvDefault,
	"tighter":     disambiguateAdjAdvDefault,
	"tightest":    disambiguateAdjAdvDefault,
	"underground": disambiguateAdjAdvDefault,
	"wide":        disambiguateAdjAdvDefault,
	"wider":       disambiguateAdjAdvDefault,
	"widest":      disambiguateAdjAdvDefault,
	"wrong":       disambiguateAdjAdvDefault,
	// adj/adv special cases
	"prior":   disambiguatePrior,
	"likely":  disambiguateLikely,
	"later":   disambiguateLaterGroup,
	"earlier": disambiguateLaterGroup,
	"further": disambiguateLaterGroup,
	"early":   disambiguateEarlyLate,
	"late":    disambiguateEarlyLate,
	"dead":    disambiguateDead,
	"best":    disambiguateBest,
	"better":  disambiguateBetter,
	// adv/noun group
	"way":   disambiguateAdvNoun,
	"brand": disambiguateAdvNoun,
	"lot":   disambiguateAdvNoun,
	// adj/adv: only, little
	"only":   disambiguateOnly,
	"little": disambiguateOnly,
	// ADP/SCONJ
	"as":     disambiguateAs,
	"after":  disambiguateAfter,
	"before": disambiguateAfter,
	"until":  disambiguateAfter,
	// degree quantifiers
	"more":   disambiguateMore,
	"most":   disambiguateMore,
	"much":   disambiguateMore,
	"less":   disambiguateMore,
	"twice":  disambiguateMore,
	"enough": disambiguateMore,
	// common NOUN/VERB finite forms
	"says":   disambiguateVerbForms,
	"say":    disambiguateVerbForms,
	"remains": disambiguateVerbForms,
	"calls":  disambiguateVerbForms,
	"rose":   disambiguateVerbForms,
	"fell":   disambiguateVerbForms,
	"runs":   disambiguateVerbForms,
	"turns":  disambiguateVerbForms,
	"holds":  disambiguateVerbForms,
	"needs":  disambiguateVerbForms,
	"wants":  disambiguateVerbForms,
	"plans":  disambiguateVerbForms,
	"shows":  disambiguateVerbForms,
	"leads":  disambiguateVerbForms,
	"leaves": disambiguateVerbForms,
	"means":  disambiguateVerbForms,
	"takes":  disambiguateVerbForms,
	"makes":  disambiguateVerbForms,
	"comes":  disambiguateVerbForms,
	"goes":   disambiguateVerbForms,
	"gives":  disambiguateVerbForms,
	"brings": disambiguateVerbForms,
	"adds":   disambiguateVerbForms,
	"argues": disambiguateVerbForms,
	// adv/noun
	"back": disambiguateBack,
	"well": disambiguateWell,
	// ADP/ADV directional
	"down": disambiguateDown,
	"near": disambiguateDown,
	// quantifiers: both/neither/either/all single-handler; each/any dual-handler
	"both":    disambiguateQuantifiers,
	"neither": disambiguateQuantifiers,
	"either":  disambiguateQuantifiers,
	"all":     disambiguateQuantifiers,
	"each":    disambiguateEachAll,
	"any":     disambiguateAnyAll,
	// such
	"such": disambiguateSuch,
	// spatial adverbs
	"outside": disambiguateAbove,
	"above":   disambiguateAbove,
	"inside":  disambiguateAbove,
	// yet
	"yet": disambiguateYet,
	// past
	"past": disambiguatePast,
	// pro
	"pro": disambiguatePro,
	// following
	"following": disambiguateFollowing,
	// one
	"one": disambiguateOne,
	// spatial prepositions
	"about":  disambiguatePrepositions,
	"around": disambiguatePrepositions,
	"below":  disambiguatePrepositions,
	"behind": disambiguatePrepositions,
	"out":    disambiguatePrepositions,
	// mine, u
	"mine": disambiguateMine,
	"u":    disambiguateU,
	// ADJ/VERB words requiring possessive-pronoun context
	"own":           disambiguateOwn,
	"live":          disambiguateOwn,
	"separate":      disambiguateOwn,
	"complete":      disambiguateOwn,
	"correct":       disambiguateOwn,
	"dry":           disambiguateOwn,
	"warm":          disambiguateOwn,
	"smooth":        disambiguateOwn,
	"secure":        disambiguateOwn,
	"frequent":      disambiguateOwn,
	"lasting":       disambiguateOwn,
	"varying":       disambiguateOwn,
	"corresponding": disambiguateOwn,
	"marked":        disambiguateOwn,
	"elaborate":     disambiguateOwn,
	"engaging":      disambiguateOwn,
	// DET/PRON: some/this/these/those/another/what (each/any handled above)
	"some":    disambiguateDetPron,
	"this":    disambiguateDetPron,
	"these":   disambiguateDetPron,
	"those":   disambiguateDetPron,
	"another": disambiguateDetPron,
	"what":    disambiguateDetPron,
	// plus
	"plus": disambiguatePlus,
	// ADP/VERB
	"save":        disambiguateSave,
	"respecting":  disambiguateSave,
	// to
	"to": disambiguateTo,
	// SCONJ → ADP
	"since":   disambiguateSCONJasADP,
	"despite": disambiguateSCONJasADP,
	"upon":    disambiguateSCONJasADP,
}

// disambiguateUpAll chains the particle handler (ADP|ADV) with the NOUN-bit handler.
func disambiguateUpAll(tokens []Token, i int) {
	disambiguateParticles(tokens, i)
	disambiguateUp(tokens, i)
}

// disambiguateEachAll chains quantifier and det/pron handlers for "each".
func disambiguateEachAll(tokens []Token, i int) {
	disambiguateQuantifiers(tokens, i)
	disambiguateDetPron(tokens, i)
}

// disambiguateAnyAll chains quantifier and det/pron handlers for "any".
func disambiguateAnyAll(tokens []Token, i int) {
	disambiguateQuantifiers(tokens, i)
	disambiguateDetPron(tokens, i)
}

// DisambiguateWords applies all word-specific disambiguation in a single pass.
// It must run after LexicalTag and TagUnknowns, and before sentence segmentation
// and context disambiguation.
//
// ContractionFragments and ApostropheS run as separate pre-passes: the former
// has no tag guard, the latter mutates neighboring tokens and must complete
// before the main pass reads neighbors.
//
// To add a new word-specific disambiguator: implement a func([]Token, int) helper
// and register it in wordHandlers above.
func DisambiguateWords(tokens []Token) []Token {
	tokens = DisambiguateContractionFragments(tokens)
	tokens = DisambiguateApostropheS(tokens)
	for i := range tokens {
		lw := strings.ToLower(tokens[i].Word)
		if h, ok := wordHandlers[lw]; ok {
			h(tokens, i)
		}
		disambiguateAdjVerb(tokens, i)
	}
	return tokens
}
