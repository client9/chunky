package tok

import (
	"github.com/client9/chunky"
)

// Mask bits for ContextRule.Mask — which of the four neighboring positions are checked.
const (
	maskPrev2 uint8 = 1 << 3
	maskPrev  uint8 = 1 << 2
	maskNext  uint8 = 1 << 1
	maskNext2 uint8 = 1 << 0
)

// ContextRule fires when a token's tag set contains all bits in Tags
// and the active context slots match. Mask selects which slots are checked;
// unchecked slots are wildcards. TagUNK in an active slot matches an absent
// neighbor (sentence boundary).
type ContextRule struct {
	Tags        chunky.Tag
	Prev2, Prev chunky.Tag
	Next, Next2 chunky.Tag
	Mask        uint8
	Resolve     chunky.Tag
}

// tokenAt returns the token at index i, or a zero Token if i is out of bounds.
func tokenAt(tokens []Token, i int) Token {
	if i < 0 || i >= len(tokens) {
		return Token{}
	}
	return tokens[i]
}

// resolvedAs reports whether tok is fully resolved to exactly tag.
func resolvedAs(tok Token, tag Tag) bool {
	return tok.IsResolved() && tok.Tags == tag
}

func matchSlot(want chunky.Tag, tok Token, active bool) bool {
	if !active {
		return true
	}
	if want == chunky.TagUNK {
		// Active TagUNK means "must be sentence boundary/absent."
		return tok.Word == ""
	}
	return tok.IsResolved() && tok.Tags&want != 0
}

// CopyTags returns a snapshot of each token's Tags for change detection.
func CopyTags(tokens []Token) []chunky.Tag {
	out := make([]chunky.Tag, len(tokens))
	for i, t := range tokens {
		out[i] = t.Tags
	}
	return out
}

// TagsEqual reports whether tokens have the same Tags as the snapshot produced by CopyTags.
func TagsEqual(tokens []Token, snap []chunky.Tag) bool {
	if len(tokens) != len(snap) {
		return false
	}
	for i, t := range tokens {
		if t.Tags != snap[i] {
			return false
		}
	}
	return true
}

// applyRules makes one left-to-right pass, firing the first matching rule for
// each ambiguous token. Returns true if any token changed.
func applyRules(tokens []Token, rules []ContextRule) bool {
	changed := false
	for i, tok := range tokens {
		if tok.IsResolved() || tok.IsUnknownTag() {
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
			if tok.Tags&r.Tags != r.Tags {
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
			tokens[i].Tags = r.Resolve
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

// contextRules is the combined rule table: globally-sorted generated rules
// (all pairs, most-specific-first) followed by 1-slot broad fallbacks.
var contextRules = append(
	append(append(append(append(append(
		generatedRules,
		detPronBroadRules...),
		advDetBroadRules...),
		adpSconjBroadRules...),
		nounVerbBroadRules...),
		adjNounBroadRules...),
)

// DisambiguateContext resolves ambiguous tokens using the compiled context rule
// table. It runs to fixed point so that tokens resolved in one pass can unblock
// neighbors in the next.
func DisambiguateContext(tokens []Token) []Token {
	return disambiguateWith(tokens, contextRules)
}
