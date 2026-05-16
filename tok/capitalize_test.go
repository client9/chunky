package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestRetagCapitalized(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// Mid-sentence capitalized known word → PROPN.
		{"I visited London yesterday.", "London", chunky.TagPROPN},
		// Sentence-initial capitalized word: tag unchanged (VERB stays VERB).
		{"Walked the dog.", "Walked", chunky.TagVERB},
		// "I" is PRON and must never be promoted to PROPN.
		{"She and I went.", "I", chunky.TagPRON},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var found *Token
		for i := range sents {
			for j := range sents[i].Tokens {
				if sents[i].Tokens[j].Word == tc.word {
					found = &sents[i].Tokens[j]
				}
			}
		}
		if found == nil {
			t.Errorf("Parse(%q): token %q not found", tc.input, tc.word)
			continue
		}
		if found.IsUnknownTag() || !found.HasTag(tc.want) {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, found.Tags, tc.want)
		}
	}
}
