package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// neverPropn is the set of lowercase words that should never be promoted to
// PROPN by RetagCapitalized even when capitalized mid-sentence. These are
// common nouns and abbreviations that appear capitalized for reasons other
// than being proper names (legal terms like "Chapter 11", acronyms like "TV").
var neverPropn = map[string]bool{
	"chapter": true, "section": true, "article": true,
	"tv": true, "cds": true, "cd": true,
	"adrs": true, "adr": true, "dvd": true, "dvds": true,
	"remics": true, "remic": true,
}

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
			// Exception: a NOUN, ADJ, or {ADJ,NOUN} token immediately before a
			// capitalized open-class word is a name component — promote to PROPN.
			// "Robert Edward Turner": Robert→NOUN, Edward→{PROPN,ADJ}
			// "Great American Bank": Great→ADJ, American→NOUN→promoted
			// "Northeast Brazil": Northeast→{ADJ,NOUN}, Brazil→NOUN→promoted
			isNominalCandidate := t.Tags == chunky.TagNOUN ||
				t.Tags == chunky.TagADJ ||
				t.Tags == chunky.TagADJ|chunky.TagNOUN
			if i+1 < len(tokens) && isNominalCandidate {
				next := tokens[i+1]
				if len(next.Word) > 0 && next.Word[0] >= 'A' && next.Word[0] <= 'Z' &&
					next.Tags&^closedClassOnly != 0 {
					tokens[i].Tags = chunky.TagPROPN
					tokens[i].Rule = t.Rule + "+caps"
				}
			}
			continue
		}
		// Skip words whose entire candidate set is closed-class: function words
		// (ADP, SCONJ, CCONJ, AUX, DET, PART, PRON, PUNCT …) are never proper
		// nouns regardless of capitalization. E.g. `` In the morning" should
		// keep "In" as ADP, not become PROPN.
		if t.Tags&^closedClassOnly == 0 {
			continue
		}
		// Skip words that are never proper names regardless of capitalization.
		if neverPropn[strings.ToLower(t.Word)] {
			continue
		}
		// Non-sentence-initial open-class capitalized word → PROPN.
		if !t.IsUnknownTag() {
			tokens[i].Tags = chunky.TagPROPN
			tokens[i].Rule = t.Rule + "+caps"
		}
	}
	return tokens
}
