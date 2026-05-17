package tok

import (
	"testing"

	"github.com/client9/chunky"
)

// TestDisambiguateContractionFragments checks that Penn-treebank contraction
// fragments "ca" and "wo" are resolved to AUX when followed by "n't" or "'t".
func TestDisambiguateContractionFragments(t *testing.T) {
	cases := []struct {
		input string
		word  string
	}{
		// Penn "can't" split: ca + n't
		{"I ca n't do it .", "ca"},
		// Penn "won't" split: wo + n't
		{"She wo n't be here .", "wo"},
		// Uppercase variants
		{"He Ca n't believe it .", "Ca"},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want AUX", tc.input, tc.word)
			continue
		}
		if got != chunky.TagAUX {
			t.Errorf("Parse(%q) %q: got %v, want AUX", tc.input, tc.word, got)
		}
	}
}
