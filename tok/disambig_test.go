package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func makeToken(word string, tags ...chunky.Tag) Token {
	var combined chunky.Tag
	for _, t := range tags {
		combined |= t
	}
	return Token{Word: word, Tags: combined}
}

// Single-rule tests: prev tag controls NOUN vs VERB resolution.
var singleSlotRules = []ContextRule{
	// DET before → NOUN
	{Tags: chunky.TagNOUN | chunky.TagVERB, Prev: chunky.TagDET, Mask: maskPrev, Resolve: chunky.TagNOUN},
	// AUX before → VERB
	{Tags: chunky.TagNOUN | chunky.TagVERB, Prev: chunky.TagAUX, Mask: maskPrev, Resolve: chunky.TagVERB},
	// NOUN next → VERB
	{Tags: chunky.TagNOUN | chunky.TagVERB, Next: chunky.TagNOUN, Mask: maskNext, Resolve: chunky.TagVERB},
}

func TestDisambiguateWith_DETprev(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if !out[1].IsResolved() || out[1].Tags != chunky.TagNOUN {
		t.Errorf("DET prev: got %v, want NOUN", out[1].Tags)
	}
}

func TestDisambiguateWith_AUXprev(t *testing.T) {
	tokens := []Token{
		makeToken("will", chunky.TagAUX),
		makeToken("run", chunky.TagNOUN, chunky.TagVERB),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if !out[1].IsResolved() || out[1].Tags != chunky.TagVERB {
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
	if out[1].IsResolved() {
		t.Errorf("no match: got %v, want unchanged NOUN|VERB", out[1].Tags)
	}
}

func TestDisambiguateWith_AlreadyUnambiguous(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("cat", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if !out[1].IsResolved() || out[1].Tags != chunky.TagNOUN {
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
	out := disambiguateWith(tokens, singleSlotRules)
	if out[1].IsResolved() {
		t.Errorf("ambiguous neighbor: got %v, want unchanged NOUN|VERB", out[1].Tags)
	}
}

// Multi-pass cascade: right-to-left dependency needs a second pass.
func TestDisambiguateWith_MultiPassCascade(t *testing.T) {
	cascadeRules := []ContextRule{
		{Tags: chunky.TagNOUN | chunky.TagVERB, Next: chunky.TagNOUN, Mask: maskNext, Resolve: chunky.TagVERB},
		{Tags: chunky.TagNOUN | chunky.TagVERB, Next: chunky.TagVERB, Mask: maskNext, Resolve: chunky.TagVERB},
	}
	tokens := []Token{
		makeToken("state", chunky.TagNOUN, chunky.TagVERB),
		makeToken("time", chunky.TagNOUN, chunky.TagVERB),
		makeToken("laws", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, cascadeRules)
	if !out[0].IsResolved() || out[0].Tags != chunky.TagVERB {
		t.Errorf("cascade token[0]: got %v, want VERB", out[0].Tags)
	}
	if !out[1].IsResolved() || out[1].Tags != chunky.TagVERB {
		t.Errorf("cascade token[1]: got %v, want VERB", out[1].Tags)
	}
}

// Sentence boundary: TagUNK in an active slot must match only an absent neighbor.
func TestDisambiguateWith_BoundarySlot(t *testing.T) {
	boundaryRules := []ContextRule{
		{Tags: chunky.TagNOUN | chunky.TagVERB, Prev: chunky.TagUNK, Next: chunky.TagDET, Mask: maskPrev | maskNext, Resolve: chunky.TagNOUN},
	}
	tokens := []Token{
		makeToken("State", chunky.TagNOUN, chunky.TagVERB),
		makeToken("the", chunky.TagDET),
	}
	out := disambiguateWith(tokens, boundaryRules)
	if !out[0].IsResolved() || out[0].Tags != chunky.TagNOUN {
		t.Errorf("boundary match: got %v, want NOUN", out[0].Tags)
	}

	tokens2 := []Token{
		makeToken("will", chunky.TagAUX),
		makeToken("state", chunky.TagNOUN, chunky.TagVERB),
		makeToken("the", chunky.TagDET),
	}
	out2 := disambiguateWith(tokens2, boundaryRules)
	if out2[1].IsResolved() {
		t.Errorf("boundary non-match: got %v, want unchanged NOUN|VERB", out2[1].Tags)
	}
}

// Specificity: a 2-slot rule must not override a matching 4-slot rule.
func TestDisambiguateWith_SpecificityOrder(t *testing.T) {
	rules := []ContextRule{
		{Tags: chunky.TagNOUN | chunky.TagVERB, Prev2: chunky.TagDET, Prev: chunky.TagADJ, Next: chunky.TagADP, Next2: chunky.TagNOUN, Mask: 0x0f, Resolve: chunky.TagVERB},
		{Tags: chunky.TagNOUN | chunky.TagVERB, Prev: chunky.TagADJ, Mask: maskPrev, Resolve: chunky.TagNOUN},
	}
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		makeToken("big", chunky.TagADJ),
		makeToken("state", chunky.TagNOUN, chunky.TagVERB),
		makeToken("of", chunky.TagADP),
		makeToken("mind", chunky.TagNOUN),
	}
	out := disambiguateWith(tokens, rules)
	if !out[2].IsResolved() || out[2].Tags != chunky.TagVERB {
		t.Errorf("specificity: got %v, want VERB (4-slot rule wins)", out[2].Tags)
	}
}

// RuleField: resolved token should have "+ctx" appended to its Rule.
func TestDisambiguateWith_RuleField(t *testing.T) {
	tokens := []Token{
		makeToken("the", chunky.TagDET),
		{Word: "run", Tags: chunky.TagNOUN | chunky.TagVERB, Rule: "lexicon"},
	}
	out := disambiguateWith(tokens, singleSlotRules)
	if out[1].Rule != "lexicon+ctx" {
		t.Errorf("Rule field: got %q, want %q", out[1].Rule, "lexicon+ctx")
	}
}
