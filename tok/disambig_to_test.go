package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateTo(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// PART: infinitive after AUX (prev=unambiguous AUX → PART)
		// Note: "to VERB" where VERB is ambiguous with NOUN is resolved post-chunk
		// by DisambiguateByChunk — see disambig_chunk_test.go for those cases.
		// Here we test the prev=AUX rule with an unambiguous auxiliary.
		{"She ought to go.", "to", chunky.TagPART},

		// ADP: before DET
		{"He went to the store.", "to", chunky.TagADP},
		{"She drove to the city.", "to", chunky.TagADP},

		// ADP: before PRON
		{"Give it to him.", "to", chunky.TagADP},
		{"She spoke to them.", "to", chunky.TagADP},

		// ADP: before PROPN
		{"He traveled to London.", "to", chunky.TagADP},
		{"She moved to Paris.", "to", chunky.TagADP},

		// ADP: before NUM
		{"Up to five people.", "to", chunky.TagADP},

		// ADP: before resolved pure NOUN
		{"He went to war.", "to", chunky.TagADP},
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
