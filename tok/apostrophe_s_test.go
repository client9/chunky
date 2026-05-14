package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateApostropheS(t *testing.T) {
	cases := []struct {
		input string
		want  chunky.Tag
	}{
		// copula / auxiliary
		{"it's", chunky.TagAUX},
		{"he's", chunky.TagAUX},
		{"she's", chunky.TagAUX},
		{"that's", chunky.TagAUX},
		{"what's", chunky.TagAUX},
		{"there's", chunky.TagAUX},
		{"here's", chunky.TagAUX},
		{"who's", chunky.TagAUX},
		{"how's", chunky.TagAUX},
		// possessive
		{"John's", chunky.TagPART},
		{"company's", chunky.TagPART},
		{"everyone's", chunky.TagPART},
		{"one's", chunky.TagPART},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var tokens []Token
		for _, s := range sents {
			tokens = append(tokens, s.Tokens...)
		}
		var apostropheS *Token
		for i := range tokens {
			if tokens[i].Word == "'s" {
				apostropheS = &tokens[i]
				break
			}
		}
		if apostropheS == nil {
			t.Errorf("Parse(%q): no \"'s\" token found", tc.input)
			continue
		}
		if len(apostropheS.Tags) == 0 || apostropheS.Tags[0] != tc.want {
			t.Errorf("Parse(%q) \"'s\": got %v, want %v", tc.input, apostropheS.Tags, tc.want)
		}
	}
}
