package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateThat(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// prev=VERB, next=DET → complementizer
		{"He said that the car was fast.", "that", chunky.TagSCONJ},
		{"She knew that a solution existed.", "that", chunky.TagSCONJ},
		// prev=VERB, next=ambiguous-DET ("more") → complementizer
		{"the report said that more borrowers defaulted.", "that", chunky.TagSCONJ},
		// prev=ADJ, next=PRON → complementizer
		{"She was confident that he would agree.", "that", chunky.TagSCONJ},
		// appositive: prev=NOUN, next=PRON → complementizer
		{"the fact that he lied shocked everyone.", "that", chunky.TagSCONJ},
		{"evidence that she had lied emerged.", "that", chunky.TagSCONJ},
		// next=DET (standalone DET trigger, prev irrelevant)
		{"the claim that the vote was rigged", "that", chunky.TagSCONJ},
		// "that" before NOUN at sentence start → DET (demonstrative)
		{"That car is nice.", "That", chunky.TagDET},
		// "that" before PRON → still ambiguous (anaphoric "I see that")
		{"I see that.", "that", chunky.TagPRON},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			if tc.want != chunky.TagPRON {
				t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			}
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
