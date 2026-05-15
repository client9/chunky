package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestHyphenCompoundParticiple(t *testing.T) {
	cases := []struct {
		word string
		want chunky.Tag
	}{
		// Past-participle compounds → ADJ
		{"environmental-themed", chunky.TagADJ},
		{"evidence-based", chunky.TagADJ},
		{"computer-aided", chunky.TagADJ},
		{"open-minded", chunky.TagADJ},
		{"red-haired", chunky.TagADJ},
		{"ill-fated", chunky.TagADJ},
		// Present-participle compounds → ADJ
		{"forward-looking", chunky.TagADJ},
		{"fast-moving", chunky.TagADJ},
		{"record-breaking", chunky.TagADJ},
		// Existing adj-suffix rules still work
		{"life-like", chunky.TagADJ},
		{"tax-free", chunky.TagADJ},
	}
	for _, tc := range cases {
		tags, _ := HyphenCandidates(tc.word)
		if len(tags) == 0 {
			t.Errorf("HyphenCandidates(%q): no tags returned", tc.word)
			continue
		}
		if len(tags) != 1 || tags[0] != tc.want {
			t.Errorf("HyphenCandidates(%q) = %v, want [%v]", tc.word, tags, tc.want)
		}
	}
}

func TestInflectionNoPronoun(t *testing.T) {
	// "themed" must not return PRON (false stem "them") after filtering
	tags, _ := InflectionCandidates("themed")
	for _, tag := range tags {
		if tag == chunky.TagPRON {
			t.Errorf("InflectionCandidates(%q): got PRON, want only open-class tags", "themed")
		}
	}
}
