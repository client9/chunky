package tok

import "github.com/client9/chunky"

// DisambiguateByChunk refines POS tags using chunk membership as evidence.
// It runs after an initial Chunk() pass and resolves residual ambiguities
// where the chunk position is a reliable signal:
//
//   - NOUN/VERB inside a VP chunk → VERB
//   - NOUN/VERB inside an NP chunk → NOUN
//   - ADJ/VERB  inside a VP chunk → VERB
//   - ADJ/VERB  inside an NP chunk → ADJ
func DisambiguateByChunk(tokens []Token) []Token {
	for i, tok := range tokens {
		if len(tok.Tags) <= 1 {
			continue
		}
		kind := tok.Chunk.Kind
		if kind == chunky.ChunkO {
			continue
		}
		switch {
		case tok.HasTag(chunky.TagNOUN) && tok.HasTag(chunky.TagVERB):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = []chunky.Tag{chunky.TagVERB}
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkNP {
				tokens[i].Tags = []chunky.Tag{chunky.TagNOUN}
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		case tok.HasTag(chunky.TagADJ) && tok.HasTag(chunky.TagVERB):
			if kind == chunky.ChunkVP {
				tokens[i].Tags = []chunky.Tag{chunky.TagVERB}
				tokens[i].Rule = tok.Rule + "+chunk"
			} else if kind == chunky.ChunkNP {
				tokens[i].Tags = []chunky.Tag{chunky.TagADJ}
				tokens[i].Rule = tok.Rule + "+chunk"
			}
		}
	}
	return tokens
}
