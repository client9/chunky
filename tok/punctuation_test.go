package tok

import "testing"

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
		// em dash splits into three tokens
		{"business—Turner", []string{"business", "—", "Turner"}},
		// en dash splits into three tokens
		{"1990–2000", []string{"1990", "–", "2000"}},
		// leading/trailing em dash
		{"—word", []string{"—", "word"}},
		{"word—", []string{"word", "—"}},
		// em dash with surrounding punctuation handled correctly
		{"million—after", []string{"million", "—", "after"}},
		// multiple dashes
		{"a—b—c", []string{"a", "—", "b", "—", "c"}},
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

func TestSplitPunctuationEmDashOffsets(t *testing.T) {
	// "go—now" — em dash is 3 bytes (U+2014 = 0xE2 0x80 0x94)
	tokens := SplitPunctuation(Tokenize("go—now"))
	want := []struct {
		word   string
		offset int
	}{
		{"go", 0}, {"—", 2}, {"now", 5},
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
