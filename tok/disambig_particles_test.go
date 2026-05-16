package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateParticles(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: verb-particle constructions
		{"She picked up the phone.", "up", chunky.TagADP},
		{"He gave up the fight.", "up", chunky.TagADP},
		{"They blew up the bridge.", "up", chunky.TagADP},
		{"She turned off the lights.", "off", chunky.TagADP},
		{"He cut off the supply.", "off", chunky.TagADP},
		{"They called off the search.", "off", chunky.TagADP},

		// next=CCONJ → ADV
		{"Up and away.", "Up", chunky.TagADV},
		// next=NOUN → ADP
		{"Up the hill they ran.", "Up", chunky.TagADP},
		// next=PUNCT → ADV (no object: intransitive particle)
		{"Prices went up.", "up", chunky.TagADV},

		// off: ADV before PUNCT
		{"The alarm went off.", "off", chunky.TagADV},
		// off: ADV before CCONJ
		{"The deal fell off and died.", "off", chunky.TagADV},

		// up: ADV before CCONJ
		{"Stand up and fight.", "up", chunky.TagADV},
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
