package tok

import (
	"testing"

	"github.com/client9/chunky"
)

func TestTokenize(t *testing.T) {
	cases := []struct {
		input string
		want  []Token
	}{
		{"hello world", []Token{{Word: "hello", Offset: 0}, {Word: "world", Offset: 6}}},
		// non-breaking space (U+00A0, 2 bytes in UTF-8) treated as whitespace
		{"hello world", []Token{{Word: "hello", Offset: 0}, {Word: "world", Offset: 7}}},
		// no splitting of punctuation — dumb tokenizer only splits on whitespace
		{"hello,", []Token{{Word: "hello,", Offset: 0}}},
		{"", nil},
		{"one", []Token{{Word: "one", Offset: 0}}},
	}
	for _, tc := range cases {
		got := Tokenize(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("Tokenize(%q): got %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i].Word != tc.want[i].Word || got[i].Offset != tc.want[i].Offset {
				t.Errorf("Tokenize(%q)[%d]: got {%q, %d}, want {%q, %d}",
					tc.input, i, got[i].Word, got[i].Offset, tc.want[i].Word, tc.want[i].Offset)
			}
		}
	}
}

func TestSurfaceTokenize(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"The dog runs.", []string{"The", "dog", "runs", "."}},
		{"hello world", []string{"hello", "world"}},
		{"hello, world.", []string{"hello", ",", "world", "."}},
		{"(hello)", []string{"(", "hello", ")"}},
		{"(hello), world", []string{"(", "hello", ")", ",", "world"}},
		{"hello: world;", []string{"hello", ":", "world", ";"}},
		{"world!", []string{"world", "!"}},
		{"really?", []string{"really", "?"}},
		{"", []string{}},
		{"one", []string{"one"}},
	}
	for _, tc := range cases {
		got := SurfaceTokenize(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("SurfaceTokenize(%q) = %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("SurfaceTokenize(%q)[%d] = %q, want %q", tc.input, i, got[i], tc.want[i])
			}
		}
	}
}

func tokWords(tokens []Token) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = t.Word
	}
	return out
}

func hasTag(tags chunky.Tag, want chunky.Tag) bool {
	return tags&want != 0
}
