package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateVerbForms(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// VERB: subject is NOUN
		{"The committee says the plan is ready.", "says", chunky.TagVERB},
		{"The price rose sharply.", "rose", chunky.TagVERB},
		{"The stock fell yesterday.", "fell", chunky.TagVERB},
		{"The company needs more funding.", "needs", chunky.TagVERB},
		{"The law means everyone must comply.", "means", chunky.TagVERB},

		// VERB: subject is PROPN
		{"Smith says the deal is done.", "says", chunky.TagVERB},

		// VERB: subject is PRON
		{"She remains the director.", "remains", chunky.TagVERB},
		{"He leads the team.", "leads", chunky.TagVERB},
		{"It takes time.", "takes", chunky.TagVERB},
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
