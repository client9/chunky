package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAdjVerb(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// VERB: after AUX (passive/perfect)
		{"The pilot was experienced.", "experienced", chunky.TagVERB},
		{"The flag was lowered.", "lowered", chunky.TagVERB},
		{"The report is marked confidential.", "marked", chunky.TagVERB},

		// ADJ: prenominal after DET
		{"The experienced pilot landed safely.", "experienced", chunky.TagADJ},
		{"A dry run confirmed the plan.", "dry", chunky.TagADJ},
		{"The marked difference was clear.", "marked", chunky.TagADJ},

		// ADJ: prenominal after ADJ
		{"The highly experienced team won.", "experienced", chunky.TagADJ},
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
