package tok

import "testing"

func TestNormalizeText(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"hello", "hello"},
		// curly double quotes → straight ASCII quotes
		{"“hello”", "\"hello\""},
	}
	for _, tc := range cases {
		tokens := NormalizeText(Tokenize(tc.input))
		if len(tokens) != 1 {
			t.Errorf("NormalizeText(%q): got %d tokens, want 1", tc.input, len(tokens))
			continue
		}
		if tokens[0].Word != tc.want {
			t.Errorf("NormalizeText(%q): got %q, want %q", tc.input, tokens[0].Word, tc.want)
		}
	}
}

func TestNormalizeTextDashes(t *testing.T) {
	// em/en dashes pass through unchanged; splitting is SplitPunctuation's job
	cases := []struct {
		input string
		want  string
	}{
		{"word—word", "word—word"},
		{"word–word", "word–word"},
	}
	for _, tc := range cases {
		tokens := NormalizeText(Tokenize(tc.input))
		if len(tokens) != 1 {
			t.Errorf("NormalizeText(%q): got %d tokens, want 1", tc.input, len(tokens))
			continue
		}
		if tokens[0].Word != tc.want {
			t.Errorf("NormalizeText(%q): got %q, want %q", tc.input, tokens[0].Word, tc.want)
		}
	}
}
