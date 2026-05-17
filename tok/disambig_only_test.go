package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateOnly(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// only: ADV before VERB/AUX
		{"She only runs on weekends.", "only", chunky.TagADV},
		{"He can only try.", "only", chunky.TagADV},

		// only: ADV before ADV
		{"It took only about an hour.", "only", chunky.TagADV},

		// only: ADV before DET
		{"He kept only the essentials.", "only", chunky.TagADV},

		// only: ADV before NUM
		{"She won only three times.", "only", chunky.TagADV},

		// only: ADV before ADP
		{"It works only for small files.", "only", chunky.TagADV},

		// only: ADJ prenominal after DET
		{"She was the only survivor.", "only", chunky.TagADJ},

		// little: DET before NOUN
		{"She had little time left.", "little", chunky.TagDET},

		// little: ADV before ADJ
		{"He was a little tired.", "little", chunky.TagADV},
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
