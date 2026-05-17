package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateEach(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// DET: next=NOUN
		{"each player scored .", "each", chunky.TagDET},
		{"each team received a trophy .", "each", chunky.TagDET},
		// DET: next=ADJ
		{"each additional vote counts .", "each", chunky.TagDET},
		// PRON: next=ADP ("each of")
		{"each of the players scored .", "each", chunky.TagPRON},
		// PRON: next=VERB (floating quantifier)
		{"they each scored .", "each", chunky.TagPRON},
		// PRON: next=AUX
		{"they each will contribute .", "each", chunky.TagPRON},
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
