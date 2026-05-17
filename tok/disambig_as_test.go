package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAs(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: role/capacity after verb, next=DET (prev≠PUNCT)
		{"he served as the president .", "as", chunky.TagADP},
		{"she is regarded as an expert .", "as", chunky.TagADP},
		{"it was seen as a problem .", "as", chunky.TagADP},
		// ADP: sentence-initial before DET (prev=sentence boundary, not PUNCT)
		{"As the storm approached , they fled .", "As", chunky.TagADP},
		{"As a pioneer , she led the way .", "As", chunky.TagADP},
		// ADP: before NOUN (no article) — role assignment
		{"he served as chairman of the board .", "as", chunky.TagADP},
		{"she was known as Smith .", "as", chunky.TagADP},

		// ADV: after PUNCT, before ADV (", as yet")
		{"the details are not clear , as yet .", "as", chunky.TagADV},

		// SCONJ: after PUNCT, before PRON (", as he said")
		{"the plan failed , as he predicted .", "as", chunky.TagSCONJ},
		{"it failed , as she expected .", "as", chunky.TagSCONJ},
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
