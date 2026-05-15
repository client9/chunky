package chunky

import "fmt"

// ChunkKind identifies the phrase type of a chunk.
type ChunkKind int

const (
	ChunkO    ChunkKind = iota // outside any chunk
	ChunkNP                    // noun phrase
	ChunkVP                    // verb phrase
	ChunkPP                    // prepositional phrase
	ChunkADVP                  // adverb phrase
	ChunkADJP                  // adjective phrase
)

func (k ChunkKind) String() string {
	switch k {
	case ChunkO:
		return "O"
	case ChunkNP:
		return "NP"
	case ChunkVP:
		return "VP"
	case ChunkPP:
		return "PP"
	case ChunkADVP:
		return "ADVP"
	case ChunkADJP:
		return "ADJP"
	default:
		panic(fmt.Sprintf("unknown ChunkKind %d", int(k)))
	}
}

// ParseChunkKind converts a string like "NP" or "VP" to a ChunkKind.
// Returns ChunkO and an error if unrecognized.
func ParseChunkKind(s string) (ChunkKind, error) {
	switch s {
	case "O":
		return ChunkO, nil
	case "NP":
		return ChunkNP, nil
	case "VP":
		return ChunkVP, nil
	case "PP":
		return ChunkPP, nil
	case "ADVP":
		return ChunkADVP, nil
	case "ADJP":
		return ChunkADJP, nil
	default:
		return ChunkO, fmt.Errorf("chunk: unknown kind %q", s)
	}
}

// ChunkTag is an IOB-encoded chunk label attached to a token.
// The zero value is outside (O).
type ChunkTag struct {
	IOB  byte      // 'B' = begin, 'I' = inside, 'O' = outside
	Kind ChunkKind // phrase type; meaningful only when IOB != 'O'
}

// String returns the CoNLL-style label: "O", "B-NP", "I-VP", etc.
func (c ChunkTag) String() string {
	if c.IOB == 'O' || c.IOB == 0 {
		return "O"
	}
	return string(c.IOB) + "-" + c.Kind.String()
}

// ParseChunkTag parses a CoNLL IOB label like "B-NP", "I-VP", or "O".
func ParseChunkTag(s string) (ChunkTag, error) {
	if s == "O" {
		return ChunkTag{IOB: 'O'}, nil
	}
	if len(s) < 3 || s[1] != '-' {
		return ChunkTag{}, fmt.Errorf("chunk: invalid tag %q", s)
	}
	iob := s[0]
	if iob != 'B' && iob != 'I' {
		return ChunkTag{}, fmt.Errorf("chunk: invalid IOB byte %q in %q", iob, s)
	}
	kind, err := ParseChunkKind(s[2:])
	if err != nil {
		return ChunkTag{}, err
	}
	return ChunkTag{IOB: iob, Kind: kind}, nil
}

// Span is a contiguous sequence of tokens forming a single phrase.
type Span struct {
	Kind   ChunkKind
	Tokens []Token
}

// ToSpans groups a token sequence into phrase spans, discarding outside tokens.
// Tokens with ChunkTag zero value (IOB==0) are treated as outside.
func ToSpans(tokens []Token) []Span {
	var spans []Span
	for i := 0; i < len(tokens); {
		tok := tokens[i]
		if tok.Chunk.IOB != 'B' {
			i++
			continue
		}
		kind := tok.Chunk.Kind
		j := i + 1
		for j < len(tokens) && tokens[j].Chunk.IOB == 'I' && tokens[j].Chunk.Kind == kind {
			j++
		}
		spans = append(spans, Span{Kind: kind, Tokens: tokens[i:j]})
		i = j
	}
	return spans
}
