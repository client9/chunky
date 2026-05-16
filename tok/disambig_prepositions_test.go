package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguatePrepositions(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// about: ADP before DET/NOUN/PRON/PROPN
		{"She knew about the problem.", "about", chunky.TagADP},
		{"He asked about it.", "about", chunky.TagADP},
		{"The book is about Paris.", "about", chunky.TagADP},
		{"She worried about time.", "about", chunky.TagADP},

		// around: ADP before DET/NOUN/PROPN
		{"They walked around the lake.", "around", chunky.TagADP},
		{"The road runs around London.", "around", chunky.TagADP},

		// around: ADV before PUNCT
		{"She turned around.", "around", chunky.TagADV},

		// below: ADP before DET/NUM/NOUN
		{"Temperatures fell below zero.", "below", chunky.TagADP},
		{"She hid below the desk.", "below", chunky.TagADP},

		// below: ADV before PUNCT
		{"See the chart below.", "below", chunky.TagADV},

		// behind: ADP before DET/NOUN/PROPN
		{"He hid behind the curtain.", "behind", chunky.TagADP},
		{"The motive behind the decision was unclear.", "behind", chunky.TagADP},

		// out: ADP before DET/NOUN
		{"She looked out the window.", "out", chunky.TagADP},
		{"They moved out the furniture.", "out", chunky.TagADP},

		// out: ADV after pure verb
		{"The fire burned out.", "out", chunky.TagADV},
		{"He went out.", "out", chunky.TagADV},

		// around: ADV after pure verb
		{"She looked around.", "around", chunky.TagADV},

		// behind: ADV after pure verb
		{"The runner fell behind.", "behind", chunky.TagADV},
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
