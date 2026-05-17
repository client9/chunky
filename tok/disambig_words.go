package tok

// DisambiguateWords applies all word-specific disambiguation passes in order.
// It must run after LexicalTag and TagUnknowns, and before sentence segmentation
// and context disambiguation.
//
// Add new word-specific passes here — this is the single registration point.
func DisambiguateWords(tokens []Token) []Token {
	tokens = DisambiguateContractionFragments(tokens)
	tokens = DisambiguateApostropheS(tokens)
	tokens = DisambiguateThere(tokens)
	tokens = DisambiguateMay(tokens)
	tokens = DisambiguateThat(tokens)
	tokens = DisambiguateThen(tokens)
	tokens = DisambiguateWill(tokens)
	tokens = DisambiguateLike(tokens)
	tokens = DisambiguateParticles(tokens)
	tokens = DisambiguateSpatial(tokens)
	tokens = DisambiguateDirectionals(tokens)
	tokens = DisambiguateSo(tokens)
	tokens = DisambiguateOrdinals(tokens)
	tokens = DisambiguateHalf(tokens)
	tokens = DisambiguateFree(tokens)
	tokens = DisambiguateHigh(tokens)
	tokens = DisambiguateRight(tokens)
	tokens = DisambiguateAdjAdv(tokens)
	tokens = DisambiguateAdvNoun(tokens)
	tokens = DisambiguateAs(tokens)
	tokens = DisambiguateMore(tokens)
	tokens = DisambiguateAfter(tokens)
	tokens = DisambiguateVerbForms(tokens)
	tokens = DisambiguateWellBack(tokens)
	tokens = DisambiguateDown(tokens)
	tokens = DisambiguateUp(tokens)
	tokens = DisambiguateLong(tokens)
	tokens = DisambiguateQuantifiers(tokens)
	tokens = DisambiguateSuch(tokens)
	tokens = DisambiguateAbove(tokens)
	tokens = DisambiguateYet(tokens)
	tokens = DisambiguatePast(tokens)
	tokens = DisambiguateFollowing(tokens)
	tokens = DisambiguateOne(tokens)
	tokens = DisambiguatePrepositions(tokens)
	tokens = DisambiguateMine(tokens)
	tokens = DisambiguateOwn(tokens)
	tokens = DisambiguateDetPron(tokens)
	tokens = DisambiguatePlus(tokens)
	tokens = DisambiguateSave(tokens)
	tokens = DisambiguateTo(tokens)
	tokens = DisambiguateAdjVerb(tokens)
	tokens = DisambiguateSCONJasADP(tokens)
	return tokens
}
