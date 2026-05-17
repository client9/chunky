package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateAdjNounDefault(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// ADJ: pre-nominal position
		{"The social media campaign launched.", "social", chunky.TagADJ},
		{"An unknown quantity remains.", "unknown", chunky.TagADJ},
		{"The red cars are popular.", "red", chunky.TagADJ},
		{"The white bread is common.", "white", chunky.TagADJ},
		{"The native population grew.", "native", chunky.TagADJ},
		{"The public opinion shifted.", "public", chunky.TagADJ},
		{"The general manager resigned.", "general", chunky.TagADJ},
		{"An equivalent amount was paid.", "equivalent", chunky.TagADJ},
		{"The evil empire fell.", "evil", chunky.TagADJ},
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

func TestDisambiguateChief(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"The chief executive resigned.", "chief", chunky.TagADJ},
		{"The chief engineer filed the report.", "chief", chunky.TagADJ},
		{"The editor in chief approved it.", "chief", chunky.TagNOUN},
		{"The commander in chief decided.", "chief", chunky.TagNOUN},
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

func TestDisambiguateSideNoun(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"The right side was damaged.", "side", chunky.TagNOUN},
		{"The other side agreed.", "side", chunky.TagNOUN},
		{"The side of the road was clear.", "side", chunky.TagNOUN},
		{"She moved to the side.", "side", chunky.TagNOUN},
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

func TestDisambiguateNearby(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"A nearby town was flooded.", "nearby", chunky.TagADJ},
		{"The nearby hospital closed.", "nearby", chunky.TagADJ},
		{"She lived nearby.", "nearby", chunky.TagADV},
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
