package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateWell(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// well: ADV after verb
		{"She did well on the exam.", "well", chunky.TagADV},
		{"The plan worked well.", "well", chunky.TagADV},

		// well: NOUN after determiner
		{"The old well ran dry.", "well", chunky.TagNOUN},
		{"They dug a well.", "well", chunky.TagNOUN},
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
