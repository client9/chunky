package tok

// DisambiguateWords applies all word-specific disambiguation passes in order.
// It must run after LexicalTag and TagUnknowns, and before sentence segmentation
// and context disambiguation.
//
// Add new word-specific passes here — this is the single registration point.
func DisambiguateWords(tokens []Token) []Token {
	tokens = DisambiguateApostropheS(tokens)
	tokens = DisambiguateThere(tokens)
	tokens = DisambiguateMay(tokens)
	tokens = DisambiguateThat(tokens)
	tokens = DisambiguateThen(tokens)
	tokens = DisambiguateWill(tokens)
	return tokens
}
