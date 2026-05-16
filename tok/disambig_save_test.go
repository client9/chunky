package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateSave(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// save: VERB before DET
		{"Please save the file.", "save", chunky.TagVERB},
		// save: VERB before NOUN
		{"We must save lives.", "save", chunky.TagVERB},
		// save: VERB before PRON
		{"Can you save it?", "save", chunky.TagVERB},

		// respecting: ADP (formal preposition) before DET
		{"All rules respecting the process apply.", "respecting", chunky.TagADP},
		// respecting: ADP before NOUN
		{"Laws respecting privacy were enacted.", "respecting", chunky.TagADP},
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
