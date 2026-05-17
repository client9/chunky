package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateFree(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: prenominal
		{"free speech is protected .", "free", chunky.TagADJ},
		// ADJ: predicative
		{"the service is free .", "free", chunky.TagADJ},
		// ADJ: before PUNCT
		{"admission is free , no tickets required .", "free", chunky.TagADJ},
		// ADJ: before ADP ("free of charge")
		{"the room was free of charge .", "free", chunky.TagADJ},
		// ADJ: before PART ("free to go")
		{"you are free to leave .", "free", chunky.TagADJ},
		// VERB: next=DET
		{"they fought to free the prisoners .", "free", chunky.TagVERB},
		// VERB: next=PRON
		{"free them immediately .", "free", chunky.TagVERB},
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
