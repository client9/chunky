package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestDisambiguateOrdinals(t *testing.T) {
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		// NUM: prenominal ordinal (DET|PRON + ordinal + NOUN)
		{"The first chapter was long.", "first", chunky.TagNUM},
		{"Her second attempt succeeded.", "second", chunky.TagNUM},
		{"His third novel won an award.", "third", chunky.TagNUM},
		{"My first car was old.", "first", chunky.TagNUM},

		// ADV: sentential / sequential adverb (next=VERB)
		{"We must first decide the scope.", "first", chunky.TagADV},
		{"They will second the motion.", "second", chunky.TagADV},

		// ADV: discourse marker at sentence start before comma
		{"First, consider the options.", "First", chunky.TagADV},
		{"Second, we reviewed the data.", "Second", chunky.TagADV},

		// ADV: ordinal after ADP ("first of all", "second of the month")
		{"First of all, thank you.", "First", chunky.TagADV},

		// NUM: prenominal after NOUN ("June first", "a split second")
		{"The event was on June first.", "first", chunky.TagNUM},

		// NUM: prenominal after ADJ ("a close second")
		{"She finished a close second.", "second", chunky.TagNUM},

		// NOUN: standalone after DET — time unit ("a second", "the third")
		{"Wait a second.", "second", chunky.TagNOUN},
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
