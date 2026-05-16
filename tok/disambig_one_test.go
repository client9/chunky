package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateOne(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// NUM: before ADP
		{"One of the students won.", "One", chunky.TagNUM},
		{"He scored one of his best goals.", "one", chunky.TagNUM},

		// NUM: before NOUN
		{"She took one step forward.", "one", chunky.TagNUM},
		{"One day she returned.", "One", chunky.TagNUM},

		// NUM: before ADJ
		{"Only one large section remained.", "one", chunky.TagNUM},

		// NUM: before NUM
		{"Chapter one, section one begins here.", "one", chunky.TagNUM},

		// NUM: before CCONJ
		{"Choose one or the other.", "one", chunky.TagNUM},

		// NUM: before PUNCT
		{"She chose option one.", "one", chunky.TagNUM},
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
