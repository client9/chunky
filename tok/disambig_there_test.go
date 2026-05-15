package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateThere(t *testing.T) {
	cases := []struct {
		input string
		word  string // "there" or "There"
		want  chunky.Tag
	}{
		// Existential: followed by AUX
		{"There is a problem.", "There", chunky.TagPRON},
		{"There are several reasons.", "There", chunky.TagPRON},
		{"There was no sign.", "There", chunky.TagPRON},
		{"There were no survivors.", "There", chunky.TagPRON},
		{"There will be consequences.", "There", chunky.TagPRON},
		// Existential mid-sentence
		{"He said there is no way.", "there", chunky.TagPRON},
		{"She knew there were problems.", "there", chunky.TagPRON},
		// Locative: after a verb
		{"He went there.", "there", chunky.TagADV},
		{"She is over there.", "there", chunky.TagADV},
		{"Put it there.", "there", chunky.TagADV},
		// Locative: at end / before punctuation
		{"We have been there before.", "there", chunky.TagADV},
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
