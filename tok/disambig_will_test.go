package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateWill(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// AUX: modal before main verb
		{"He will go tomorrow.", "will", chunky.TagAUX},
		{"She will return.", "will", chunky.TagAUX},
		{"They will win.", "will", chunky.TagAUX},
		// AUX: modal before auxiliary (perfect/passive chains)
		{"She will be there.", "will", chunky.TagAUX},
		{"It will have ended.", "will", chunky.TagAUX},
		// AUX: modal separated by adverb
		{"He will not go.", "will", chunky.TagAUX},
		{"She will never return.", "will", chunky.TagAUX},
		{"They will also attend.", "will", chunky.TagAUX},
		// AUX: modal before PART ("will not" where not=PART)
		{"He will not leave.", "will", chunky.TagAUX},
		// AUX: interrogative inversion
		{"Will he attend?", "Will", chunky.TagAUX},
		{"Will they agree?", "Will", chunky.TagAUX},

		// NOUN: preceded by DET
		{"He read the will aloud.", "will", chunky.TagNOUN},
		{"She contested a will in court.", "will", chunky.TagNOUN},
		{"The will was disputed.", "will", chunky.TagNOUN},
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
