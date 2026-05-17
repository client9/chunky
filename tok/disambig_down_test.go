package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateDown(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// down: ADP before DET/NOUN/PROPN
		{"She walked down the stairs.", "down", chunky.TagADP},
		{"Drive down Main Street.", "down", chunky.TagADP},

		// down: ADV — particle after verb, or bare number (financial "closed down 1.05")
		{"The system went down.", "down", chunky.TagADV},
		{"The car broke down.", "down", chunky.TagADV},
		{"She counted down 10 seconds.", "down", chunky.TagADV}, // phrasal verb particle

		// near: ADP before DET/NOUN/PROPN
		{"She sat near the window.", "near", chunky.TagADP},
		{"He moved near London.", "near", chunky.TagADP},

		// near: ADV before ADJ/ADV
		{"The shot was near perfect.", "near", chunky.TagADV},
		{"The result was near perfect.", "near", chunky.TagADV},

		// near: ADV after pure verb
		{"He stood near.", "near", chunky.TagADV},
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
