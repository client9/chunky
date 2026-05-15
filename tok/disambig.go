package tok

import (
	"sort"

	"github.com/client9/chunky"
)

// Mask bits for ContextRule.Mask — which of the four neighboring positions are checked.
const (
	maskPrev2 uint8 = 1 << 3
	maskPrev  uint8 = 1 << 2
	maskNext  uint8 = 1 << 1
	maskNext2 uint8 = 1 << 0
)

// ContextRule fires when a token's tag set contains both Tag1 and Tag2
// and the active context slots match. Mask selects which slots are checked;
// unchecked slots are wildcards. TagUNK in an active slot matches an absent
// neighbor (sentence boundary).
type ContextRule struct {
	Tag1, Tag2      chunky.Tag
	Prev2, Prev     chunky.Tag
	Next, Next2     chunky.Tag
	Mask            uint8
	Resolve         chunky.Tag
}

func (r ContextRule) specificity() int {
	n := 0
	for _, b := range []uint8{maskPrev2, maskPrev, maskNext, maskNext2} {
		if r.Mask&b != 0 {
			n++
		}
	}
	return n
}

func matchSlot(want chunky.Tag, tok Token, active bool) bool {
	if !active {
		return true
	}
	if want == chunky.TagUNK {
		// Active TagUNK means "must be sentence boundary/absent."
		return tok.Word == ""
	}
	return len(tok.Tags) == 1 && tok.Tags[0] == want
}

// applyRules makes one left-to-right pass, firing the first matching rule for
// each ambiguous token. Returns true if any token changed.
func applyRules(tokens []Token, rules []ContextRule) bool {
	changed := false
	for i, tok := range tokens {
		if len(tok.Tags) <= 1 {
			continue
		}
		var prev2, prev, next, next2 Token
		if i >= 2 {
			prev2 = tokens[i-2]
		}
		if i >= 1 {
			prev = tokens[i-1]
		}
		if i+1 < len(tokens) {
			next = tokens[i+1]
		}
		if i+2 < len(tokens) {
			next2 = tokens[i+2]
		}
		for _, r := range rules {
			if !tok.HasTag(r.Tag1) || !tok.HasTag(r.Tag2) {
				continue
			}
			if !matchSlot(r.Prev2, prev2, r.Mask&maskPrev2 != 0) {
				continue
			}
			if !matchSlot(r.Prev, prev, r.Mask&maskPrev != 0) {
				continue
			}
			if !matchSlot(r.Next, next, r.Mask&maskNext != 0) {
				continue
			}
			if !matchSlot(r.Next2, next2, r.Mask&maskNext2 != 0) {
				continue
			}
			tokens[i].Tags = []chunky.Tag{r.Resolve}
			tokens[i].Rule = tokens[i].Rule + "+ctx"
			changed = true
			break
		}
	}
	return changed
}

// disambiguateWith resolves ambiguous tokens using the given rules, iterating
// until fixed point. Rules must be sorted most-specific-first.
func disambiguateWith(tokens []Token, rules []ContextRule) []Token {
	for applyRules(tokens, rules) {
	}
	return tokens
}

// contextRules is the combined, specificity-sorted rule table built from all
// generated per-pair rule sets. Rules from different pairs don't compete for
// 2-way ambiguous tokens (Tag1/Tag2 are checked first), but 3-way ambiguous
// tokens (e.g. ADJ/NOUN/VERB) can match rules from multiple pairs, so the
// global sort ensures the most-specific rule always fires first.
var contextRules []ContextRule

func init() {
	contextRules = make([]ContextRule, 0,
		len(nounVerbRules)+len(adjNounRules)+len(adpPartRules)+len(auxVerbRules)+
			len(detPronRules)+len(adpSconjRules)+len(adjVerbRules)+len(adjAdvRules)+
			len(advDetRules)+len(adpAdvRules)+len(advNumRules))
	contextRules = append(contextRules, nounVerbRules...)
	contextRules = append(contextRules, adjNounRules...)
	contextRules = append(contextRules, adpPartRules...)
	contextRules = append(contextRules, auxVerbRules...)
	contextRules = append(contextRules, detPronRules...)
	contextRules = append(contextRules, adpSconjRules...)
	contextRules = append(contextRules, adjVerbRules...)
	contextRules = append(contextRules, adjAdvRules...)
	contextRules = append(contextRules, advDetRules...)
	contextRules = append(contextRules, adpAdvRules...)
	contextRules = append(contextRules, advNumRules...)
	sort.SliceStable(contextRules, func(i, j int) bool {
		return contextRules[i].specificity() > contextRules[j].specificity()
	})
}

// DisambiguateContext resolves ambiguous tokens using the compiled context rule
// table. It runs to fixed point so that tokens resolved in one pass can unblock
// neighbors in the next.
func DisambiguateContext(tokens []Token) []Token {
	return disambiguateWith(tokens, contextRules)
}
