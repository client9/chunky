package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateThat(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// "that" before DET → SCONJ (complementizer)
		{"He said that the car was fast.", "that", chunky.TagSCONJ},
		{"She knew that a solution existed.", "that", chunky.TagSCONJ},
		// "that" before NOUN → still ambiguous (leave as Tags[0])
		{"That car is nice.", "That", chunky.TagPRON},
		// "that" before PRON → still ambiguous
		{"I see that.", "that", chunky.TagPRON},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			// still ambiguous is OK for non-SCONJ cases
			if tc.want != chunky.TagPRON {
				t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			}
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
