package chunky

import (
	"math/bits"
	"strings"
)

// Token is a single word-like unit produced by the tagging pipeline.
// Offset records the token's byte position in the original source string.
// Tags holds the candidate tag set as a bitfield; zero means untagged.
// A single set bit means the tag is fully resolved.
// Rule identifies which pipeline step or rule assigned the tags.
type Token struct {
	Word   string
	Offset int
	Tags   Tag
	Rule   string
	Chunk  ChunkTag
}

// Sentence is an ordered slice of tokens forming a single sentence, with
// the byte offset of the first token in the original source string.
type Sentence struct {
	Tokens []Token
	Offset int
}

// IsUnknownTag reports whether the token has no assigned tags.
func (t Token) IsUnknownTag() bool {
	return t.Tags == 0
}

// IsResolved reports whether the token has exactly one tag assigned.
func (t Token) IsResolved() bool {
	return bits.OnesCount32(uint32(t.Tags)) == 1
}

// HasTag reports whether tag x is in the token's tag set.
func (t Token) HasTag(x Tag) bool {
	return t.Tags&x != 0
}

// String returns the token in "word/TAG" format for a single tag, or
// "word/{TAG1,TAG2,...}" for multiple candidates.
func (t Token) String() string {
	if t.IsResolved() {
		return t.Word + "/" + t.Tags.String()
	}
	var parts []string
	for _, tag := range AllTags {
		if t.Tags&tag != 0 {
			parts = append(parts, tag.String())
		}
	}
	if len(parts) == 0 {
		return t.Word + "/<UNK>"
	}
	return t.Word + "/{" + strings.Join(parts, ",") + "}"
}
