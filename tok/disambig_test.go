package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func makeToken(word string, tags ...chunky.Tag) Token {
	return Token{Word: word, Tags: tags}
}

// Single-rule tests: prev tag controls NOUN vs VERB resolution.
var singleSlotRules = []ContextRule{
	// DET before → NOUN
	{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Prev: chunky.TagDET, Mask: maskPrev, Resolve: chunky.TagNOUN},
	// AUX before → VERB
	{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Prev: chunky.TagAUX, Mask: maskPrev, Resolve: chunky.TagVERB},
	// NOUN next → VERB
	{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Next: chunky.TagNOUN, Mask: maskNext, Resolve: chunky.TagVERB},
}

func TestDisambiguateWith_DETprev(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if len(out[1].Tags) != 1 || out[1].Tags[0] != chunky.TagNOUN {
		t.Errorf("DET prev: got %v, want NOUN", out[1].Tags)
	}
}

func TestDisambiguateWith_AUXprev(t *testing.T) {
	tokens := []Token{
		makeToken("will", chunky.TagAUX),
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if len(out[1].Tags) != 1 || out[1].Tags[0] != chunky.TagVERB {
		t.Errorf("AUX prev: got %v, want VERB", out[1].Tags)
	}
}

func TestDisambiguateWith_NoMatch(t *testing.T) {
	// PROPN before → no rule fires, token stays ambiguous.
	tokens := []Token{
		makeToken("London", chunky.TagPROPN),
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if len(out[1].Tags) != 2 {
		t.Errorf("no match: got %v, want unchanged {NOUN,VERB}", out[1].Tags)
	}
}

func TestDisambiguateWith_AlreadyUnambiguous(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("cat", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if len(out[1].Tags) != 1 || out[1].Tags[0] != chunky.TagNOUN {
		t.Errorf("already unambiguous: got %v, want NOUN", out[1].Tags)
	}
	if out[1].Rule != "" {
		t.Errorf("unambiguous token Rule changed: %q", out[1].Rule)
	}
}

// Ambiguous neighbor: rule requires prev to be unambiguous, so it must not fire.
func TestDisambiguateWith_AmbiguousNeighborBlocks(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagNOUN, chunky.TagVERB), // ambiguous neighbor
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	// Even though DET would resolve token[1] to NOUN, token[0] is ambiguous
	// so the rule cannot fire for token[1].
	out := disambiguateWith(tokens, singleSlotRules)
	if len(out[1].Tags) != 2 {
		t.Errorf("ambiguous neighbor: got %v, want unchanged {NOUN,VERB}", out[1].Tags)
	}
}

// Multi-pass cascade: right-to-left dependency needs a second pass.
// tokens: [NOUN/VERB, NOUN/VERB, NOUN]
// rule "next=NOUN → VERB": token[1] resolves in pass 1 (next=NOUN),
// then token[0] resolves in pass 2 (next now = VERB, which is unambiguous).
func TestDisambiguateWith_MultiPassCascade(t *testing.T) {
	cascadeRules := []ContextRule{
		// next=NOUN → VERB
		{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Next: chunky.TagNOUN, Mask: maskNext, Resolve: chunky.TagVERB},
		// next=VERB → VERB (so token[0] resolves once token[1] is VERB)
		{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Next: chunky.TagVERB, Mask: maskNext, Resolve: chunky.TagVERB},
	}
	tokens := []Token{
		makeToken("state", chunky.TagNOUN, chunky.TagVERB), // pass 2: next=VERB → VERB
		makeToken("time", chunky.TagNOUN, chunky.TagVERB),  // pass 1: next=NOUN → VERB
		makeToken("laws", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, cascadeRules)
	if len(out[0].Tags) != 1 || out[0].Tags[0] != chunky.TagVERB {
		t.Errorf("cascade token[0]: got %v, want VERB", out[0].Tags)
	}
	if len(out[1].Tags) != 1 || out[1].Tags[0] != chunky.TagVERB {
		t.Errorf("cascade token[1]: got %v, want VERB", out[1].Tags)
	}
}

// Sentence boundary: TagUNK in an active slot must match only an absent neighbor.
func TestDisambiguateWith_BoundarySlot(t *testing.T) {
	boundaryRules := []ContextRule{
		// sentence-initial (no prev) + next=DET → NOUN
		{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Prev: chunky.TagUNK, Next: chunky.TagDET, Mask: maskPrev | maskNext, Resolve: chunky.TagNOUN},
	}
	// Token at position 0: no prev, next=DET → should fire.
	tokens := []Token{
		makeToken("State", chunky.TagNOUN, chunky.TagVERB),
		makeToken("the", chunky.TagDET),
	}
	out := disambiguateWith(tokens, boundaryRules)
	if len(out[0].Tags) != 1 || out[0].Tags[0] != chunky.TagNOUN {
		t.Errorf("boundary match: got %v, want NOUN", out[0].Tags)
	}

	// Same rule should NOT fire when there is a real prev token.
	tokens2 := []Token{
		makeToken("will", chunky.TagAUX),
		makeToken("state", chunky.TagNOUN, chunky.TagVERB),
		makeToken("the", chunky.TagDET),
	}
	out2 := disambiguateWith(tokens2, boundaryRules)
	if len(out2[1].Tags) != 2 {
		t.Errorf("boundary non-match: got %v, want unchanged {NOUN,VERB}", out2[1].Tags)
	}
}

// Specificity: a 2-slot rule must not override a matching 4-slot rule.
// The 4-slot rule resolves to VERB; the 2-slot rule would resolve to NOUN.
// Since rules are ordered most-specific-first, VERB wins.
func TestDisambiguateWith_SpecificityOrder(t *testing.T) {
	rules := []ContextRule{
		// 4-slot: DET+ADJ+ADP+NOUN → VERB (very specific)
		{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Prev2: chunky.TagDET, Prev: chunky.TagADJ, Next: chunky.TagADP, Next2: chunky.TagNOUN, Mask: 0x0f, Resolve: chunky.TagVERB},
		// 1-slot: prev=ADJ → NOUN (less specific)
		{Tag1: chunky.TagNOUN, Tag2: chunky.TagVERB, Prev: chunky.TagADJ, Mask: maskPrev, Resolve: chunky.TagNOUN},
	}
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("big", chunky.TagADJ),
		makeToken("state", chunky.TagNOUN, chunky.TagVERB),
		makeToken("of", chunky.TagADP),
		makeToken("mind", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, rules)
	if len(out[2].Tags) != 1 || out[2].Tags[0] != chunky.TagVERB {
		t.Errorf("specificity: got %v, want VERB (4-slot rule wins)", out[2].Tags)
	}
}

// RuleField: resolved token should have "+ctx" appended to its Rule.
func TestDisambiguateWith_RuleField(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		{Word: "run", Tags: []chunky.Tag{chunky.TagNOUN, chunky.TagVERB}, Rule: "lexicon"},
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if out[1].Rule != "lexicon+ctx" {
		t.Errorf("Rule field: got %q, want %q", out[1].Rule, "lexicon+ctx")
	}
}
