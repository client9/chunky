package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateLong(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: next=NOUN
		{"it was a long time ago .", "long", chunky.TagADJ},
		{"they took the long road .", "long", chunky.TagADJ},
		// ADJ: prev=DET
		{"the long debate continued .", "long", chunky.TagADJ},
		// ADV: next=AUX
		{"how long will this take ?", "long", chunky.TagADV},
		// ADV: next=ADV
		{"not long ago they arrived .", "long", chunky.TagADV},
		// ADV: next=ADV ("so long ago")
		{"so long ago they lived there .", "long", chunky.TagADV},
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
