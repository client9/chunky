package tok

import (
	"testing"
)

func TestFilterBrackets(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		// footnote markers
		{"word [1] more", []string{"word", "more"}},
		{"word [1][2] more", []string{"word", "more"}},
		// single-token named bracket
		{"word [sic] more", []string{"word", "more"}},
		// multi-token span
		{"the [critical section] and more", []string{"the", "and", "more"}},
		// curly braces
		{"word {x} more", []string{"word", "more"}},
		// no brackets — passthrough
		{"the dog runs", []string{"the", "dog", "runs"}},
		// bracket at end of sentence
		{"see here [1].", []string{"see", "here", "."}},
		// unclosed bracket — pass through rather than eating rest of stream
		{"word [unclosed more text", []string{"word", "[unclosed", "more", "text"}},
		// inline citation embedded in word: tokenizer strips [8] before splitting
		{"Planeteers.[8] The", []string{"Planeteers", ".", "The"}},
		{"word.[1][2] more", []string{"word", ".", "more"}},
	}
	for _, tc := range cases {
		tokens := FilterBrackets(TagString(tc.input))
		if len(tokens) != len(tc.want) {
			t.Errorf("FilterBrackets(%q): got %v, want %v", tc.input, tokenWords(tokens), tc.want)
			continue
		}
		for i, tok := range tokens {
			if tok.Word != tc.want[i] {
				t.Errorf("FilterBrackets(%q)[%d]: got %q, want %q", tc.input, i, tok.Word, tc.want[i])
			}
		}
	}
}

func tokenWords(tokens []Token) []string {
	out := make([]string, len(tokens))
	for i, t := range tokens {
		out[i] = t.Word
	}
	return out
}
