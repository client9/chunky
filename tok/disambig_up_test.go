package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateUp(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: next=DET
		{"they walked up the hill .", "up", chunky.TagADP},
		// ADP: next=NOUN
		{"costs scaled up production .", "up", chunky.TagADP},
		// ADP: next=PART ("up to")
		{"he lived up to expectations .", "up", chunky.TagADP},
		// ADP: prev=VERB (phrasal particle)
		{"she picked up the package .", "up", chunky.TagADP},
		{"the company scaled up last year .", "up", chunky.TagADP},
		// ADV: prev=AUX
		{"they will end up alone .", "up", chunky.TagADV},
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
