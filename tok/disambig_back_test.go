package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateBack(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// back: ADV after verb
		{"She went back to the office.", "back", chunky.TagADV},
		{"He came back yesterday.", "back", chunky.TagADV},
		{"Please step back.", "back", chunky.TagADV},

		// back: ADJ before noun (not after verb)
		{"The back door was open.", "back", chunky.TagADJ},
		{"She took the back seat.", "back", chunky.TagADJ},
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
