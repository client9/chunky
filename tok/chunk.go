package tok

import (
	"strings"

	"github.com/client9/chunky"
)

// Chunk assigns IOB chunk tags to tokens using left-to-right pattern matching
// over the UD tag sequence.
func Chunk(tokens []Token) []Token {
	for i := 0; i < len(tokens); {
		if tokens[i].IsUnknownTag() {
			i++
			continue
		}
		switch tokens[i].Tags {
		case chunky.TagADJ:
			if isPredicateADJ(tokens, i) {
				i = markADJP(tokens, i)
			} else {
				i = markNP(tokens, i)
			}
		case chunky.TagDET, chunky.TagNOUN, chunky.TagPROPN, chunky.TagNUM, chunky.TagPRON:
			i = markNP(tokens, i)
		case chunky.TagPART:
			if isInfinitival(tokens, i) {
				i = markVP(tokens, i)
			} else if tokens[i].Word == "'s" {
				// Possessive 's starts the possessed NP in CoNLL-2000 convention:
				// "the company 's board" → NP(the company) + NP('s board).
				i = markNP(tokens, i)
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
			if tokens[i].Word == "'s" && !isAuxVP(tokens, i) {
				i++
			} else {
				i = markVP(tokens, i)
			}
		default:
			// Handle common ambiguous combinations.
			switch {
			case tokens[i].HasTag(chunky.TagAUX):
				// AUX|NOUN (will, may) → start VP
				if tokens[i].Word == "'s" && !isAuxVP(tokens, i) {
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
			case tokens[i].HasTag(chunky.TagNOUN) || tokens[i].HasTag(chunky.TagPROPN) || tokens[i].HasTag(chunky.TagDET) || tokens[i].HasTag(chunky.TagPRON):
				// NOUN|VERB, DET|PRON, PRON|X, etc. — try NP
				i = markNP(tokens, i)
			default:
				i++
			}
		}
	}
	return tokens
}

// isAuxVP reports whether "'s" at position i should start a VP — i.e., it is
// the contracted "is/has" auxiliary followed by a verbal or adjectival complement.
func isAuxVP(tokens []Token, i int) bool {
	next := tokenAt(tokens, i+1)
	return next.HasTag(chunky.TagVERB) || next.HasTag(chunky.TagAUX) || next.HasTag(chunky.TagADJ) || next.HasTag(chunky.TagADV)
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
	return tokens[i+1].HasTag(chunky.TagVERB) || tokens[i+1].HasTag(chunky.TagAUX)
}

// copulaVerbs are linking verbs that can take a predicate adjective complement.
var copulaVerbs = map[string]bool{
	"remain":   true,
	"remains":  true,
	"remained": true,
	"seem":     true,
	"seems":    true,
	"seemed":   true,
	"appear":   true,
	"appears":  true,
	"appeared": true,
	"become":   true,
	"becomes":  true,
	"became":   true,
	"prove":    true,
	"proves":   true,
	"proved":   true,
	"proven":   true,
	"feel":     true,
	"feels":    true,
	"felt":     true,
	"look":     true,
	"looks":    true,
	"looked":   true,
	"sound":    true,
	"sounds":   true,
	"sounded":  true,
	"taste":    true,
	"tastes":   true,
	"tasted":   true,
	"smell":    true,
	"smells":   true,
	"smelled":  true,
	"stay":     true,
	"stays":    true,
	"stayed":   true,
	"stand":    true,
	"stands":   true,
	"stood":    true,
	"keep":     true,
	"keeps":    true,
	"kept":     true,
	"turn":     true,
	"turns":    true,
	"turned":   true,
	"grow":     true,
	"grows":    true,
	"grew":     true,
	"grown":    true,
	"get":      true,
	"gets":     true,
	"got":      true,
}

// isPredicateADJ reports whether the pure ADJ at i is a predicate adjective —
// i.e., it follows a resolved AUX or copula VERB. "is unchanged", "remained likely", etc.
func isPredicateADJ(tokens []Token, i int) bool {
	prev := tokenAt(tokens, i-1)
	if prev.Tags == chunky.TagAUX {
		return true
	}
	if prev.Tags == chunky.TagVERB && copulaVerbs[prev.Word] {
		return true
	}
	return false
}

func markADJP(tokens []Token, start int) int {
	i := start + 1
	for i < len(tokens) && tokens[i].Tags == chunky.TagADJ {
		i++
	}
	tokens[start].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkADJP}
	for j := start + 1; j < i; j++ {
		tokens[j].Chunk = chunky.ChunkTag{IOB: 'I', Kind: chunky.ChunkADJP}
	}
	return i
}

// whPronouns are WH-words that head single-token NPs in CoNLL-2000 convention.
// They introduce relative clauses or questions and never extend into I-NP.
var whPronouns = map[string]bool{
	"who": true, "whom": true, "whose": true,
	"which": true, "what": true,
	"whoever": true, "whichever": true, "whatever": true,
}

func markNP(tokens []Token, start int) int {
	// WH-pronouns head single-token NPs: "who left" → B-NP(who) B-VP(left).
	if tokens[start].Tags == chunky.TagPRON && whPronouns[strings.ToLower(tokens[start].Word)] {
		tokens[start].Chunk = chunky.ChunkTag{IOB: 'B', Kind: chunky.ChunkNP}
		return start + 1
	}
	i := start + 1
	for i < len(tokens) && isNPCont(tokens[i]) {
		i++
	}
	// Suppress NPs that consist only of a bare DET — "the", "a", "an" alone
	// are never chunk heads in CoNLL-2000 style. Exception: partitive quantifiers
	// ("most of", "more of", "much of") head NPs with a following PP complement.
	if i == start+1 && tokens[start].Tags == chunky.TagDET {
		next := tokenAt(tokens, i)
		if next.Word != "of" {
			return i
		}
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
	return tok.HasTag(chunky.TagVERB) || tok.HasTag(chunky.TagAUX)
}
