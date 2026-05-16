package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateMine(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// PRON: after AUX
		{"The victory is mine.", "mine", chunky.TagPRON},
		{"That car was mine.", "mine", chunky.TagPRON},

		// NOUN: after DET
		{"Workers sealed the mine.", "mine", chunky.TagNOUN},
		{"They found a mine on the road.", "mine", chunky.TagNOUN},

		// NOUN: before ADP
		{"The mine of coal was exhausted.", "mine", chunky.TagNOUN},

		// NOUN: before NOUN (compound)
		{"The mine shaft collapsed.", "mine", chunky.TagNOUN},
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
