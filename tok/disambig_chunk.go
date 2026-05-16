package tok

import "github.com/client9/chunky"

// DisambiguateByChunk refines POS tags using chunk membership as evidence.
// It runs after an initial Chunk() pass and resolves residual ambiguities
// where the chunk position is a reliable signal:
//
//   - NOUN/VERB inside a VP → VERB; inside an NP → NOUN
//   - ADJ/VERB  inside a VP → VERB; inside an NP → ADJ
//   - ADP/PART  inside a VP → PART (infinitival "to")
func DisambiguateByChunk(tokens []Token) []Token {
	for i, tok := range tokens {
		if tok.IsResolved() || tok.IsUnknownTag() {
			continue
		}
		kind := tok.Chunk.Kind
		if kind == chunky.ChunkO {
			continue
		}
		switch {
		case tok.HasTag(chunky.TagNOUN) && tok.HasTag(chunky.TagVERB):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = chunky.TagVERB
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkNP {
				tokens[i].Tags = chunky.TagNOUN
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		case tok.HasTag(chunky.TagADJ) && tok.HasTag(chunky.TagVERB):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = chunky.TagVERB
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkNP {
				tokens[i].Tags = chunky.TagADJ
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		case tok.HasTag(chunky.TagADP) && tok.HasTag(chunky.TagPART):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = chunky.TagPART
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkPP {
				tokens[i].Tags = chunky.TagADP
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		case tok.HasTag(chunky.TagAUX) && tok.HasTag(chunky.TagNOUN):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = chunky.TagAUX
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkNP {
				tokens[i].Tags = chunky.TagNOUN
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		}
	}
	return tokens
}
