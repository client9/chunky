package tok

import (
	"testing"

	"github.com/client9/chunky"
)

// TestSentenceInitialOOVNoun checks that an unknown capitalized word at the
// start of a non-first sentence is not given the caps (PROPN|ADJ) signal.
// Sentence-initial capitalisation carries no POS information — "HFRS" and
// "Ribavirin" are common nouns, not proper names. TagUnknowns approximates
// sentence boundaries by checking for a preceding sentence-end punctuation
// mark, suppressing the MorphCandidates caps rule in that position.
func TestSentenceInitialOOVNoun(t *testing.T) {
	cases := []struct {
		input string
		word  string
	}{
		// Medical abbreviation — NOUN, not PROPN.
		{"Fever is common. HFRS is due to a virus.", "HFRS"},
		// Drug name used as a common noun in medical prose.
		{"Treatment varies. Ribavirin may be useful.", "Ribavirin"},
		// Invented unknown word at sentence start — caps signal suppressed.
		{"First sentence here. Zalbrutix was administered.", "Zalbrutix"},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want NOUN", tc.input, tc.word)
			continue
		}
		if got != chunky.TagNOUN {
			t.Errorf("Parse(%q) %q: got %v, want NOUN", tc.input, tc.word, got)
		}
	}
}

func TestSentenceInitialPropnChain(t *testing.T) {
	cases := []struct {
		input string
		word  string
	}{
		// OOV first name followed by a known capitalized name component.
		{"Robert Edward Turner III was born.", "Robert"},
		// OOV first name followed by an OOV last name that gets caps→PROPN.
		{"Zalbrutix Morwick arrived.", "Zalbrutix"},
		// Sentence-initial word resolved to ADJ followed by PROPN → PROPN.
		{"Great American Bank failed.", "Great"},
		{"Eastern Airlines reported losses.", "Eastern"},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want PROPN", tc.input, tc.word)
			continue
		}
		if got != chunky.TagPROPN {
			t.Errorf("Parse(%q) %q: got %v, want PROPN", tc.input, tc.word, got)
		}
	}
}

func TestRetagCapitalized(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// Mid-sentence capitalized known word → PROPN.
		{"I visited London yesterday.", "London", chunky.TagPROPN},
		// Sentence-initial capitalized word: tag unchanged (VERB stays VERB).
		{"Walked the dog.", "Walked", chunky.TagVERB},
		// "I" is PRON and must never be promoted to PROPN.
		{"She and I went.", "I", chunky.TagPRON},
		// Pure closed-class words must NOT be promoted even when capitalized mid-sentence.
		// Occurs after opening quotes: `` In the morning, ...
		{"She said `` In the morning.", "In", chunky.TagADP},
		{"She said `` By contrast.", "By", chunky.TagADP},
		{"She said `` For example.", "For", chunky.TagADP},
		{"She said `` With that.", "With", chunky.TagADP},

		// neverPropn words: common nouns and abbreviations that stay NOUN even
		// when capitalized mid-sentence (legal terms, acronyms).
		{"filed for Chapter 11 protection.", "Chapter", chunky.TagNOUN},
		{"see Section 3 for details.", "Section", chunky.TagNOUN},
		{"he watches TV every night.", "TV", chunky.TagNOUN},
		{"she bought some CDs yesterday.", "CDs", chunky.TagNOUN},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var found *Token
		for i := range sents {
			for j := range sents[i].Tokens {
				if sents[i].Tokens[j].Word == tc.word {
					found = &sents[i].Tokens[j]
				}
			}
		}
		if found == nil {
			t.Errorf("Parse(%q): token %q not found", tc.input, tc.word)
			continue
		}
		if found.IsUnknownTag() || !found.HasTag(tc.want) {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, found.Tags, tc.want)
		}
	}
}
