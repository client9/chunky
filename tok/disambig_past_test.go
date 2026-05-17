package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguatePast(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// past: ADJ before NOUN/NUM/ADJ
		{"She had no past experience.", "past", chunky.TagADJ},
		{"The past decade saw rapid growth.", "past", chunky.TagADJ},
		{"That was a past president.", "past", chunky.TagADJ},

		// past: NOUN before PUNCT or AUX
		{"He lived in the past.", "past", chunky.TagNOUN},
		{"The past was different.", "past", chunky.TagNOUN},

		// past: ADP before DET
		{"She drove past the school.", "past", chunky.TagADP},

	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
