package tok

import "github.com/client9/chunky"

// closedClassOnly is the set of tags that are purely functional — words with
// only these tags can never be proper nouns regardless of capitalization.
const closedClassOnly = chunky.TagPRON | chunky.TagDET | chunky.TagADP | chunky.TagAUX |
	chunky.TagSCONJ | chunky.TagCCONJ | chunky.TagPART | chunky.TagPUNCT |
	chunky.TagINTJ | chunky.TagSYM

// RetagCapitalized promotes mid-sentence capitalized known words to PROPN.
// It runs per-sentence (inside Segment) so i==0 means sentence-initial.
func RetagCapitalized(tokens []Token) []Token {
	for i, t := range tokens {
		if len(t.Word) == 0 || t.Word[0] < 'A' || t.Word[0] > 'Z' {
			continue
		}
		if i == 0 {
			// Sentence-initial capitalization is grammatical, not semantic.
			// We cannot distinguish a proper noun from a sentence-initial verb or
			// participle (e.g. "Walked the dog." vs "Ted is here."), so we leave
			// the pipeline tag unchanged.
			continue
		}
		// Skip words whose entire candidate set is closed-class: function words
		// (ADP, SCONJ, CCONJ, AUX, DET, PART, PRON, PUNCT …) are never proper
		// nouns regardless of capitalization. E.g. `` In the morning" should
		// keep "In" as ADP, not become PROPN.
		if t.Tags&^closedClassOnly == 0 {
			continue
		}
		// Non-sentence-initial: known open-class capitalized word → PROPN.
		if !t.IsUnknownTag() {
			tokens[i].Tags = chunky.TagPROPN
			tokens[i].Rule = t.Rule + "+caps"
		}
	}
	return tokens
}
