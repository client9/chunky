package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateHigh(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: prenominal before NOUN
		{"a high speed chase .", "high", chunky.TagADJ},
		{"high temperatures were recorded .", "high", chunky.TagADJ},
		// ADJ: before ADJ
		{"a high powered engine .", "high", chunky.TagADJ},
		// ADJ: before CCONJ
		{"prices were high and rising .", "high", chunky.TagADJ},
		// ADJ: before PUNCT
		{"the score was surprisingly high .", "high", chunky.TagADJ},
		// ADJ: after DET
		{"the high cost of living .", "high", chunky.TagADJ},
		// ADJ: after ADP
		{"at high tide the water rises .", "high", chunky.TagADJ},
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
