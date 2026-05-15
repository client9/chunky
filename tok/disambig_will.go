package tok

import "github.com/client9/chunky"

// DisambiguateWill resolves the AUX/NOUN ambiguity on "will" and "Will".
//
// Modal use (AUX): "will go", "will not", "will never", "will he?"
// Noun use (NOUN): "his will", "the will of the people"
//
// Signal: next token carries VERB, AUX, ADV, PART, or PRON → AUX.
// Anything else (PUNCT, ADP, DET, NOUN, …) → leave ambiguous; the chunker
// will resolve via chunk position.
func DisambiguateWill(tokens []Token) []Token {
	for i, t := range tokens {
		if t.Word != "will" && t.Word != "Will" {
			continue
		}
		if !t.HasTag(chunky.TagAUX) || !t.HasTag(chunky.TagNOUN) {
			continue
		}
		if i+1 >= len(tokens) {
			continue
		}
		next := tokens[i+1]
		// DET immediately before "will" → always a noun phrase head.
		var resolve chunky.Tag
		if i > 0 && len(tokens[i-1].Tags) == 1 && tokens[i-1].Tags[0] == chunky.TagDET {
			resolve = chunky.TagNOUN
		} else {
			switch {
			case next.HasTag(chunky.TagVERB):
				// "will go", "will return" — modal before main verb
				resolve = chunky.TagAUX
			case next.HasTag(chunky.TagADV), next.HasTag(chunky.TagPART):
				// "will not", "will never", "will also" — adverb separates modal from verb
				resolve = chunky.TagAUX
			case next.HasTag(chunky.TagPRON):
				// "will he?", "will they?" — interrogative inversion
				resolve = chunky.TagAUX
			}
		}
		if resolve != 0 {
			tokens[i].Tags = []chunky.Tag{resolve}
			tokens[i].Rule = t.Rule + "+will"
		}
	}
	return tokens
}
