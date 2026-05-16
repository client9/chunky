package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateMore(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADV: degree modifier before adjective or adverb
		{"The result was more important.", "more", chunky.TagADV},
		{"She ran more quickly.", "more", chunky.TagADV},
		{"This is most likely.", "most", chunky.TagADV},
		{"He worked much harder.", "much", chunky.TagADV},
		{"She arrived less often.", "less", chunky.TagADV},

		// DET: quantifier before noun
		{"More people attended.", "More", chunky.TagDET},
		{"Most cities have parks.", "Most", chunky.TagDET},
		{"She had much time.", "much", chunky.TagDET},
		{"Less effort was needed.", "Less", chunky.TagDET},
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
