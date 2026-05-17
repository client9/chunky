package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateByChunkADPPART(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// Infinitival "to" inside VP → PART
		{"ranches to re-popularize bison meat", "to", chunky.TagPART},
		{"They plan to expand operations.", "to", chunky.TagPART},
		// Prepositional "to" → ADP
		{"He went to the store.", "to", chunky.TagADP},
		{"She drove to the city.", "to", chunky.TagADP},
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

func TestDisambiguateByChunkNounVerb(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// NOUN/VERB inside NP → NOUN (resolved by chunk position)
		{"the bison herd in the world", "herd", chunky.TagNOUN},
		{"of cable television", "cable", chunky.TagNOUN},

		// VP must not swallow a NOUN/VERB object after a non-auxiliary VERB.
		// "made plans" → plans is direct object, not part of VP.
		{"Noriega made plans to fly", "plans", chunky.TagNOUN},
		{"the team made record profits this year", "record", chunky.TagNOUN},
		{"Lilly attributed record third-quarter results", "record", chunky.TagNOUN},

		// Two consecutive {NOUN,VERB} tokens: second must not extend VP.
		// "had operating profit" → operating extends VP, profit is NP head.
		{"He had operating profit of ten million", "profit", chunky.TagNOUN},
		{"it had operating profit of ten million", "profit", chunky.TagNOUN},
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
