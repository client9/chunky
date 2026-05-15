package tok

import "github.com/client9/chunky"

// Chunk assigns IOB chunk tags to tokens using left-to-right pattern matching
// over the UD tag sequence.
func Chunk(tokens []Token) []Token {
	for i := 0; i < len(tokens); {
		if len(tokens[i].Tags) == 0 {
			i++
			continue
		}
		switch tokens[i].Tags[0] {
		case chunky.TagDET, chunky.TagADJ, chunky.TagNOUN, chunky.TagPROPN, chunky.TagNUM, chunky.TagPRON:
			i = markNP(tokens, i)
		case chunky.TagPART:
			if isInfinitival(tokens, i) {
				// Infinitival "to" before a VERB starts a VP.
				i = markVP(tokens, i)
			} else {
				i++
			}
		case chunky.TagADP:
			if isInfinitival(tokens, i) {
				// "to/ADP" before a VERB is infinitival; start VP.
				i = markVP(tokens, i)
			} else {
				tokens[i].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkPP}
				i++
			}
		case chunky.TagAUX, chunky.TagVERB:
			// Possessive 's is always O — never a chunk head.
			if tokens[i].Word == "'s" {
				i++
			} else {
				i = markVP(tokens, i)
			}
		default:
			i++
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
	if i+1 >= len(tokens) || len(tokens[i+1].Tags) == 0 {
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
	if i == start+1 && len(tokens[start].Tags) > 0 && tokens[start].Tags[0] == chunky.TagDET {
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
	if len(tok.Tags) == 0 {
		return false
	}
	switch tok.Tags[0] {
	case chunky.TagADJ, chunky.TagNOUN, chunky.TagPROPN, chunky.TagNUM:
		return true
	}
	return false
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
// ADV, PART, and infinitival ADP "to" only extend a VP when a verbal token
// follows — a trailing negation particle or bare adverb does not stay inside.
// NOUN/VERB ambiguous tokens are treated as verbal when inside an active VP.
func isVPCont(tokens []Token, i int) bool {
	if i >= len(tokens) || len(tokens[i].Tags) == 0 {
		return false
	}
	tok := tokens[i]
	switch tok.Tags[0] {
	case chunky.TagAUX, chunky.TagVERB:
		return true
	case chunky.TagNOUN, chunky.TagADJ:
		return tok.HasTag(chunky.TagVERB)
	case chunky.TagADV:
		// "n't" / "not" are always O in CoNLL-2000, never inside VP.
		if tok.Word == "n't" || tok.Word == "not" {
			return false
		}
		// Other adverbs continue VP only when a verbal token follows.
		return i+1 < len(tokens) && isVerbal(tokens[i+1])
	case chunky.TagPART:
		// Infinitival "to" (PART) continues VP only when a VERB follows.
		return i+1 < len(tokens) && isVerbal(tokens[i+1])
	case chunky.TagADP:
		// Infinitival "to" (ADP) inside a VP chain: "has helped to prevent".
		return tok.Word == "to" && i+1 < len(tokens) && isVerbal(tokens[i+1])
	}
	return false
}

func isVerbal(tok Token) bool {
	if len(tok.Tags) == 0 {
		return false
	}
	return tok.Tags[0] == chunky.TagVERB || tok.Tags[0] == chunky.TagAUX
}
