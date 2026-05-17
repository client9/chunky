package tok

import (
	"testing"

	"github.com/client9/chunky"
)

// tagOf returns the resolved tag of the first token in sents whose Word matches word.
func tagOf(sents []Sentence, word string) (chunky.Tag, bool) {
	for _, s := range sents {
		for _, tok := range s.Tokens {
			if tok.Word == word {
				if tok.IsResolved() {
					return tok.Tags, true
				}
				return 0, false // ambiguous
			}
		}
	}
	return 0, false // not found
}

func TestDisambiguateApostropheS(t *testing.T) {
	cases := []struct {
		input string
		want  chunky.Tag
	}{
		// copula / auxiliary
		{"it's", chunky.TagAUX},
		{"he's", chunky.TagAUX},
		{"she's", chunky.TagAUX},
		{"that's", chunky.TagAUX},
		{"what's", chunky.TagAUX},
		{"there's", chunky.TagAUX},
		{"here's", chunky.TagAUX},
		{"who's", chunky.TagAUX},
		{"how's", chunky.TagAUX},
		// possessive
		{"John's", chunky.TagPART},
		{"company's", chunky.TagPART},
		{"everyone's", chunky.TagPART},
		{"one's", chunky.TagPART},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var tokens []Token
		for _, s := range sents {
			tokens = append(tokens, s.Tokens...)
		}
		var apostropheS *Token
		for i := range tokens {
			if tokens[i].Word == "'s" {
				apostropheS = &tokens[i]
				break
			}
		}
		if apostropheS == nil {
			t.Errorf("Parse(%q): no \"'s\" token found", tc.input)
			continue
		}
		if apostropheS.Tags != tc.want {
			t.Errorf("Parse(%q) \"'s\": got %v, want %v", tc.input, apostropheS.Tags, tc.want)
		}
	}
}

// TestPossessiveApostrophe checks that a bare ' after a plural noun is tagged PART.
// Penn Treebank tokenizes "analysts'" as "analysts" + "'" (tag POS).
// UD always tags the possessive marker as PART.
func TestPossessiveApostrophe(t *testing.T) {
	cases := []struct {
		input string
		word  string // the word before '
	}{
		{"analysts ' expectations were high.", "analysts"},
		{"countries ' leaders met yesterday.", "countries"},
		{"creditors ' approval was needed.", "creditors"},
		{"the companies ' profits rose.", "companies"},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var apostrophe *Token
		for _, s := range sents {
			for i, tok := range s.Tokens {
				if tok.Word == "'" && i > 0 && s.Tokens[i-1].Word == tc.word {
					apostrophe = &s.Tokens[i]
					break
				}
			}
		}
		if apostrophe == nil {
			t.Errorf("Parse(%q): no \"'\" token after %q", tc.input, tc.word)
			continue
		}
		if !apostrophe.IsResolved() || apostrophe.Tags != chunky.TagPART {
			t.Errorf("Parse(%q) \"'\": got %v (resolved=%v), want PART",
				tc.input, apostrophe.Tags, apostrophe.IsResolved())
		}
	}
}

// TestPossessiveNeighbors checks that NOUN/VERB-ambiguous tokens adjacent to a
// possessive "'s" are resolved to NOUN. PART conflates possessive "'s" and
// infinitival "to", so their corpus statistics cancel out and no corpus-derived
// rule clears the 10x ratio threshold. These are handled directly in
// DisambiguateApostropheS instead.
func TestPossessiveNeighbors(t *testing.T) {
	cases := []struct {
		input string
		word  string // the NOUN/VERB-ambiguous word adjacent to "'s"
		want  chunky.Tag
	}{
		// Possessor (before "'s"): "father" is NOUN/VERB; should resolve to NOUN.
		{"after his father 's suicide", "father", chunky.TagNOUN},
		// Possessed head (after "'s"): "board" is NOUN/VERB; should resolve to NOUN.
		{"the Foundation 's board of directors", "board", chunky.TagNOUN},
		// Possessor that is a PROPN via caps should be unaffected (already unambiguous).
		{"after Turner 's death", "Turner", chunky.TagPROPN},
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
