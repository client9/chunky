package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguatePlus(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"three plus four equals seven", "plus", chunky.TagCCONJ},
		{"the cost plus tax", "plus", chunky.TagCCONJ},
		{"a bonus plus benefits", "plus", chunky.TagCCONJ},
		// capitalized in brand names resolved by RetagCapitalized
		{"Disney Plus launched", "Plus", chunky.TagPROPN},
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
