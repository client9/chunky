package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateDetPron(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// each: next=ADP → PRON
		{"each of them", "each", chunky.TagPRON},
		// each: floating after subject → PRON
		{"they each have a role", "each", chunky.TagPRON},
		// each: before noun → DET
		{"each team played well", "each", chunky.TagDET},

		// some: next=ADP → PRON
		{"some of us agree", "some", chunky.TagPRON},
		// some: before noun → DET
		{"some water spilled", "some", chunky.TagDET},

		// any: next=ADP → PRON
		{"any of these works", "any", chunky.TagPRON},
		// any: before noun → DET
		{"any questions welcome", "any", chunky.TagDET},

		// this: before noun → DET
		{"This decision affected everyone.", "This", chunky.TagDET},
		{"I made this choice.", "this", chunky.TagDET},
		// this: before resolved AUX → PRON
		{"This is important.", "This", chunky.TagPRON},
		{"This was unexpected.", "This", chunky.TagPRON},

		// these: before unambiguous noun → DET
		{"These findings support the theory.", "These", chunky.TagDET},
		// these: before resolved AUX → PRON
		{"These are done.", "These", chunky.TagPRON},

		// those: before noun → DET
		{"Those players performed well.", "Those", chunky.TagDET},
		// those: before resolved AUX → PRON
		{"Those were difficult times.", "Those", chunky.TagPRON},

		// another: before noun → DET
		{"Another day passed.", "Another", chunky.TagDET},
		// another: next=ADP → PRON
		{"Another of the reports was missing.", "Another", chunky.TagPRON},

		// what: before unambiguous noun → DET
		{"What year was this?", "What", chunky.TagDET},
		// what: before resolved verb → PRON
		{"What happened next?", "What", chunky.TagPRON},
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
