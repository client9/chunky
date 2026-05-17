package chunky

import (
	"math/bits"
	"testing"
)

func TestTagFromPennTag(t *testing.T) {
	cases := []struct {
		penn     string
		wantBits int // expected number of UD bits (1 = unambiguous, >1 = ambiguous)
		wantTag  Tag // exact match when wantBits == 1
		mustHave Tag // for ambiguous: all these bits must be set
	}{
		// Unambiguous — single UD mapping.
		{"NN", 1, TagNOUN, 0},
		{"NNS", 1, TagNOUN, 0},
		{"NNP", 1, TagPROPN, 0},
		{"NNPS", 1, TagPROPN, 0},
		{"MD", 1, TagAUX, 0},
		{"CC", 1, TagCCONJ, 0},
		{"RP", 1, TagPART, 0},
		{"POS", 1, TagPART, 0},
		{"CD", 1, TagNUM, 0},
		{"UH", 1, TagINTJ, 0},
		{"FW", 1, TagX, 0},
		{".", 1, TagPUNCT, 0},
		{",", 1, TagPUNCT, 0},
		{"$", 1, TagSYM, 0},

		// Pronouns — unambiguous in Penn.
		{"PRP", 1, TagPRON, 0},
		{"PRP$", 1, TagPRON, 0},
		{"WP", 1, TagPRON, 0},
		{"EX", 1, TagPRON, 0},

		// DT/PDT — Penn uses DT for both prenominal determiners and standalone
		// pronouns ("this is good" → PRON in UD). Map to DET|PRON.
		{"DT", 2, 0, TagDET | TagPRON},
		{"PDT", 2, 0, TagDET | TagPRON},

		// Ambiguous verb tags — Penn has no AUX; "was/VBD" is AUX in UD.
		{"VB", 2, 0, TagVERB | TagAUX},
		{"VBD", 2, 0, TagVERB | TagAUX},
		{"VBG", 2, 0, TagVERB | TagAUX},
		{"VBN", 2, 0, TagVERB | TagAUX},
		{"VBP", 2, 0, TagVERB | TagAUX},
		{"VBZ", 2, 0, TagVERB | TagAUX},

		// Ambiguous adjective tags.
		{"JJ", 0, 0, TagADJ | TagNOUN | TagDET | TagADV}, // 4 bits
		{"JJR", 2, 0, TagADJ | TagADV},
		{"JJS", 2, 0, TagADJ | TagADV},

		// Adverbs.
		{"RB", 2, 0, TagADV | TagPART}, // "not/n't" → PART in UD
		{"RBR", 1, TagADV, 0},
		{"RBS", 2, 0, TagADV | TagDET}, // "most" → DET in UD

		// Interrogative/relative adverbs — when/where as subordinators → SCONJ.
		{"WRB", 2, 0, TagADV | TagSCONJ},

		// Determiner/pronoun ambiguity.
		{"WDT", 2, 0, TagDET | TagPRON},

		// Preposition/subordinator ambiguity.
		{"IN", 2, 0, TagADP | TagSCONJ},

		// Infinitive marker.
		{"TO", 2, 0, TagPART | TagADP},

		// Unknown Penn tag → zero.
		{"UNKNOWN", 0, 0, 0},
		{"", 0, 0, 0},
	}

	for _, tc := range cases {
		got := TagFromPennTag(tc.penn)
		gotBits := bits.OnesCount32(uint32(got))

		if tc.wantBits == 0 && tc.mustHave == 0 {
			// Expecting zero (unknown Penn tag).
			if got != 0 {
				t.Errorf("TagFromPennTag(%q) = %v, want 0", tc.penn, got)
			}
			continue
		}

		if tc.wantBits == 1 {
			if got != tc.wantTag {
				t.Errorf("TagFromPennTag(%q) = %v, want %v", tc.penn, got, tc.wantTag)
			}
		} else {
			// Ambiguous: check bit count and required bits.
			if gotBits < 2 {
				t.Errorf("TagFromPennTag(%q): got %v (%d bits), want ambiguous (≥2 bits)", tc.penn, got, gotBits)
			}
			if got&tc.mustHave != tc.mustHave {
				t.Errorf("TagFromPennTag(%q) = %v, must include %v", tc.penn, got, tc.mustHave)
			}
		}
	}
}
