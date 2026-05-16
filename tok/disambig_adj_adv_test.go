package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAdjAdv(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// prior: ADV before ADP ("prior to")
		{"Prior to the meeting, she prepared.", "Prior", chunky.TagADV},
		{"He resigned prior to the announcement.", "prior", chunky.TagADV},

		// prior: ADJ before noun
		{"She had no prior experience.", "prior", chunky.TagADJ},
		{"The prior arrangement was cancelled.", "prior", chunky.TagADJ},

		// likely: ADJ before PART ("likely to V")
		{"She is likely to win.", "likely", chunky.TagADJ},
		{"It is likely to rain.", "likely", chunky.TagADJ},

		// later: ADV before VERB
		{"She later regretted the decision.", "later", chunky.TagADV},

		// later: ADJ before noun
		{"A later version fixed the bug.", "later", chunky.TagADJ},

		// early: ADV before ADP
		{"They arrived early in the morning.", "early", chunky.TagADV},

		// early: ADJ before noun
		{"The early results were promising.", "early", chunky.TagADJ},

		// late: ADV before ADP
		{"He submitted late in the term.", "late", chunky.TagADV},

		// hard: ADJ before noun
		{"It was a hard decision.", "hard", chunky.TagADJ},

		// far: ADJ before noun
		{"They reached the far shore.", "far", chunky.TagADJ},

		// earlier: ADV before VERB / ADJ before noun
		{"She had earlier announced her plans.", "earlier", chunky.TagADV},
		{"An earlier version was simpler.", "earlier", chunky.TagADJ},

		// longer: ADV before VERB / ADJ before noun
		{"She no longer works there.", "longer", chunky.TagADV},
		{"The longer route was scenic.", "longer", chunky.TagADJ},

		// further: ADV before VERB / ADJ before noun
		{"The investigation further revealed flaws.", "further", chunky.TagADV},
		{"Further details are available.", "Further", chunky.TagADJ},

		// short: ADJ before noun
		{"It was a short trip.", "short", chunky.TagADJ},

		// higher: ADJ before noun
		{"Higher taxes were proposed.", "Higher", chunky.TagADJ},
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
