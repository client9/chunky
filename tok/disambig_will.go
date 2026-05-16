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
		var resolve chunky.Tag
		if i > 0 && tokens[i-1].IsResolved() && tokens[i-1].Tags == chunky.TagDET {
			resolve = chunky.TagNOUN
		} else {
			switch {
			case next.HasTag(chunky.TagVERB), next.HasTag(chunky.TagAUX):
				// "will go", "will be" — modal before verbal
				resolve = chunky.TagAUX
			case next.HasTag(chunky.TagADV), next.HasTag(chunky.TagPART):
				// "will not", "will never"
				resolve = chunky.TagAUX
			case next.HasTag(chunky.TagPRON):
				// "will he?" — interrogative inversion
				resolve = chunky.TagAUX
			}
		}
		if resolve != 0 {
			tokens[i].Tags = resolve
			tokens[i].Rule = t.Rule + "+will"
		}
	}
	return tokens
}
