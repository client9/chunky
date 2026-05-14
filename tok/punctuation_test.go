package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestSplitPunctuation(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"hello world", []string{"hello", "world"}},
		{"hello, world.", []string{"hello", ",", "world", "."}},
		{"(hello)", []string{"(", "hello", ")"}},
		{"(hello), world", []string{"(", "hello", ")", ",", "world"}},
		{"hello: world;", []string{"hello", ":", "world", ";"}},
		{"world!", []string{"world", "!"}},
		{"really?", []string{"really", "?"}},
		// dotted abbreviation: dot stays attached
		{"Dr. Smith", []string{"Dr.", "Smith"}},
		{"etc. more", []string{"etc.", "more"}},
		// internal spaces from bracket replacement: punct found via last-non-space scan
		{"word   .", []string{"word", "."}},
		// trailing spaces with dot: dot is split off (not an abbreviation)
		{"Planeteers.   ", []string{"Planeteers", "."}},
	}
	for _, tc := range cases {
		tokens := SplitPunctuation(Tokenize(tc.input))
		if len(tokens) != len(tc.want) {
			t.Errorf("SplitPunctuation(%q): got %v, want %v", tc.input, tokWords(tokens), tc.want)
			continue
		}
		for i, tok := range tokens {
			if tok.Word != tc.want[i] {
				t.Errorf("SplitPunctuation(%q)[%d]: got %q, want %q", tc.input, i, tok.Word, tc.want[i])
			}
		}
	}
}

func TestSplitPunctuationOffsets(t *testing.T) {
	// "hello, world." — verify offsets of split tokens
	tokens := SplitPunctuation(Tokenize("hello, world."))
	want := []struct {
		word   string
		offset int
	}{
		{"hello", 0}, {",", 5}, {"world", 7}, {".", 12},
	}
	if len(tokens) != len(want) {
		t.Fatalf("got %d tokens, want %d: %v", len(tokens), len(want), tokWords(tokens))
	}
	for i, w := range want {
		if tokens[i].Word != w.word || tokens[i].Offset != w.offset {
			t.Errorf("[%d]: got {%q, %d}, want {%q, %d}", i, tokens[i].Word, tokens[i].Offset, w.word, w.offset)
		}
	}
}

func TestContractionSplit(t *testing.T) {
	cases := []struct {
		input string
		words []string
		tags  []chunky.Tag
	}{
		{"They're fast.", []string{"They", "'re", "fast", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagADJ, chunky.TagPUNCT}},
		{"They'll go.", []string{"They", "'ll", "go", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPUNCT}},
		{"I've seen it.", []string{"I", "'ve", "seen", "it", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPRON, chunky.TagPUNCT}},
		{"I'm here.", []string{"I", "'m", "here", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagADV, chunky.TagPUNCT}},
		{"He'd gone.", []string{"He", "'d", "gone", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPUNCT}},
		{"can't stop.", []string{"can", "'t", "stop", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"don't go.", []string{"do", "n't", "go", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"shouldn't leave.", []string{"should", "n't", "leave", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"won't go.", []string{"will", "n't", "go", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"ain't right.", []string{"ain't", "right", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADJ, chunky.TagPUNCT}},
		{"John's book.", []string{"John", "'s", "book", "."}, []chunky.Tag{0, chunky.TagAUX, chunky.TagNOUN, chunky.TagPUNCT}},
	}
	for _, tc := range cases {
		sents := Parse(tc.input)
		var tokens []Token
		for _, s := range sents {
			tokens = append(tokens, s.Tokens...)
		}
		if len(tokens) != len(tc.words) {
			t.Errorf("Parse(%q): got %d tokens %v, want %d %v", tc.input, len(tokens), tokWords(tokens), len(tc.words), tc.words)
			continue
		}
		for i, w := range tc.words {
			if tokens[i].Word != w {
				t.Errorf("Parse(%q)[%d]: word = %q, want %q", tc.input, i, tokens[i].Word, w)
			}
			if tc.tags[i] != 0 && !hasTag(tokens[i].Candidates, tc.tags[i]) {
				t.Errorf("Parse(%q)[%d] %q: candidates = %v, want %v", tc.input, i, w, tokens[i].Candidates, tc.tags[i])
			}
		}
	}
}
