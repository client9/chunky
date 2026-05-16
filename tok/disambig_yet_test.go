package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateYet(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADV: followed by ADV/ADJ/DET/NUM
		{"Yet another attempt failed.", "Yet", chunky.TagADV},
		{"She had not yet finished.", "yet", chunky.TagADV},
		{"It is not yet complete.", "yet", chunky.TagADV},
		{"Yet three more problems arose.", "Yet", chunky.TagADV},

		// ADV: preceded by AUX or VERB
		{"He has not yet arrived.", "yet", chunky.TagADV},
		{"She arrived yet.", "yet", chunky.TagADV},

		// ADV: before ADJ — fires even in adversative position; CCONJ not resolvable by single-token lookahead
		{"The task was simple, yet challenging.", "yet", chunky.TagADV},
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
