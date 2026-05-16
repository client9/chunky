package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateOwn(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// possessive pronoun → ADJ
		{"his own house", "own", chunky.TagADJ},
		{"my own home", "own", chunky.TagADJ},
		{"their own decision", "own", chunky.TagADJ},
		{"her own car", "own", chunky.TagADJ},
		{"its own merits", "own", chunky.TagADJ},
		{"our own land", "own", chunky.TagADJ},
		{"your own choice", "own", chunky.TagADJ},

		// possessive pronoun + live → ADJ
		{"their live broadcast", "live", chunky.TagADJ},

		// possessive pronoun + separate → ADJ
		{"my separate account", "separate", chunky.TagADJ},

		// DET + own already resolved by context rules
		{"the own goal", "own", chunky.TagADJ},

		// PRON + own + DET → VERB (context rules handle)
		{"they own the house", "own", chunky.TagVERB},
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
