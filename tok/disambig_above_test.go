package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAbove(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: followed by DET/NOUN/PROPN/NUM
		{"She walked outside the building.", "outside", chunky.TagADP},
		{"The plane flew above the clouds.", "above", chunky.TagADP},
		{"He hid inside the house.", "inside", chunky.TagADP},
		{"Temperatures rose above 40 degrees.", "above", chunky.TagADP},

		// ADV: after pure verb
		{"She stayed outside.", "outside", chunky.TagADV},
		{"The balloon rose above.", "above", chunky.TagADV},
		{"He waited inside.", "inside", chunky.TagADV},
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
