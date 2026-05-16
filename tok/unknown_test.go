package tok

import (
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

func TestNumericCandidates(t *testing.T) {
	tests := []struct {
		word    string
		wantTag chunky.Tag
	}{
		{"1st", chunky.TagADJ},
		{"20th", chunky.TagADJ},
		{"1980s", chunky.TagNOUN},
		{"42", chunky.TagNUM},
		{"3.14", chunky.TagNUM},
		{"15%", chunky.TagNUM},
		// currency
		{"$1", chunky.TagNUM},
		{"$1,000", chunky.TagNUM},
		{"$3.50", chunky.TagNUM},
		{"£5", chunky.TagNUM},
		{"€1,000", chunky.TagNUM},
		{"¥500", chunky.TagNUM},
		// fractions
		{"3/8", chunky.TagNUM},
		{"1/2", chunky.TagNUM},
		{"3/4", chunky.TagNUM},
	}
	for _, tc := range tests {
		tags, rule := NumericCandidates(tc.word)
		if tags == 0 {
			t.Errorf("NumericCandidates(%q) = nil, want tags including %v", tc.word, tc.wantTag)
			continue
		}
		if !hasTag(tags, tc.wantTag) {
			t.Errorf("NumericCandidates(%q) = %v (rule=%q), want %v in result", tc.word, tags, rule, tc.wantTag)
		}
	}
}

func TestMorphCandidates(t *testing.T) {
	tests := []struct {
		word    string
		isFirst bool
		wantTag chunky.Tag
		wantNil bool
	}{
		{"London", false, chunky.TagPROPN, false},
		{"London", true, chunky.TagUNK, true},
		{"quickly", false, chunky.TagADV, false},
		{"transportation", false, chunky.TagNOUN, false},
		{"mechanisms", false, chunky.TagNOUN, false},
		{"organize", false, chunky.TagVERB, false},
		{"dangerous", false, chunky.TagADJ, false},
		{"running", false, chunky.TagVERB, false},
		{"walked", false, chunky.TagVERB, false},
		{"Japanese-American", false, chunky.TagADJ, false},
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
			if tags != 0 {
				t.Errorf("MorphCandidates(%q, %v) = %v, want nil", tc.word, tc.isFirst, tags)
			}
			continue
		}
		if tags == 0 {
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
		if tags == 0 {
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
		{"flu-like", chunky.TagADJ},
		{"war-like", chunky.TagADJ},
		{"sugar-free", chunky.TagADJ},
		{"industry-wide", chunky.TagADJ},
	}
	for _, tc := range tests {
		tags, rule := HyphenCandidates(tc.word)
		if tags == 0 {
			t.Errorf("HyphenCandidates(%q) = nil, want tags including %v", tc.word, tc.wantTag)
			continue
		}
		if !hasTag(tags, tc.wantTag) {
			t.Errorf("HyphenCandidates(%q) = %v (rule=%q), want %v in result", tc.word, tags, rule, tc.wantTag)
		}
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
		{"she accelerates quickly", "accelerates", chunky.TagVERB},
		{"a co-chairman spoke", "co-chairman", chunky.TagNOUN},
		{"the reforestation effort", "reforestation", chunky.TagNOUN},
		{"the xyzzy thing", "xyzzy", chunky.TagNOUN},
	}
	for _, tc := range tests {
		sents := Parse(tc.sentence)
		var tokens []Token
		for _, s := range sents {
			tokens = append(tokens, s.Tokens...)
		}
		found := false
		for _, tok := range tokens {
			if tok.Word != tc.word {
				continue
			}
			found = true
			if !hasTag(tok.Tags, tc.wantTag) {
				t.Errorf("word %q in %q: candidates = %v (rule=%q), want %v", tc.word, tc.sentence, tok.Tags, tok.Rule, tc.wantTag)
			}
			break
		}
		if !found {
			t.Errorf("word %q not found in Parse(%q)", tc.word, tc.sentence)
		}
	}
}
