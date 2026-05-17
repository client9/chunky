package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateSo(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// so: ADV — intensifier before adjective or adverb
		{"The result was so good.", "so", chunky.TagADV},
		{"She ran so quickly.", "so", chunky.TagADV},
		{"It was so very cold.", "so", chunky.TagADV},

		// so: "so that" is merged by MergeLexical before this disambiguator runs;
		// the next=SCONJ rule handles any remaining so+SCONJ sequences.
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if tc.want == 0 {
			if resolved {
				t.Errorf("Parse(%q) %q: expected ambiguous, got resolved %v", tc.input, tc.word, got)
			}
			continue
		}
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
