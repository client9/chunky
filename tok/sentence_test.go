package tok

import (
	"strings"
	"testing"
)

func TestSegment(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"The dog runs.", []string{"The dog runs ."}},
		{"The dog runs. The cat sleeps.", []string{"The dog runs .", "The cat sleeps ."}},
		{"Stop! Are you sure?", []string{"Stop !", "Are you sure ?"}},
		{"Dr. Smith arrived.", []string{"Dr. Smith arrived ."}},
		{"Ted C. Turner spoke.", []string{"Ted C . Turner spoke ."}},
		{"The dog runs", []string{"The dog runs"}},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		if len(sents) != len(tc.want) {
			t.Errorf("Parse(%q): got %d sentences, want %d", tc.input, len(sents), len(tc.want))
			for i, s := range sents {
				t.Logf("  [%d] %s", i, joinWords(s.Tokens))
			}
			continue
		}
		for i, s := range sents {
			got := joinWords(s.Tokens)
			if got != tc.want[i] {
				t.Errorf("Parse(%q)[%d]: got %q, want %q", tc.input, i, got, tc.want[i])
			}
		}
	}
}

func TestSegmentOffset(t *testing.T) {
	sents := Parse("Hello. World.")
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
	sents := Parse("The dog runs. The cat sleeps.")
	if len(sents) != 2 {
		t.Fatalf("want 2 sentences, got %d", len(sents))
	}
	first := sents[1].Tokens[0]
	if first.Word != "The" {
		t.Fatalf("expected 'The', got %q", first.Word)
	}
	if len(first.Candidates) != 1 || first.Candidates[0].String() != "DET" {
		t.Errorf("'The' at sentence start: got %v, want [DET]", first.Candidates)
	}
}

func TestSegmentSentenceInitialNoun(t *testing.T) {
	sents := Parse("Hantaviruses are dangerous.")
	if len(sents) != 1 {
		t.Fatalf("want 1 sentence, got %d", len(sents))
	}
	first := sents[0].Tokens[0]
	if first.Word != "Hantaviruses" {
		t.Fatalf("expected 'Hantaviruses', got %q", first.Word)
	}
	if len(first.Candidates) != 1 || first.Candidates[0].String() != "NOUN" {
		t.Errorf("'Hantaviruses' at sentence start: got %v, want [NOUN]", first.Candidates)
	}
}

func TestSegmentSentenceInitialLexiconWord(t *testing.T) {
	sents := Parse("Run fast.")
	if len(sents) != 1 {
		t.Fatalf("want 1 sentence, got %d", len(sents))
	}
	first := sents[0].Tokens[0]
	if first.Word != "Run" {
		t.Fatalf("expected 'Run', got %q", first.Word)
	}
	for _, c := range first.Candidates {
		if c.String() == "PROPN" {
			t.Errorf("'Run' at sentence start: got PROPN, want lexicon tag (VERB/NOUN)")
		}
	}
}

func joinWords(tokens []Token) string {
	var b strings.Builder
	for i, t := range tokens {
		if i > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(t.Word)
	}
	return b.String()
}
