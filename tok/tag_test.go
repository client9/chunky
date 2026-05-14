package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestLexicalTag(t *testing.T) {
	tokens := LexicalTag(SplitPunctuation(NormalizeText(StripBrackets(Tokenize("the cat")))))
	if len(tokens) != 2 {
		t.Fatalf("LexicalTag returned %d tokens, want 2", len(tokens))
	}
	if tokens[0].Rule != "lexicon" {
		t.Errorf("tokens[0].Rule = %q, want 'lexicon'", tokens[0].Rule)
	}
	if !hasTag(tokens[0].Candidates, chunky.TagDET) {
		t.Errorf("tokens[0] ('the') candidates = %v, want DET", tokens[0].Candidates)
	}
}

func TestLexicalTagPreservesCompounds(t *testing.T) {
	// MergeLexical sets candidates on compound tokens; LexicalTag must not overwrite them.
	tokens := LexicalTag(MergeLexical(SplitPunctuation(NormalizeText(StripBrackets(Tokenize("such as"))))))
	if len(tokens) != 1 {
		t.Fatalf("got %d tokens, want 1 (compound): %v", len(tokens), tokWords(tokens))
	}
	if tokens[0].Rule != "compound" {
		t.Errorf("compound token Rule = %q, want 'compound'", tokens[0].Rule)
	}
	if !hasTag(tokens[0].Candidates, chunky.TagADP) {
		t.Errorf("'such as' candidates = %v, want ADP", tokens[0].Candidates)
	}
}
