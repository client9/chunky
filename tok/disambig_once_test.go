package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateOnce(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// once: ADV — preceded by auxiliary
		{"He had once been famous.", "once", chunky.TagADV},
		{"She was once a teacher.", "once", chunky.TagADV},
		{"They have once again failed.", "once", chunky.TagADV},

		// once: ADV — followed by adverb ("once more", "once again")
		{"Try once more.", "once", chunky.TagADV},
		{"He tried once again.", "once", chunky.TagADV},

		// once: ADV — followed by ADP ("once upon", "once in", "once before")
		{"Once upon a time there was a king.", "Once", chunky.TagADV},
		{"It happens once in a while.", "once", chunky.TagADV},

		// once: ADV — frequentive "once a X"
		{"She visits once a month.", "once", chunky.TagADV},

		// once: ADV — sentence-final
		{"He visited once.", "once", chunky.TagADV},

		// once: SCONJ — temporal clause
		{"Once the war ended, trade resumed.", "Once", chunky.TagSCONJ},
		{"Once they arrived, the meeting began.", "Once", chunky.TagSCONJ},
		{"Once completed, the project was reviewed.", "Once", chunky.TagSCONJ},
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
