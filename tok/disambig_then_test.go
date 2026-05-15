package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateThen(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADV: default temporal/sequential use
		{"Outdoor Advertising, then worth a million.", "then", chunky.TagADV},
		{"He left, then returned.", "then", chunky.TagADV},
		{"She studied hard, then passed.", "then", chunky.TagADV},
		{"Then he arrived.", "Then", chunky.TagADV},
		// ADJ: pre-nominal after DET ("the then X")
		{"the then president of the company", "then", chunky.TagADJ},
		{"the then prime minister", "then", chunky.TagADJ},
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
