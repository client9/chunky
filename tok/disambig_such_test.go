package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateSuch(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: prenominal before DET
		{"Such a case requires care.", "Such", chunky.TagADJ},
		{"She had such a problem.", "such", chunky.TagADJ},

		// ADJ: prenominal before NOUN
		{"Such cases are rare.", "Such", chunky.TagADJ},
		{"He showed such courage.", "such", chunky.TagADJ},

		// ADJ: prenominal before ADJ
		{"Such extreme conditions were unusual.", "Such", chunky.TagADJ},

		// "such as" is merged to ADP by MergeLexical before this runs
		{"Fruits such as apples are healthy.", "such as", chunky.TagADP},
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
