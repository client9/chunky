package tok

import (
	"strings"
	"unicode"

	"github.com/client9/chunky"
)

// TagUnknowns fills in candidate tags for tokens that have none, trying rules
// in order: inflection → hyphen → morphology → alpha fallback.
func TagUnknowns(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.IsUnknownTag() {
			continue
		}
		if candidates, rule := InflectionCandidates(t.Word); candidates != nil {
			tokens[i].Candidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := HyphenCandidates(t.Word); candidates != nil {
			tokens[i].Candidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := MorphCandidates(t.Word, i == 0); candidates != nil {
			tokens[i].Candidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if isAlpha(t.Word) {
			tokens[i].Candidates = []chunky.Tag{chunky.TagNOUN}
			tokens[i].Rule = "unk:word"
		}
	}
	return tokens
}

// InflectionCandidates looks up candidate tags by stripping common inflectional
// suffixes and checking the resulting stem against the lexicon. Returns
// candidates and a rule ID of the form "inflect:<suffix>:<stem>".
func InflectionCandidates(word string) ([]chunky.Tag, string) {
	lower := strings.ToLower(word)

	seen := make(map[chunky.Tag]bool)
	matchedStem := ""
	matchedSuffix := ""

	try := func(suffix, stem string) {
		tags := wordtagmap[stem]
		if len(tags) == 0 {
			return
		}
		for _, t := range tags {
			seen[t] = true
		}
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

	if len(seen) == 0 {
		return nil, ""
	}
	out := make([]chunky.Tag, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	return out, "inflect:" + matchedSuffix + ":" + matchedStem
}

// hyphenAdjSuffixes always produce ADJ in a hyphenated compound.
var hyphenAdjSuffixes = map[string]bool{
	"like": true,
	"free": true,
	"wide": true,
}

// HyphenCandidates handles hyphenated compounds by looking up the last
// component in the lexicon and applying morph rules as a fallback.
func HyphenCandidates(word string) ([]chunky.Tag, string) {
	i := strings.LastIndex(word, "-")
	if i < 0 || i == len(word)-1 {
		return nil, ""
	}
	last := strings.ToLower(word[i+1:])

	if hyphenAdjSuffixes[last] {
		return []chunky.Tag{chunky.TagADJ}, "hyphen:adj-suffix:" + last
	}
	if tags, ok := wordtagmap[last]; ok {
		return tags, "hyphen:lexicon:" + last
	}
	if tags, rule := InflectionCandidates(last); tags != nil {
		return tags, "hyphen:" + rule
	}
	if tags, rule := MorphCandidates(last, false); tags != nil {
		return tags, "hyphen:" + rule
	}
	return nil, ""
}

// TODO: evaluate additional suffix rules using corpus data:
//
//	-ium, -ine  mostly NOUN (physics/chemistry: calcium, chlorine)
//	-ary        NOUN/ADJ ambiguous (library, military)
//	-ogy        NOUN (biology, theology)
//	-ies        NOUN/VERB ambiguous when not a plural (series, species)

// MorphCandidates returns candidate tags based on morphological features.
// isFirst suppresses the capitalization rule for sentence-initial words.
func MorphCandidates(word string, isFirst bool) ([]chunky.Tag, string) {
	seen := make(map[chunky.Tag]bool)
	add := func(tags ...chunky.Tag) {
		for _, t := range tags {
			seen[t] = true
		}
	}

	if len(word) == 0 {
		return nil, ""
	}

	lower := strings.ToLower(word)
	var rules []string

	if isOrdinal(lower) {
		return []chunky.Tag{chunky.TagADJ}, "morph:ordinal"
	}
	if isDecade(lower) {
		return []chunky.Tag{chunky.TagNOUN, chunky.TagNUM}, "morph:decade"
	}
	if isNumber(lower) {
		return []chunky.Tag{chunky.TagNUM}, "morph:num"
	}
	if strings.HasSuffix(lower, "%") && isNumber(lower[:len(lower)-1]) {
		return []chunky.Tag{chunky.TagNUM}, "morph:percent"
	}
	if strings.Contains(word, "-") && word[0] >= 'A' && word[0] <= 'Z' {
		return []chunky.Tag{chunky.TagADJ}, "morph:hyphen+caps"
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
		add(chunky.TagVERB, chunky.TagNOUN, chunky.TagADJ)
		suffix = "-ing"
	case strings.HasSuffix(norm, "ed"):
		add(chunky.TagVERB, chunky.TagADJ)
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
	case strings.HasPrefix(lower, "re"),
		strings.HasPrefix(lower, "over"),
		strings.HasPrefix(lower, "under"):
		add(chunky.TagVERB)
		prefix = "re-/over-/under-"
	case strings.HasPrefix(lower, "un"),
		strings.HasPrefix(lower, "non"),
		strings.HasPrefix(lower, "anti"),
		strings.HasPrefix(lower, "pre"),
		strings.HasPrefix(lower, "post"),
		strings.HasPrefix(lower, "inter"),
		strings.HasPrefix(lower, "intra"),
		strings.HasPrefix(lower, "trans"),
		strings.HasPrefix(lower, "extra"):
		add(chunky.TagADJ, chunky.TagNOUN)
		prefix = "un-/non-/anti-/pre-/..."
	}
	if prefix != "" {
		rules = append(rules, "morph:"+prefix)
	}

	if len(seen) == 0 {
		return nil, ""
	}
	out := make([]chunky.Tag, 0, len(seen))
	for t := range seen {
		out = append(out, t)
	}
	return out, strings.Join(rules, "+")
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
