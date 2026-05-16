package tok

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/client9/chunky"
)

// TagUnknowns fills in candidate tags for tokens that have none, trying rules
// in order: numeric → inflection → hyphen → morphology → alpha fallback.
func TagUnknowns(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.IsUnknownTag() {
			continue
		}
		if candidates, rule := NumericCandidates(t.Word); candidates != 0 {
			tokens[i].Tags = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := InflectionCandidates(t.Word); candidates != 0 {
			tokens[i].Tags = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := HyphenCandidates(t.Word); candidates != 0 {
			tokens[i].Tags = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := MorphCandidates(t.Word, i == 0); candidates != 0 {
			tokens[i].Tags = candidates
			tokens[i].Rule = rule
			continue
		}
		if isAlpha(t.Word) {
			tokens[i].Tags = chunky.TagNOUN
			tokens[i].Rule = "unk:word"
		}
	}
	return tokens
}

// closedClassMask covers tags that inflection never produces.
const closedClassMask = chunky.TagPRON | chunky.TagDET | chunky.TagADP | chunky.TagAUX |
	chunky.TagSCONJ | chunky.TagCCONJ | chunky.TagPART | chunky.TagPUNCT |
	chunky.TagINTJ | chunky.TagSYM

// InflectionCandidates looks up candidate tags by stripping common inflectional
// suffixes and checking the resulting stem against the lexicon. Returns
// candidates and a rule ID of the form "inflect:<suffix>:<stem>".
func InflectionCandidates(word string) (chunky.Tag, string) {
	lower := strings.ToLower(word)

	var tags chunky.Tag
	matchedStem := ""
	matchedSuffix := ""

	try := func(suffix, stem string) {
		t := wordtagmap[stem]
		if t == 0 {
			return
		}
		tags |= t
		if matchedStem == "" {
			matchedStem = stem
			matchedSuffix = suffix
		}
	}
	tryDoubled := func(suffix, stem string) {
		if len(stem) >= 2 && stem[len(stem)-1] == stem[len(stem)-2] {
			try(suffix, stem[:len(stem)-1])
		}
	}

	if strings.HasSuffix(lower, "ies") && len(lower) > 3 {
		try("-ies", lower[:len(lower)-3]+"y")
	}
	if strings.HasSuffix(lower, "s") && len(lower) > 2 {
		try("-s", lower[:len(lower)-1])
	}
	if strings.HasSuffix(lower, "es") && len(lower) > 3 {
		try("-es", lower[:len(lower)-2])
	}
	if strings.HasSuffix(lower, "ing") && len(lower) > 4 {
		stem := lower[:len(lower)-3]
		try("-ing", stem)
		try("-ing+e", stem+"e")
		tryDoubled("-ing+double", stem)
	}
	if strings.HasSuffix(lower, "ed") && len(lower) > 3 {
		stem := lower[:len(lower)-2]
		try("-ed", stem)
		try("-ed+e", stem+"e")
		tryDoubled("-ed+double", stem)
	}
	if strings.HasSuffix(lower, "er") && len(lower) > 3 {
		stem := lower[:len(lower)-2]
		try("-er", stem)
		try("-er+e", stem+"e")
		tryDoubled("-er+double", stem)
	}
	if strings.HasSuffix(lower, "est") && len(lower) > 4 {
		stem := lower[:len(lower)-3]
		try("-est", stem)
		try("-est+e", stem+"e")
	}

	if tags == 0 {
		return 0, ""
	}
	// Strip closed-class tags that crept in via false stem matches.
	tags &^= closedClassMask
	if tags == 0 {
		return 0, ""
	}
	return tags, "inflect:" + matchedSuffix + ":" + matchedStem
}

// hyphenAdjSuffixes always produce ADJ in a hyphenated compound.
var hyphenAdjSuffixes = map[string]bool{
	"like": true,
	"free": true,
	"wide": true,
}

// isCompoundAdj reports whether the last component of a hyphenated word
// signals a compound adjective.
func isCompoundAdj(last string) bool {
	if len(last) > 3 && strings.HasSuffix(last, "ed") {
		return true
	}
	if len(last) > 4 && strings.HasSuffix(last, "ing") {
		return true
	}
	return false
}

// HyphenCandidates handles hyphenated compounds by looking up the last
// component in the lexicon and applying morph rules as a fallback.
func HyphenCandidates(word string) (chunky.Tag, string) {
	i := strings.LastIndex(word, "-")
	if i < 0 || i == len(word)-1 {
		return 0, ""
	}
	last := strings.ToLower(word[i+1:])

	if hyphenAdjSuffixes[last] {
		return chunky.TagADJ, "hyphen:adj-suffix:" + last
	}
	if isCompoundAdj(last) {
		return chunky.TagADJ, "hyphen:compound-participle"
	}
	if tags := wordtagmap[last]; tags != 0 {
		return tags, "hyphen:lexicon:" + last
	}
	if tags, rule := NumericCandidates(last); tags != 0 {
		return tags, "hyphen:" + rule
	}
	if tags, rule := InflectionCandidates(last); tags != 0 {
		return tags, "hyphen:" + rule
	}
	if tags, rule := MorphCandidates(last, false); tags != 0 {
		return tags, "hyphen:" + rule
	}
	return 0, ""
}

// currencySymbols is the set of leading characters that mark a currency amount.
const currencySymbols = "$£€¥¢₹₩₪"

// NumericCandidates tags numeric forms: integers, decimals, ordinals, decades,
// percentages, and currency amounts ($1, £5.50, €1,000). Returns 0 if not numeric.
func NumericCandidates(word string) (chunky.Tag, string) {
	lower := strings.ToLower(word)

	// Strip a leading currency symbol and treat the remainder as a number.
	if r, size := utf8.DecodeRuneInString(lower); r != utf8.RuneError && strings.ContainsRune(currencySymbols, r) && len(lower) > size {
		if isNumber(lower[size:]) {
			return chunky.TagNUM, "morph:currency"
		}
	}

	if isOrdinal(lower) {
		return chunky.TagADJ, "morph:ordinal"
	}
	if isDecade(lower) {
		return chunky.TagNOUN, "morph:decade"
	}
	if isNumber(lower) {
		return chunky.TagNUM, "morph:num"
	}
	if strings.HasSuffix(lower, "%") && isNumber(lower[:len(lower)-1]) {
		return chunky.TagNUM, "morph:percent"
	}
	if isFraction(lower) {
		return chunky.TagNUM, "morph:fraction"
	}
	return 0, ""
}

// MorphCandidates returns candidate tags based on morphological features.
// isFirst suppresses the capitalization rule for sentence-initial words.
func MorphCandidates(word string, isFirst bool) (chunky.Tag, string) {
	var tags chunky.Tag
	add := func(ts ...chunky.Tag) {
		for _, t := range ts {
			tags |= t
		}
	}

	if len(word) == 0 {
		return 0, ""
	}

	lower := strings.ToLower(word)
	var rules []string

	if strings.Contains(word, "-") && word[0] >= 'A' && word[0] <= 'Z' {
		return chunky.TagADJ, "morph:hyphen+caps"
	}
	if !isFirst && word[0] >= 'A' && word[0] <= 'Z' {
		add(chunky.TagPROPN, chunky.TagADJ)
		rules = append(rules, "morph:caps")
	}

	norm := lower
	switch {
	case strings.HasSuffix(lower, "ies") && len(lower) > 4:
		norm = lower[:len(lower)-3] + "y"
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") &&
		!strings.HasSuffix(lower, "ous") && !strings.HasSuffix(lower, "us") &&
		!strings.HasSuffix(lower, "is") && len(lower) > 3:
		norm = lower[:len(lower)-1]
	}

	suffix := ""
	switch {
	case strings.HasSuffix(norm, "ly"):
		add(chunky.TagADV)
		suffix = "-ly"
	case strings.HasSuffix(norm, "ian"),
		strings.HasSuffix(norm, "ese"),
		strings.HasSuffix(norm, "ish"):
		add(chunky.TagADJ, chunky.TagPROPN, chunky.TagNOUN)
		suffix = "-ian/-ese/-ish"
	case strings.HasSuffix(norm, "tion"),
		strings.HasSuffix(norm, "sion"),
		strings.HasSuffix(norm, "ment"),
		strings.HasSuffix(norm, "ness"),
		strings.HasSuffix(norm, "ance"),
		strings.HasSuffix(norm, "ence"),
		strings.HasSuffix(norm, "ship"),
		strings.HasSuffix(norm, "hood"),
		strings.HasSuffix(norm, "dom"),
		strings.HasSuffix(norm, "ism"),
		strings.HasSuffix(norm, "ure"),
		strings.HasSuffix(norm, "sis"),
		strings.HasSuffix(norm, "ia"),
		len(norm) >= 6 && strings.HasSuffix(norm, "ity"):
		add(chunky.TagNOUN)
		suffix = "-tion/-ment/-ness/..."
	case strings.HasSuffix(norm, "ist"):
		add(chunky.TagNOUN, chunky.TagADJ)
		suffix = "-ist/-ists"
	case strings.HasSuffix(norm, "ize"),
		strings.HasSuffix(norm, "ise"),
		strings.HasSuffix(norm, "ify"):
		add(chunky.TagVERB)
		suffix = "-ize/-ise/-ify"
	case strings.HasSuffix(norm, "ate"):
		add(chunky.TagVERB, chunky.TagNOUN, chunky.TagADJ)
		suffix = "-ate/-ates"
	case strings.HasSuffix(norm, "ous"),
		strings.HasSuffix(norm, "ful"),
		strings.HasSuffix(norm, "less"),
		strings.HasSuffix(norm, "able"),
		strings.HasSuffix(norm, "ible"),
		strings.HasSuffix(norm, "ive"),
		strings.HasSuffix(norm, "ical"):
		add(chunky.TagADJ)
		suffix = "-ous/-ful/-less/-able/-ible/-ive/-ical"
	case strings.HasSuffix(norm, "ic"):
		add(chunky.TagADJ)
		suffix = "-ic"
	case strings.HasSuffix(norm, "ing"):
		add(chunky.TagVERB, chunky.TagNOUN)
		suffix = "-ing"
	case strings.HasSuffix(norm, "ed"):
		add(chunky.TagVERB)
		suffix = "-ed"
	case strings.HasSuffix(norm, "al"):
		add(chunky.TagADJ, chunky.TagNOUN)
		suffix = "-al"
	case strings.HasSuffix(norm, "er"),
		strings.HasSuffix(norm, "or"):
		add(chunky.TagNOUN)
		suffix = "-er/-or"
	}
	if suffix != "" {
		rules = append(rules, "morph:"+suffix)
	}

	prefix := ""
	switch {
	case suffix != "-ly" && (strings.HasPrefix(lower, "re") ||
		strings.HasPrefix(lower, "over") ||
		strings.HasPrefix(lower, "under")):
		add(chunky.TagVERB)
		prefix = "re-/over-/under-"
	case suffix != "-ly" && (strings.HasPrefix(lower, "un") ||
		strings.HasPrefix(lower, "non") ||
		strings.HasPrefix(lower, "anti") ||
		strings.HasPrefix(lower, "pre") ||
		strings.HasPrefix(lower, "post") ||
		strings.HasPrefix(lower, "inter") ||
		strings.HasPrefix(lower, "intra") ||
		strings.HasPrefix(lower, "trans") ||
		strings.HasPrefix(lower, "extra")):
		add(chunky.TagADJ, chunky.TagNOUN)
		prefix = "un-/non-/anti-/pre-/..."
	}
	if prefix != "" {
		rules = append(rules, "morph:"+prefix)
	}

	if tags == 0 {
		return 0, ""
	}
	return tags, strings.Join(rules, "+")
}

func isDecade(s string) bool {
	if !strings.HasSuffix(s, "s") || len(s) < 3 {
		return false
	}
	stem := s[:len(s)-1]
	for _, r := range stem {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func isOrdinal(s string) bool {
	if strings.HasSuffix(s, "st") || strings.HasSuffix(s, "nd") ||
		strings.HasSuffix(s, "rd") || strings.HasSuffix(s, "th") {
		stem := s[:len(s)-2]
		if len(stem) == 0 {
			return false
		}
		for _, r := range stem {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	return false
}

func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	if s[0] == '+' || s[0] == '-' {
		s = s[1:]
		if len(s) == 0 {
			return false
		}
	}
	hasDigit := false
	for i, r := range s {
		switch {
		case r >= '0' && r <= '9':
			hasDigit = true
		case (r == '.' || r == ',') && i > 0 && i < len(s)-1:
		default:
			return false
		}
	}
	return hasDigit
}

// isFraction reports whether s is a fraction of the form NUM/NUM (e.g. "3/8", "1/2").
// Also handles the Penn Treebank escaped form "3\/8" where the slash is preceded
// by a literal backslash.
func isFraction(s string) bool {
	sep := "/"
	i := strings.Index(s, "\\/")
	if i > 0 {
		sep = "\\/"
	} else {
		i = strings.Index(s, "/")
	}
	if i <= 0 || i+len(sep) >= len(s) {
		return false
	}
	return isNumber(s[:i]) && isNumber(s[i+len(sep):])
}

func isAlpha(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
