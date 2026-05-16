package tok

import (
	"testing"

	"github.com/client9/chunky"
)

// chunkOf returns the ChunkTag of the first token matching word in sents.
func chunkOf(sents []Sentence, word string) chunky.ChunkTag {
	for _, s := range sents {
		for _, tok := range s.Tokens {
			if tok.Word == word {
				return tok.Chunk
			}
		}
	}
	return chunky.ChunkTag{}
}

// chunkKindOf returns the chunk kind (NP/VP/PP/O) of the first matching token.
func chunkKindOf(sents []Sentence, word string) chunky.ChunkKind {
	return chunkOf(sents, word).Kind
}

func TestIsInfinitival_ToVerb(t *testing.T) {
	// Original: to + VERB should produce B-VP
	cases := []struct {
		input string
		want  chunky.ChunkKind
	}{
		{"They plan to sell shares.", chunky.ChunkVP},
		{"She wants to run the company.", chunky.ChunkVP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, "to")
		if got != tc.want {
			t.Errorf("Parse(%q) \"to\": chunk=%v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestIsInfinitival_ToAux(t *testing.T) {
	// to + AUX (be, have, do) should produce B-VP, not B-PP
	cases := []struct {
		input string
		want  chunky.ChunkKind
	}{
		{"It needs to be done.", chunky.ChunkVP},
		{"They ought to have left.", chunky.ChunkVP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, "to")
		if got != tc.want {
			t.Errorf("Parse(%q) \"to\": chunk=%v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestNPSpan_PossessiveS(t *testing.T) {
	// CoNLL-2000 convention: 's starts the POSSESSED NP, not the possessor NP.
	// "the company 's board" → NP(the company) + NP('s board).
	cases := []struct {
		input string
		word  string
		iob   byte
	}{
		{"the company 's board", "the", 'B'},
		{"the company 's board", "company", 'I'},
		{"the company 's board", "'s", 'B'}, // starts new NP
		{"the company 's board", "board", 'I'},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkOf(sents, tc.word)
		if got.IOB != tc.iob || (tc.iob != 'O' && got.Kind != chunky.ChunkNP) {
			t.Errorf("Parse(%q) %q: chunk=%v, want %c-NP", tc.input, tc.word, got, tc.iob)
		}
	}
}

func TestChunkDefault_AmbiguousPRON(t *testing.T) {
	// Personal pronouns should always produce a NP chunk.
	cases := []string{"he", "she", "it", "they", "we", "I"}
	templates := []string{
		"%s said nothing.",
		"Then %s left.",
	}
	for _, pron := range cases {
		for _, tmpl := range templates {
			input := tmpl
			_ = input // use pron in sentence
			sents := Parse(pron + " said nothing.")
			got := chunkKindOf(sents, pron)
			if got != chunky.ChunkNP {
				t.Errorf("Parse(%q) %q: chunk=%v, want NP", pron+" said nothing.", pron, got)
			}
		}
	}
}

func TestChunkDefault_NounVerbPrefersVP(t *testing.T) {
	// NOUN|VERB words following a subject NP should be VP, not NP.
	cases := []struct {
		input string
		word  string
		want  chunky.ChunkKind
	}{
		{"The company says it will expand.", "says", chunky.ChunkVP},
		{"The board plans to sell assets.", "plans", chunky.ChunkVP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, tc.word)
		if got != tc.want {
			t.Errorf("Parse(%q) %q: chunk=%v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}

func TestIsVerbal_AmbiguousVerb(t *testing.T) {
	// isVerbal must use HasTag so NOUN|VERB tokens extend VP spans.
	// "declined to comment" should be one VP (or at least to+comment extend the VP).
	cases := []struct {
		input string
		word  string
		want  chunky.ChunkKind
	}{
		// "to" must continue the VP when next word is NOUN|VERB
		{"The company declined to comment.", "to", chunky.ChunkVP},
		{"She refused to answer.", "to", chunky.ChunkVP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, tc.word)
		if got != tc.want {
			t.Errorf("Parse(%q) %q: chunk=%v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}

func TestPredicateADJ_ADJP(t *testing.T) {
	// Pure ADJ following a resolved AUX (copula) should be B-ADJP, not B-NP.
	cases := []struct {
		input string
		word  string
		want  chunky.ChunkKind
	}{
		{"The result is unchanged.", "unchanged", chunky.ChunkADJP},
		{"The plan was unlikely.", "unlikely", chunky.ChunkADJP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, tc.word)
		if got != tc.want {
			t.Errorf("Parse(%q) %q: chunk=%v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}

func TestIsAuxVP_ContractedAux(t *testing.T) {
	// "'s" (contracted is/has) followed by AUX should be B-VP, not B-PP.
	cases := []struct {
		input string
		want  chunky.ChunkKind
	}{
		{"It 's been done.", chunky.ChunkVP},
		{"He 's been working.", chunky.ChunkVP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, "'s")
		if got != tc.want {
			t.Errorf("Parse(%q) \"'s\": chunk=%v, want %v", tc.input, got, tc.want)
		}
	}
}

func TestPredicateADJ_CopulaVerb(t *testing.T) {
	// Pure ADJ (single-tag) following a copula VERB should be B-ADJP, not B-NP.
	// Words must be unambiguously ADJ in the lexicon for isPredicateADJ to fire.
	cases := []struct {
		input string
		word  string
		want  chunky.ChunkKind
	}{
		{"The situation remained unchanged.", "unchanged", chunky.ChunkADJP},
		{"The outcome seemed reasonable.", "reasonable", chunky.ChunkADJP},
		{"The policy proved ineffective.", "ineffective", chunky.ChunkADJP},
		{"The plan appears unlikely.", "unlikely", chunky.ChunkADJP},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkKindOf(sents, tc.word)
		if got != tc.want {
			t.Errorf("Parse(%q) %q: chunk=%v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}

func TestVPSpan_NegationNt(t *testing.T) {
	// VP should extend through "n't" when followed by a verbal head.
	// "did n't know" should be a single VP span.
	cases := []struct {
		input string
		word  string
		iob   byte
	}{
		{"She did n't know the answer.", "n't", 'I'},
		{"They could n't find it.", "n't", 'I'},
		{"He would n't leave.", "n't", 'I'},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkOf(sents, tc.word)
		if got.IOB != tc.iob || got.Kind != chunky.ChunkVP {
			t.Errorf("Parse(%q) %q: chunk=%v, want %c-VP", tc.input, tc.word, got, tc.iob)
		}
	}
}

func TestWHPronoun_SingleTokenNP(t *testing.T) {
	// WH-pronouns should be single-token B-NP chunks (no I-NP extension).
	cases := []struct {
		input string
		word  string
	}{
		{"The man who arrived was late.", "who"},
		{"The policy which failed was costly.", "which"},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got := chunkOf(sents, tc.word)
		if got.IOB != 'B' || got.Kind != chunky.ChunkNP {
			t.Errorf("Parse(%q) %q: chunk=%v, want B-NP", tc.input, tc.word, got)
		}
		// Verify next token is NOT I-NP (WH-pronoun doesn't extend into relative clause)
		for _, s := range sents {
			for j, tok := range s.Tokens {
				if tok.Word == tc.word && j+1 < len(s.Tokens) {
					next := s.Tokens[j+1]
					if next.Chunk.IOB == 'I' && next.Chunk.Kind == chunky.ChunkNP {
						t.Errorf("Parse(%q) %q: next token %q is I-NP, expected new chunk", tc.input, tc.word, next.Word)
					}
				}
			}
		}
	}
}

func TestGerundNouns_AsNoun(t *testing.T) {
	// Gerunds used as noun heads or modifiers should resolve to NOUN.
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"The trading volume rose.", "trading", chunky.TagNOUN},
		{"The operating profit declined.", "operating", chunky.TagNOUN},
		{"The banking sector grew.", "banking", chunky.TagNOUN},
		{"The funding was approved.", "funding", chunky.TagNOUN},
		{"Global warming is a problem.", "warming", chunky.TagNOUN},
		{"The developing countries need aid.", "developing", chunky.TagNOUN},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}

func TestGerundNouns_AsVerb(t *testing.T) {
	// Gerunds inside a VP (AUX + gerund) should resolve to VERB.
	cases := []struct {
		input string
		word  string
		want  chunky.Tag
	}{
		{"He is trading stocks.", "trading", chunky.TagVERB},
		{"The company is operating normally.", "operating", chunky.TagVERB},
		{"They are developing a plan.", "developing", chunky.TagVERB},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		got, resolved := tagOf(sents, tc.word)
		if !resolved {
			t.Errorf("Parse(%q) %q: still ambiguous, want %v", tc.input, tc.word, tc.want)
			continue
		}
		if got != tc.want {
			t.Errorf("Parse(%q) %q: got %v, want %v", tc.input, tc.word, got, tc.want)
		}
	}
}
