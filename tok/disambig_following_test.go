package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateFollowing(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP (prepositional): before DET
		{"Following the announcement, trading halted.", "Following", chunky.TagADP},
		{"The protest erupted following the decision.", "following", chunky.TagADP},

		// ADP (prepositional): before PRON
		{"Following his resignation, the board met.", "Following", chunky.TagADP},
		{"Following their victory, the team celebrated.", "Following", chunky.TagADP},

		// NOUN: before AUX ("the following was/were")
		{"The following was announced.", "following", chunky.TagNOUN},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if tc.want == 0 {
			if resolved {
				t.Errorf("Parse(%q) %q: expected ambiguous, got resolved %v", tc.input, tc.word, got)
			}
			continue
		}
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
