package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateQuantifiers(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// both: DET before noun phrase
		{"Both teams played well.", "Both", chunky.TagDET},
		{"She visited both cities.", "both", chunky.TagDET},

		// neither: DET before noun phrase
		{"Neither team scored.", "Neither", chunky.TagDET},
		{"Neither option was ideal.", "Neither", chunky.TagDET},

		// either: DET before noun phrase
		{"Either road leads there.", "Either", chunky.TagDET},
		{"Choose either option.", "either", chunky.TagDET},

		// all: DET before noun phrase
		{"All students passed.", "All", chunky.TagDET},
		{"She read all the books.", "all", chunky.TagDET},

		// both: DET before DET ("both the teams")
		{"Both the teams advanced.", "Both", chunky.TagDET},
		// both: DET before PROPN
		{"Both Germany and France agreed.", "Both", chunky.TagDET},
		// both: PRON — next=ADP ("both of them")
		{"Both of them agreed.", "Both", chunky.TagPRON},
		// both: PRON — next=AUX ("both were")
		{"Both were present.", "Both", chunky.TagPRON},

		// neither: PRON — next=ADP ("neither of them")
		{"Neither of them arrived.", "Neither", chunky.TagPRON},
		// neither: PRON — next=AUX
		{"Neither was correct.", "Neither", chunky.TagPRON},
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
