package tok

import "testing"

func TestSegment(t *testing.T) {
	cases := []struct {
		input string
		want  []string // expected sentences reconstructed as space-joined words
	}{
		// basic single sentence
		{
			"The dog runs.",
			[]string{"The dog runs ."},
		},
		// two sentences
		{
			"The dog runs. The cat sleeps.",
			[]string{"The dog runs .", "The cat sleeps ."},
		},
		// exclamation and question marks
		{
			"Stop! Are you sure?",
			[]string{"Stop !", "Are you sure ?"},
		},
		// abbreviation does not split
		{
			"Dr. Smith arrived.",
			[]string{"Dr . Smith arrived ."},
		},
		// middle initial does not split
		{
			"Ted C. Turner spoke.",
			[]string{"Ted C . Turner spoke ."},
		},
		// no terminal punctuation
		{
			"The dog runs",
			[]string{"The dog runs"},
		},
		// sentence-initial word not promoted to PROPN
		{
			"The dog runs. The cat sleeps.",
			[]string{"The dog runs .", "The cat sleeps ."},
		},
	}

	for _, tc := range cases {
		sents := Segment(TagUnknowns(FilterBrackets(TagString(tc.input))))
		if len(sents) != len(tc.want) {
			t.Errorf("Segment(%q): got %d sentences, want %d", tc.input, len(sents), len(tc.want))
			for i, s := range sents {
				t.Logf("  [%d] %s", i, joinWords(s.Tokens))
			}
			continue
		}
		for i, s := range sents {
			got := joinWords(s.Tokens)
			if got != tc.want[i] {
				t.Errorf("Segment(%q)[%d]: got %q, want %q", tc.input, i, got, tc.want[i])
			}
		}
	}
}

func TestSegmentOffset(t *testing.T) {
	sents := Segment(TagUnknowns(FilterBrackets(TagString("Hello. World."))))
	if len(sents) != 2 {
		t.Fatalf("want 2 sentences, got %d", len(sents))
	}
	if sents[0].Offset != 0 {
		t.Errorf("sentence 0 offset: got %d, want 0", sents[0].Offset)
	}
	if sents[1].Offset != 7 {
		t.Errorf("sentence 1 offset: got %d, want 7", sents[1].Offset)
	}
}

func TestSegmentSentenceInitialDET(t *testing.T) {
	sents := Segment(TagUnknowns(FilterBrackets(TagString("The dog runs. The cat sleeps."))))
	if len(sents) != 2 {
		t.Fatalf("want 2 sentences, got %d", len(sents))
	}
	// second sentence starts with "The" which should be DET not PROPN
	first := sents[1].Tokens[0]
	if first.Word != "The" {
		t.Fatalf("expected 'The', got %q", first.Word)
	}
	if len(first.Canidates) != 1 || first.Canidates[0].String() != "DET" {
		t.Errorf("'The' at sentence start: got %v, want [DET]", first.Canidates)
	}
}

// TestSegmentSentenceInitialNoun verifies that a sentence-initial unknown word
// tagged NOUN by the unk:word fallback stays NOUN and is not promoted to PROPN.
// Capitalization at sentence start is grammatical, not a proper-noun signal.
func TestSegmentSentenceInitialNoun(t *testing.T) {
	sents := Segment(TagUnknowns(FilterBrackets(TagString("Hantaviruses are dangerous."))))
	if len(sents) != 1 {
		t.Fatalf("want 1 sentence, got %d", len(sents))
	}
	first := sents[0].Tokens[0]
	if first.Word != "Hantaviruses" {
		t.Fatalf("expected 'Hantaviruses', got %q", first.Word)
	}
	if len(first.Canidates) != 1 || first.Canidates[0].String() != "NOUN" {
		t.Errorf("'Hantaviruses' at sentence start: got %v, want [NOUN]", first.Canidates)
	}
}

// TestSegmentSentenceInitialLexiconWord verifies that a sentence-initial word
// found in the lexicon keeps its lexicon tag rather than being promoted to PROPN.
func TestSegmentSentenceInitialLexiconWord(t *testing.T) {
	sents := Segment(TagUnknowns(FilterBrackets(TagString("Run fast."))))
	if len(sents) != 1 {
		t.Fatalf("want 1 sentence, got %d", len(sents))
	}
	first := sents[0].Tokens[0]
	if first.Word != "Run" {
		t.Fatalf("expected 'Run', got %q", first.Word)
	}
	// "run" is in the lexicon; should NOT be promoted to PROPN
	for _, c := range first.Canidates {
		if c.String() == "PROPN" {
			t.Errorf("'Run' at sentence start: got PROPN, want lexicon tag (VERB/NOUN)")
		}
	}
}

func joinWords(tokens []Token) string {
	out := ""
	for i, t := range tokens {
		if i > 0 {
			out += " "
		}
		out += t.Word
	}
	return out
}
