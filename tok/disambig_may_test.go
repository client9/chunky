package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateMay(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// Lowercase "may" → always AUX
		{"We may file the report.", "may", chunky.TagAUX},
		{"She may not attend.", "may", chunky.TagAUX},
		{"They may have left.", "may", chunky.TagAUX},
		// Mid-sentence "May" → PROPN (handled by RetagCapitalized)
		{"We filed in May.", "May", chunky.TagPROPN},
		{"The report is due in May 2024.", "May", chunky.TagPROPN},
		// "May" before NUM → PROPN (date)
		{"The filing is due May 15.", "May", chunky.TagPROPN},
		{"May 1 is a holiday.", "May", chunky.TagPROPN},
		// Sentence-initial "May" before PRON → AUX (interrogative inversion)
		{"May I help you?", "May", chunky.TagAUX},
		{"May we proceed?", "May", chunky.TagAUX},
		// Sentence-initial "May" before non-PRON → PROPN (month)
		{"May flowers bloom early.", "May", chunky.TagPROPN},
		{"May brings warm weather.", "May", chunky.TagPROPN},
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
