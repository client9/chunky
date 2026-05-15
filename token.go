package chunky

import "strings"

// Token is a single word-like unit produced by the tagging pipeline.
// Offset records the token's byte position in the original source string.
// Tags holds the ordered candidate tag set (most likely first); it is empty
// for untagged tokens and has length 1 when the tag is definitive (e.g.
// corpus-assigned ground-truth or a merged compound).
// Rule identifies which pipeline step or rule assigned the tags.
type Token struct {
	Word   string
	Offset int
	Tags   []Tag
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
	return len(t.Tags) == 0 || t.Tags[0] == TagUNK
}

// HasTag reports whether tag x appears anywhere in the token's tag set.
func (t Token) HasTag(x Tag) bool {
	for _, r := range t.Tags {
		if r == x {
			return true
		}
	}
	return false
}

// String returns the token in "word/TAG" format for a single tag, or
// "word/{TAG1,TAG2,...}" for multiple candidates.
func (t Token) String() string {
	if len(t.Tags) == 1 {
		return t.Word + "/" + t.Tags[0].String()
	}
	parts := make([]string, len(t.Tags))
	for i, s := range t.Tags {
		parts[i] = s.String()
	}
	return t.Word + "/{" + strings.Join(parts, ",") + "}"
}
