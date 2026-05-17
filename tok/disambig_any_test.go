package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAny(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// DET: next=NOUN
		{"any suggestion is welcome .", "any", chunky.TagDET},
		{"any team can enter .", "any", chunky.TagDET},
		// DET: next=NUM
		{"any 10 items qualify .", "any", chunky.TagDET},
		// PRON: next=ADP ("any of")
		{"any of the options will work .", "any", chunky.TagPRON},
		// PRON: next=PUNCT ("if any.")
		{"there were no changes , if any .", "any", chunky.TagPRON},
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
