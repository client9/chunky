package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAfter(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADP: gerund/participle complement (next=VERB)
		{"After graduating, she moved abroad.", "After", chunky.TagADP},
		{"He left before finishing his work.", "before", chunky.TagADP},
		{"They waited until receiving confirmation.", "until", chunky.TagADP},

		// ADP: date/year complement (next=NUM)
		{"The war ended after 1945.", "after", chunky.TagADP},
		{"Submit your application before 2026.", "before", chunky.TagADP},
		{"The contract runs until 2030.", "until", chunky.TagADP},

		// SCONJ: subject pronoun after conjunction → SCONJ
		{"He waited until they arrived.", "until", chunky.TagSCONJ},

		// ADP: noun/propn complement — all three conjunctions
		{"After midnight, the city quieted.", "After", chunky.TagADP},
		{"Finish before noon.", "before", chunky.TagADP},
		{"Wait until victory.", "until", chunky.TagADP},

		// ADP: object pronoun complement
		{"She arrived after him.", "after", chunky.TagADP},
		{"He left before them.", "before", chunky.TagADP},

		// context rules resolve DET-starting clauses to ADP
		{"After the war ended, trade resumed.", "After", chunky.TagADP},
		{"She left before the meeting started.", "before", chunky.TagADP},
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
