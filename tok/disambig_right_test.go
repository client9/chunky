package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateRight(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: prenominal before NOUN
		{"the right approach was taken .", "right", chunky.TagADJ},
		{"make the right choice .", "right", chunky.TagADJ},
		// ADJ: predicative after AUX
		{"she is right about that .", "right", chunky.TagADJ},
		{"you were right all along .", "right", chunky.TagADJ},
		// ADV: before ADV ("right now", "right away", "right here")
		{"right now is the time .", "right", chunky.TagADV},
		{"he left right away .", "right", chunky.TagADV},
		// NOUN: before PART ("right to X")
		{"everyone has the right to vote .", "right", chunky.TagNOUN},
		{"the right to remain silent .", "right", chunky.TagNOUN},
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
