package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateHalf(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// final: ADJ before noun
		{"The final answer was correct.", "final", chunky.TagADJ},
		{"She reached the final round.", "final", chunky.TagADJ},

		// half: DET before DET/PRON ("half the X", "half a X")
		{"Half the city was flooded.", "Half", chunky.TagDET},
		{"He ran half a mile.", "half", chunky.TagDET},

		// half: NOUN before ADP ("half of X")
		{"The first half of the season was strong.", "half", chunky.TagNOUN},
		{"She spent half of her savings.", "half", chunky.TagNOUN},

		// individual: ADJ before noun
		{"Individual rights must be protected.", "Individual", chunky.TagADJ},
		{"Each individual case was reviewed.", "individual", chunky.TagADJ},

		// individual: NOUN before PART ("individual to V")
		{"Every individual to register must pay.", "individual", chunky.TagNOUN},
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
