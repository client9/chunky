package tok

import "strings"

// DisambiguateContractionFragments resolves Penn-treebank contraction fragments
// to AUX. Penn tokenizes "won't" as "wo"+"n't" and "can't" as "ca"+"n't"; our
// own splitter produces "will"+"n't" and "can"+"'t", so "wo" and "ca" only
// appear when processing externally-tokenized Penn input.
//
// "wo"/"ca" immediately before "n't" or "'t" → AUX.
func DisambiguateContractionFragments(tokens []Token) []Token {
	for i, t := range tokens {
		lower := strings.ToLower(t.Word)
		if lower != "wo" && lower != "ca" {
			continue
		}
		next := tokenAt(tokens, i+1)
		if next.Word == "n't" || next.Word == "'t" {
			tokens[i].Tags = TagAUX
			tokens[i].Rule = "contraction-fragment"
		}
	}
	return tokens
}
