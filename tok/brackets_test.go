package tok

import "testing"

func TestStripBrackets(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		// standalone bracket tokens
		{"word [1] more", []string{"word", "more"}},
		{"word [sic] more", []string{"word", "more"}},
		// consecutive bracket tokens
		{"word [1][2] more", []string{"word", "more"}},
		// curly braces
		{"word {x} more", []string{"word", "more"}},
		// multi-token span
		{"the [critical section] and more", []string{"the", "and", "more"}},
		// bracket at end
		{"see here [1]", []string{"see", "here"}},
		// unclosed bracket — pass through
		{"word [unclosed more text", []string{"word", "[unclosed", "more", "text"}},
		// no brackets — passthrough
		{"the dog runs", []string{"the", "dog", "runs"}},
		// inline citation embedded in word: [digits] replaced with spaces
		// The trailing "." becomes a separate space-pad but stays in the word;
		// SplitPunctuation (run separately) handles the final split.
		{"Planeteers.[8]", []string{"Planeteers.   "}},
		{"word.[1][2]", []string{"word.      "}},
	}
	for _, tc := range cases {
		tokens := StripBrackets(Tokenize(tc.input))
		if len(tokens) != len(tc.want) {
			t.Errorf("StripBrackets(%q): got %v, want %v", tc.input, tokWords(tokens), tc.want)
			continue
		}
		for i, tok := range tokens {
			if tok.Word != tc.want[i] {
				t.Errorf("StripBrackets(%q)[%d]: got %q, want %q", tc.input, i, tok.Word, tc.want[i])
			}
		}
	}
}

// TestStripBracketsWithPunctuation verifies the full bracket+punctuation path
// that the plan describes: inline citations are replaced with spaces, then
// SplitPunctuation correctly locates the trailing punctuation.
func TestStripBracketsWithPunctuation(t *testing.T) {
	cases := []struct {
		input string
		want  []string
	}{
		{"Planeteers.[8] The", []string{"Planeteers", ".", "The"}},
		{"word.[1][2] more", []string{"word", ".", "more"}},
		{"see here [1].", []string{"see", "here", "."}},
	}
	for _, tc := range cases {
		tokens := SplitPunctuation(NormalizeText(StripBrackets(Tokenize(tc.input))))
		if len(tokens) != len(tc.want) {
			t.Errorf("pipeline(%q): got %v, want %v", tc.input, tokWords(tokens), tc.want)
			continue
		}
		for i, tok := range tokens {
			if tok.Word != tc.want[i] {
				t.Errorf("pipeline(%q)[%d]: got %q, want %q", tc.input, i, tok.Word, tc.want[i])
			}
		}
	}
}
