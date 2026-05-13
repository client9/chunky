package tok

import (
	"slices"
	"testing"

	"github.com/client9/chunky"
)

func TestIsOrdinal(t *testing.T) {
	yes := []string{"1st", "2nd", "3rd", "4th", "10th", "20th", "21st", "100th"}
	no := []string{"", "st", "nd", "1", "10", "1stx", "th"}
	for _, s := range yes {
		if !isOrdinal(s) {
			t.Errorf("isOrdinal(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if isOrdinal(s) {
			t.Errorf("isOrdinal(%q) = true, want false", s)
		}
	}
}

func TestIsDecade(t *testing.T) {
	yes := []string{"1980s", "1960s", "2000s", "90s"}
	no := []string{"", "s", "1980", "198xs"}
	for _, s := range yes {
		if !isDecade(s) {
			t.Errorf("isDecade(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if isDecade(s) {
			t.Errorf("isDecade(%q) = true, want false", s)
		}
	}
}

func TestIsNumber(t *testing.T) {
	yes := []string{"0", "42", "3.14", "1,000", "1,000.50", "+1", "-1", "-3.14"}
	no := []string{"", ".", ",", "1.", "1,", "+", "-", "abc", "1a"}
	for _, s := range yes {
		if !isNumber(s) {
			t.Errorf("isNumber(%q) = false, want true", s)
		}
	}
	for _, s := range no {
		if isNumber(s) {
			t.Errorf("isNumber(%q) = true, want false", s)
		}
	}
}

func TestSurfaceTokenizeOffsets(t *testing.T) {
	cases := []struct {
		input string
		want  []rawToken
	}{
		{"hello world", []rawToken{{"hello", 0}, {"world", 6}}},
		{"hello, world.", []rawToken{{"hello", 0}, {",", 5}, {"world", 7}, {".", 12}}},
		{"(hello)", []rawToken{{"(", 0}, {"hello", 1}, {")", 6}}},
		{"(hello), world", []rawToken{{"(", 0}, {"hello", 1}, {")", 6}, {",", 7}, {"world", 9}}},
		// non-breaking space (U+00A0, 2 bytes in UTF-8) treated as whitespace
		{"hello world", []rawToken{{"hello", 0}, {"world", 7}}},
	}
	for _, tc := range cases {
		got := surfaceTokenizeRaw(tc.input)
		if len(got) != len(tc.want) {
			t.Errorf("surfaceTokenizeRaw(%q): got %v, want %v", tc.input, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("surfaceTokenizeRaw(%q)[%d]: got %+v, want %+v", tc.input, i, got[i], tc.want[i])
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
		if !slices.Equal(tc.want, got) {
			t.Errorf("SurfaceTokenize(%q) = %v, want %v", tc.input, got, tc.want)
		}
	}
}

func hasTag(tags []chunky.Tag, want chunky.Tag) bool {
	for _, tg := range tags {
		if tg == want {
			return true
		}
	}
	return false
}

func TestMorphCandidates(t *testing.T) {
	tests := []struct {
		word    string
		isFirst bool
		wantTag chunky.Tag
		wantNil bool
	}{
		{"1st", false, chunky.TagADJ, false},
		{"20th", false, chunky.TagADJ, false},
		{"1980s", false, chunky.TagNOUN, false},
		{"42", false, chunky.TagNUM, false},
		{"3.14", false, chunky.TagNUM, false},
		{"London", false, chunky.TagPROPN, false},
		// caps at sentence start — suppressed, no caps rule fires
		{"London", true, chunky.TagUNK, true},
		{"quickly", false, chunky.TagADV, false},
		{"transportation", false, chunky.TagNOUN, false},
		{"mechanisms", false, chunky.TagNOUN, false},
		{"organize", false, chunky.TagVERB, false},
		{"dangerous", false, chunky.TagADJ, false},
		{"running", false, chunky.TagVERB, false},
		{"walked", false, chunky.TagVERB, false},
		{"Japanese-American", false, chunky.TagADJ, false},
		// prefix rules
		{"rewrite", false, chunky.TagVERB, false},
		{"overload", false, chunky.TagVERB, false},
		{"undermine", false, chunky.TagVERB, false},
		{"unnatural", false, chunky.TagADJ, false},
		{"nonprofit", false, chunky.TagADJ, false},
		{"antiwar", false, chunky.TagADJ, false},
		{"prehistoric", false, chunky.TagADJ, false},
	}
	for _, tc := range tests {
		tags, rule := MorphCandidates(tc.word, tc.isFirst)
		if tc.wantNil {
			if tags != nil {
				t.Errorf("MorphCandidates(%q, %v) = %v, want nil", tc.word, tc.isFirst, tags)
			}
			continue
		}
		if tags == nil {
			t.Errorf("MorphCandidates(%q, %v) = nil (rule=%q), want tags including %v", tc.word, tc.isFirst, rule, tc.wantTag)
			continue
		}
		if !hasTag(tags, tc.wantTag) {
			t.Errorf("MorphCandidates(%q, %v) = %v (rule=%q), want %v in result", tc.word, tc.isFirst, tags, rule, tc.wantTag)
		}
	}
}

func TestInflectionCandidates(t *testing.T) {
	tests := []struct {
		word    string
		wantTag chunky.Tag
	}{
		// contractions and possessives are now split at the surface tokenizer;
		// InflectionCandidates no longer handles them.
		{"cats", chunky.TagNOUN},
		{"flies", chunky.TagVERB},
		{"walking", chunky.TagVERB},
		{"making", chunky.TagVERB},
		{"walked", chunky.TagVERB},
		{"faster", chunky.TagADJ},
		{"accelerates", chunky.TagVERB},
		{"running", chunky.TagVERB},
	}
	for _, tc := range tests {
		tags, rule := InflectionCandidates(tc.word)
		if tags == nil {
			t.Errorf("InflectionCandidates(%q) = nil, want tags including %v", tc.word, tc.wantTag)
			continue
		}
		if !hasTag(tags, tc.wantTag) {
			t.Errorf("InflectionCandidates(%q) = %v (rule=%q), want %v in result", tc.word, tags, rule, tc.wantTag)
		}
	}
}

func TestHyphenCandidates(t *testing.T) {
	tests := []struct {
		word    string
		wantTag chunky.Tag
	}{
		{"co-chairman", chunky.TagNOUN},
		{"mid-1990s", chunky.TagNOUN},
		// adjectival suffixes always produce ADJ regardless of standalone tag
		// adjectival suffixes always produce ADJ regardless of standalone tag
		{"flu-like", chunky.TagADJ},
		{"war-like", chunky.TagADJ},
		{"sugar-free", chunky.TagADJ},
		{"industry-wide", chunky.TagADJ},
	}
	for _, tc := range tests {
		tags, rule := HyphenCandidates(tc.word)
		if tags == nil {
			t.Errorf("HyphenCandidates(%q) = nil, want tags including %v", tc.word, tc.wantTag)
			continue
		}
		if !hasTag(tags, tc.wantTag) {
			t.Errorf("HyphenCandidates(%q) = %v (rule=%q), want %v in result", tc.word, tags, rule, tc.wantTag)
		}
	}
}

func TestTagString(t *testing.T) {
	tokens := TagString("the cat")
	if len(tokens) != 2 {
		t.Fatalf("TagString returned %d tokens, want 2", len(tokens))
	}
	if tokens[0].Rule != "lexicon" {
		t.Errorf("tokens[0].Rule = %q, want 'lexicon'", tokens[0].Rule)
	}
	if !hasTag(tokens[0].Canidates, chunky.TagDET) {
		t.Errorf("tokens[0] ('the') candidates = %v, want DET", tokens[0].Canidates)
	}
}

func TestTagUnknowns(t *testing.T) {
	tests := []struct {
		sentence string
		word     string
		wantTag  chunky.Tag
	}{
		{"run quickly", "quickly", chunky.TagADV},
		{"I see cats", "cats", chunky.TagNOUN},
		{"the 20th century", "20th", chunky.TagADJ},
		{"worth 3.14 dollars", "3.14", chunky.TagNUM},
		// inflection path
		{"she accelerates quickly", "accelerates", chunky.TagVERB},
		// hyphen path
		{"a co-chairman spoke", "co-chairman", chunky.TagNOUN},
		// morph prefix path (no lexicon or inflection match)
		{"the reforestation effort", "reforestation", chunky.TagNOUN},
		// Unk1: DET _ NOUN context
		{"the xyzzy thing", "xyzzy", chunky.TagNOUN},
	}
	for _, tc := range tests {
		tokens := TagUnknowns(TagString(tc.sentence))
		found := false
		for _, tok := range tokens {
			if tok.Word != tc.word {
				continue
			}
			found = true
			if !hasTag(tok.Canidates, tc.wantTag) {
				t.Errorf("word %q in %q: candidates = %v (rule=%q), want %v", tc.word, tc.sentence, tok.Canidates, tok.Rule, tc.wantTag)
			}
			break
		}
		if !found {
			t.Errorf("word %q not found in tokenization of %q", tc.word, tc.sentence)
		}
	}
}

func TestContractionSplit(t *testing.T) {
	cases := []struct {
		input  string
		words  []string // expected token words in order
		tags   []chunky.Tag
	}{
		// auxiliary contractions
		{"They're fast.", []string{"They", "'re", "fast", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagADJ, chunky.TagPUNCT}},
		{"They'll go.", []string{"They", "'ll", "go", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPUNCT}},
		{"I've seen it.", []string{"I", "'ve", "seen", "it", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPRON, chunky.TagPUNCT}},
		{"I'm here.", []string{"I", "'m", "here", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagADV, chunky.TagPUNCT}},
		{"He'd gone.", []string{"He", "'d", "gone", "."}, []chunky.Tag{chunky.TagPRON, chunky.TagAUX, chunky.TagVERB, chunky.TagPUNCT}},
		// negating contractions: n belongs to suffix when stem-without-n is in lexicon
		{"can't stop.", []string{"can", "'t", "stop", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"don't go.", []string{"do", "n't", "go", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		{"shouldn't leave.", []string{"should", "n't", "leave", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		// irregular: won't → will + n't
		{"won't go.", []string{"will", "n't", "go", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADV, chunky.TagVERB, chunky.TagPUNCT}},
		// ain't stays whole
		{"ain't right.", []string{"ain't", "right", "."}, []chunky.Tag{chunky.TagAUX, chunky.TagADJ, chunky.TagPUNCT}},
		// possessive 's
		{"John's book.", []string{"John", "'s", "book", "."}, []chunky.Tag{0, chunky.TagAUX, chunky.TagNOUN, chunky.TagPUNCT}},
	}

	for _, tc := range cases {
		tokens := TagUnknowns(FilterBrackets(TagString(tc.input)))
		if len(tokens) != len(tc.words) {
			t.Errorf("TagString(%q): got %d tokens %v, want %d %v", tc.input, len(tokens), tokWords(tokens), len(tc.words), tc.words)
			continue
		}
		for i, w := range tc.words {
			if tokens[i].Word != w {
				t.Errorf("TagString(%q)[%d]: word = %q, want %q", tc.input, i, tokens[i].Word, w)
			}
			if tc.tags[i] != 0 && !hasTag(tokens[i].Canidates, tc.tags[i]) {
				t.Errorf("TagString(%q)[%d] %q: candidates = %v, want %v", tc.input, i, w, tokens[i].Canidates, tc.tags[i])
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
