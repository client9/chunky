package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateSpatial(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: verb-particle constructions
		{"She looked out the window.", "out", chunky.TagADP},
		{"He fell behind the others.", "behind", chunky.TagADP},
		{"They stayed below deck.", "below", chunky.TagADP},
		{"She talked about the problem.", "about", chunky.TagADP},
		{"He walked around the block.", "around", chunky.TagADP},

		// ADV: approximation (about/around + NUM)
		{"About 50 people attended.", "About", chunky.TagADV},
		{"Around 100 years ago.", "Around", chunky.TagADV},
		{"The project cost about 200 dollars.", "about", chunky.TagADV},
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
