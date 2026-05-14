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
		{"business\u2014Turner", []string{"business", "\u2014", "Turner"}},
		// en dash splits into three tokens
		{"1990\u20132000", []string{"1990", "\u2013", "2000"}},
		// leading/trailing em dash
		{"\u2014word", []string{"\u2014", "word"}},
		{"word\u2014", []string{"word", "\u2014"}},
		// em dash with surrounding words
		{"million\u2014after", []string{"million", "\u2014", "after"}},
		// multiple em dashes
		{"a\u2014b\u2014c", []string{"a", "\u2014", "b", "\u2014", "c"}},
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
	// "go\u2014now" -- em dash is 3 bytes (U+2014 = 0xE2 0x80 0x94)
	tokens := SplitPunctuation(Tokenize("go\u2014now"))
	want := []struct {
		word   string
		offset int
	}{
		{"go", 0}, {"\u2014", 2}, {"now", 5},
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
	// "hello, world." -- verify offsets of split tokens
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
