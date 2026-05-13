package chunky

import "strings"

// CompoundTags maps space-separated lowercase word sequences to their UD tag.
// Used by both the runtime retokenizer and corpus processing tools (lexrules)
// so that rule generation and tagging operate on the same token stream.
//
// Keys are lowercase words joined by spaces.
var CompoundTags = map[string]Tag{
	// adverbial / discourse
	"such as":     TagADP,
	"as such":     TagADV,
	"as well as":  TagCCONJ,
	"as well":     TagADV,
	"rather than": TagADP,

	// prepositional
	"due to":         TagADP,
	"along with":     TagADP,
	"in addition to": TagADP,
	"in terms of":    TagADP,
	"as a result":    TagADV,
	"in order to":    TagPART,
	"according to":   TagADP,
	"as of":          TagADP,
	"out of":         TagADP,
	"instead of":     TagADP,
	"because of":     TagADP,
	"in spite of":    TagADP,
	"on behalf of":   TagADP,

	// subordinating
	"as long as":  TagSCONJ,
	"as soon as":  TagSCONJ,
	"even though": TagSCONJ,
	"even if":     TagSCONJ,
	"so that":     TagSCONJ,
	"in case":     TagSCONJ,
}

// CompoundMaxLen is the length (in words) of the longest entry in CompoundTags.
// Callers use this to bound their n-gram scan window.
var CompoundMaxLen = func() int {
	max := 0
	for k := range CompoundTags {
		n := 1 + strings.Count(k, " ")
		if n > max {
			max = n
		}
	}
	return max
}()
