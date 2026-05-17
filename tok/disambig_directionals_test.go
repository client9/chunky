package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateDirectionals(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// south: ADV when followed by ADP
		{"The army moved south of the river.", "south", chunky.TagADV},
		{"They camped south of the border.", "south", chunky.TagADV},

		// south: NOUN when preceded by DET
		{"The south was unprepared.", "south", chunky.TagNOUN},
		{"She grew up in the south.", "south", chunky.TagNOUN},

		// north: NOUN when preceded by DET
		{"The north held firm.", "north", chunky.TagNOUN},
		{"Troops advanced from the north.", "north", chunky.TagNOUN},

		// east: NOUN when preceded by DET
		{"Trade routes crossed the east.", "east", chunky.TagNOUN},

		// west: NOUN when preceded by DET
		{"The west was settled last.", "west", chunky.TagNOUN},

		// prenominal ADJ: next=NOUN blocks our DET→NOUN rule; context rules resolve to ADJ
		{"The south side of town was quiet.", "south", chunky.TagADJ},
		// prenominal ADJ: resolves to ADJ via adjNounBroadRule (before NOUN) or chunk context
		{"A north wind blew in.", "north", chunky.TagADJ},

		// directionals: NOUN when followed by ADP "of"
		{"100 km northwest of the city .", "northwest", chunky.TagNOUN},
		{"Ingolstadt , northwest of Landshut .", "northwest", chunky.TagNOUN},
		{"It lies north of the river .", "north", chunky.TagNOUN},
		{"The town sits east of the highway .", "east", chunky.TagNOUN},

		// compound directionals: NOUN when preceded by DET
		{"The northwest was hit hardest.", "northwest", chunky.TagNOUN},
		{"Storms swept through the northeast.", "northeast", chunky.TagNOUN},
		{"Trade routes crossed the southeast.", "southeast", chunky.TagNOUN},
		{"Rain fell across the southwest.", "southwest", chunky.TagNOUN},

		// compound directionals: PROPN when part of a proper name (sentence-initial caps chain)
		{"Southwest Airlines reported a loss.", "Southwest", chunky.TagPROPN},
		{"Southeast Asia saw rapid growth.", "Southeast", chunky.TagPROPN},
		{"Northeast Brazil was affected.", "Northeast", chunky.TagPROPN},
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
