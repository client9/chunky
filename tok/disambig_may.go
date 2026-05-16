package tok

import "github.com/client9/chunky"

// DisambiguateMay resolves the AUX/PROPN ambiguity on "may" and "May" tokens.
//
// Two safe cases are handled pre-sentence:
//   - Lowercase "may" is always the modal auxiliary.
//   - Capitalized "May" before a NUM is the month name: "May 15", "May 2024".
//
// Sentence-initial "May" before a PRON is handled in RetagMay (per-sentence),
// because distinguishing "May I?" (AUX) from "in May we..." (PROPN) requires
// knowing sentence position.
func DisambiguateMay(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.HasTag(chunky.TagAUX) || !t.HasTag(chunky.TagPROPN) {
			continue
		}
		switch t.Word {
		case "may":
			tokens[i].Tags = chunky.TagAUX
			tokens[i].Rule = t.Rule + "+may"
		case "May":
			if i+1 < len(tokens) && tokens[i+1].HasTag(chunky.TagNUM) {
				tokens[i].Tags = chunky.TagPROPN
				tokens[i].Rule = t.Rule + "+may"
			}
		}
	}
	return tokens
}

// RetagMay resolves sentence-initial "May" that DisambiguateMay could not
// resolve pre-sentence. Called once per sentence, before DisambiguateContext.
//
//   - "May" before PRON → AUX  ("May I help you?", "May we proceed?")
//   - "May" before anything else → PROPN  ("May flowers bloom early.")
func RetagMay(tokens []Token) []Token {
	if len(tokens) == 0 {
		return tokens
	}
	t := tokens[0]
	if t.Word != "May" || !t.HasTag(chunky.TagAUX) || !t.HasTag(chunky.TagPROPN) {
		return tokens
	}
	tag := chunky.Tag(chunky.TagPROPN)
	if len(tokens) > 1 && tokens[1].HasTag(chunky.TagPRON) {
		tag = chunky.TagAUX
	}
	tokens[0].Tags = tag
	tokens[0].Rule = t.Rule + "+may"
	return tokens
}
