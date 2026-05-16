package tok

import "github.com/client9/chunky"

// Chunk assigns IOB chunk tags to tokens using left-to-right pattern matching
// over the UD tag sequence.
func Chunk(tokens []Token) []Token {
	for i := 0; i < len(tokens); {
		if tokens[i].IsUnknownTag() {
			i++
			continue
		}
		switch tokens[i].Tags {
		case chunky.TagDET, chunky.TagADJ, chunky.TagNOUN, chunky.TagPROPN, chunky.TagNUM, chunky.TagPRON:
			i = markNP(tokens, i)
		case chunky.TagPART:
			if isInfinitival(tokens, i) {
				i = markVP(tokens, i)
			} else {
				i++
			}
		case chunky.TagADP:
			if isInfinitival(tokens, i) {
				i = markVP(tokens, i)
			} else {
				tokens[i].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkPP}
				i++
			}
		case chunky.TagAUX, chunky.TagVERB:
			if tokens[i].Word == "'s" {
				i++
			} else {
				i = markVP(tokens, i)
			}
		default:
			// Handle common ambiguous combinations.
			switch {
			case tokens[i].HasTag(chunky.TagAUX):
				// AUX|NOUN (will, may) → start VP
				if tokens[i].Word == "'s" {
					i++
				} else {
					i = markVP(tokens, i)
				}
			case tokens[i].HasTag(chunky.TagADP) || tokens[i].HasTag(chunky.TagPART):
				// ADP|PART (to) — check infinitival first
				if isInfinitival(tokens, i) {
					i = markVP(tokens, i)
				} else if tokens[i].HasTag(chunky.TagADP) {
					tokens[i].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkPP}
					i++
				} else {
					i++
				}
			case tokens[i].HasTag(chunky.TagNOUN) || tokens[i].HasTag(chunky.TagPROPN) || tokens[i].HasTag(chunky.TagDET):
				// NOUN|VERB, DET|PRON, etc. — try NP
				i = markNP(tokens, i)
			default:
				i++
			}
		}
	}
	return tokens
}

// isInfinitival reports whether the token at i is "to" (ADP or PART) used as
// an infinitive marker — i.e., followed by a VERB.
func isInfinitival(tokens []Token, i int) bool {
	tok := tokens[i]
	if tok.Word != "to" {
		return false
	}
	if i+1 >= len(tokens) || tokens[i+1].IsUnknownTag() {
		return false
	}
	return tokens[i+1].HasTag(chunky.TagVERB)
}

func markNP(tokens []Token, start int) int {
	i := start + 1
	for i < len(tokens) && isNPCont(tokens[i]) {
		i++
	}
	// Suppress NPs that consist only of a bare DET — "the", "a", "an" alone
	// are never chunk heads in CoNLL-2000 style.
	if i == start+1 && tokens[start].Tags == chunky.TagDET {
		return i
	}
	tokens[start].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkNP}
	for j := start + 1; j < i; j++ {
		tokens[j].Chunk = chunky.ChunkTag{IOB: 'I', Kind: chunky.ChunkNP}
	}
	return i
}

// isNPCont returns true if tok can continue (I-NP) a noun phrase in progress.
func isNPCont(tok Token) bool {
	if tok.IsUnknownTag() {
		return false
	}
	switch tok.Tags {
	case chunky.TagADJ, chunky.TagNOUN, chunky.TagPROPN, chunky.TagNUM:
		return true
	}
	// Ambiguous NOUN|VERB or PROPN-having tokens can continue NPs;
	// DisambiguateByChunk will resolve them using chunk context.
	return tok.HasTag(chunky.TagNOUN) || tok.HasTag(chunky.TagPROPN)
}

func markVP(tokens []Token, start int) int {
	tokens[start].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkVP}
	i := start + 1
	for i < len(tokens) && isVPCont(tokens, i) {
		tokens[i].Chunk = chunky.ChunkTag{IOB: 'I', Kind: chunky.ChunkVP}
		i++
	}
	return i
}

// isVPCont returns true if the token at i can continue (I-VP) a verb phrase.
func isVPCont(tokens []Token, i int) bool {
	if i >= len(tokens) || tokens[i].IsUnknownTag() {
		return false
	}
	tok := tokens[i]
	if tok.Tags == chunky.TagAUX || tok.Tags == chunky.TagVERB {
		return true
	}
	// Ambiguous NOUN/VERB or ADJ/VERB — treat as verbal inside VP.
	if tok.HasTag(chunky.TagVERB) && (tok.HasTag(chunky.TagNOUN) || tok.HasTag(chunky.TagADJ)) {
		return true
	}
	if tok.Tags == chunky.TagADV {
		if tok.Word == "n't" || tok.Word == "not" {
			return false
		}
		return i+1 < len(tokens) && isVerbal(tokens[i+1])
	}
	if tok.Tags == chunky.TagPART {
		return i+1 < len(tokens) && isVerbal(tokens[i+1])
	}
	if tok.Tags == chunky.TagADP {
		return tok.Word == "to" && i+1 < len(tokens) && isVerbal(tokens[i+1])
	}
	return false
}

func isVerbal(tok Token) bool {
	return tok.Tags == chunky.TagVERB || tok.Tags == chunky.TagAUX
}
