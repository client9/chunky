package tok

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/client9/chunky"
	"github.com/client9/typewriter"
)

type rawToken struct {
	word   string
	offset int
}

// surfaceTokenizeRaw splits s into tokens, recording each token's byte offset
// in the original string.
func surfaceTokenizeRaw(s string) []rawToken {

	// TODO: make a global?  It can be reused concurrently.
    	tw := typewriter.New(typewriter.Config{
		// Unclear what the right behavior is with other unicode symboles
		// so start with quotes, dashes, and spaces.
        	Categories: typewriter.Quotes | typewriter.Dashes | typewriter.Spaces,
    	})

	out := make([]rawToken, 0, 16)
	i := 0
	for i < len(s) {
		// skip unicode whitespace
		for i < len(s) {
			r, size := utf8.DecodeRuneInString(s[i:])
			if !unicode.IsSpace(r) {
				break
			}
			i += size
		}
		if i >= len(s) {
			break
		}
		// find end of whitespace-delimited field
		start := i
		for i < len(s) {
			r, size := utf8.DecodeRuneInString(s[i:])
			if unicode.IsSpace(r) {
				break
			}
			i += size
		}

		//
		// unicode normalize with client/typewriter
		//
		// TODO: add client9/demoji once it's stable.
		//
		p := s[start:i]
		p = tw.Replace(p)

		//
		// remove [items in square brackets]
		//
		p = stripInlineCitations(p)
		pos := start

		if len(p) == 0 {
			continue
		}

		if len(p) == 1 {
			out = append(out, rawToken{p, pos})
			continue
		}

		// strip leading '('
		if p[0] == '(' {
			out = append(out, rawToken{"(", pos})
			p = p[1:]
			pos++
		}
		if len(p) == 0 {
			continue
		}

		// strip trailing sentence punctuation, but keep the dot when the whole
		// word is a known dotted abbreviation (p., Dr., etc.) so it stays atomic.
		last, lastPos := "", 0
		if ch := p[len(p)-1]; ch == ',' || ch == '.' || ch == ':' || ch == ';' || ch == '!' || ch == '?' {
			if ch != '.' || !chunky.DottedAbbreviations[strings.ToLower(p)] {
				last = string(ch)
				lastPos = pos + len(p) - 1
				p = p[:len(p)-1]
			}
		}
		if len(p) == 0 {
			if last != "" {
				out = append(out, rawToken{last, lastPos})
			}
			continue
		}

		// strip trailing ')'
		if p[len(p)-1] == ')' {
			out = append(out, rawToken{p[:len(p)-1], pos})
			out = append(out, rawToken{")", pos + len(p) - 1})
			if last != "" {
				out = append(out, rawToken{last, lastPos})
			}
			continue
		}

		out = append(out, rawToken{p, pos})
		if last != "" {
			out = append(out, rawToken{last, lastPos})
		}
	}
	return splitContractions(out)
}

// contractionSuffixes is the set of apostrophe-led suffixes that trigger a split.
var contractionSuffixes = map[string]bool{
	"'ll": true, "'re": true, "'ve": true, "'m": true, "'d": true, "'s": true,
	"'t": true,
}

// splitContractions expands contraction tokens into (stem, suffix) pairs.
// It fires on words whose suffix (from the apostrophe onward) is a known
// contraction suffix. Irregular forms in chunky.ContractionNorm are handled
// first. Words in chunky.AbbreviationTags that are NOT in ContractionNorm
// stay whole (e.g., "ain't").
func splitContractions(tokens []rawToken) []rawToken {
	out := make([]rawToken, 0, len(tokens)+4)
	for _, t := range tokens {
		lower := strings.ToLower(t.word)

		// Irregular forms: won't → will + n't, shan't → shall + n't.
		if parts, ok := chunky.ContractionNorm[lower]; ok {
			out = append(out, rawToken{parts[0], t.offset})
			out = append(out, rawToken{parts[1], t.offset + len(t.word) - len(parts[1])})
			continue
		}

		// Words that stay whole (ain't, o'clock, etc.) — already in AbbreviationTags.
		if _, ok := chunky.AbbreviationTags[lower]; ok {
			out = append(out, t)
			continue
		}

		ap := strings.IndexByte(t.word, '\'')
		if ap <= 0 {
			out = append(out, t)
			continue
		}

		suffix := lower[ap:]

		// n't: move the n to the suffix when the stem-without-n is in the lexicon
		// (don't→do+n't, shouldn't→should+n't); keep n in stem otherwise
		// (can't→can+n't, where "ca" is not a word).
		if suffix == "'t" && ap >= 2 && (t.word[ap-1] == 'n' || t.word[ap-1] == 'N') {
			stemNoN := strings.ToLower(t.word[:ap-1])
			if len(wordtagmap[stemNoN]) > 0 || len(chunky.AbbreviationTags[stemNoN]) > 0 {
				// n belongs to suffix: don't→do+n't, shouldn't→should+n't
				out = append(out, rawToken{t.word[:ap-1], t.offset})
				out = append(out, rawToken{"n't", t.offset + ap - 1})
			} else {
				// n belongs to stem: can't→can+'t (no byte overlap)
				out = append(out, rawToken{t.word[:ap], t.offset})
				out = append(out, rawToken{"'t", t.offset + ap})
			}
			continue
		}

		if contractionSuffixes[suffix] {
			out = append(out, rawToken{t.word[:ap], t.offset})
			out = append(out, rawToken{t.word[ap:], t.offset + ap})
			continue
		}

		out = append(out, t)
	}
	return out
}

// SurfaceTokenize converts a string into a list of "surface tokens".
func SurfaceTokenize(s string) []string {
	raw := surfaceTokenizeRaw(s)
	out := make([]string, len(raw))
	for i, r := range raw {
		out[i] = r.word
	}
	return out
}

type Token struct {
	Word      string
	Offset    int
	Canidates []chunky.Tag
	Rule      string // which rule assigned the candidates
}

func (t Token) IsUnknownTag() bool {
	return len(t.Canidates) == 0 || t.Canidates[0] == chunky.TagUNK
}

func (t Token) HasTag(x chunky.Tag) bool {
	for _, r := range t.Canidates {
		if r == x {
			return true
		}
	}
	return false
}

func (t Token) String() string {
	if len(t.Canidates) == 1 {
		return t.Word + "/" + t.Canidates[0].String()
	}
	parts := make([]string, len(t.Canidates))
	for i, s := range t.Canidates {
		parts[i] = s.String()
	}
	return t.Word + "/{" + strings.Join(parts, ",") + "}"
}

func TagString(s string) []Token {
	raw := surfaceTokenizeRaw(s)
	out := make([]Token, len(raw))

	for i, r := range raw {
		lower := strings.ToLower(r.word)
		candidates := wordtagmap[lower]
		rule := "lexicon"
		if len(candidates) == 0 {
			// AbbreviationTags is the runtime-editable override layer; entries
			// added there (e.g., contraction suffixes) don't require regenerating
			// the compiled lexicon.
			if tags, ok := chunky.AbbreviationTags[lower]; ok {
				candidates = tags
			} else {
				rule = ""
			}
		}
		out[i] = Token{
			Word:      r.word,
			Offset:    r.offset,
			Canidates: candidates,
			Rule:      rule,
		}
	}
	return out
}

// isDecade returns true for decade forms: 1980s, 1960s, 2000s, ...
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

// isOrdinal returns true for numeric ordinals: 1st, 2nd, 3rd, 4th, 20th, ...
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

// isNumber returns true for numeric strings: integers, decimals, and
// comma-formatted numbers (e.g. 1,000 or 3.14).
func isNumber(s string) bool {
	if len(s) == 0 {
		return false
	}
	// strip optional leading sign
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
			// allow interior separators only
		default:
			return false
		}
	}
	return hasDigit
}

// TODO: evaluate additional suffix rules using corpus data:
//   -ium, -ine  mostly NOUN (physics/chemistry: calcium, chlorine)
//   -ary        NOUN/ADJ ambiguous (library, military)
//   -ogy        NOUN (biology, theology)
//   -ies        NOUN/VERB ambiguous when not a plural (series, species)

// MorphCandidates returns candidate tags for an unknown word based on
// morphological features. isFirst is true when the word is sentence-initial
// (suppresses the capitalization rule). Returns candidates and a rule ID.
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

	// Ordinals: 1st, 2nd, 3rd, 20th, ... → ADJ (corpus confirms universally)
	if isOrdinal(lower) {
		return []chunky.Tag{chunky.TagADJ}, "morph:ordinal"
	}

	// Decades: 1980s, 1960s → NOUN/NUM (corpus ~2:1 NOUN)
	if isDecade(lower) {
		return []chunky.Tag{chunky.TagNOUN, chunky.TagNUM}, "morph:decade"
	}

	// Numbers: plain integers, decimals (3.14), and formatted (1,000.50).
	if isNumber(lower) {
		return []chunky.Tag{chunky.TagNUM}, "morph:num"
	}

	// Percentages: 15%, 1.5%, 100%
	if strings.HasSuffix(lower, "%") && isNumber(lower[:len(lower)-1]) {
		return []chunky.Tag{chunky.TagNUM}, "morph:percent"
	}

	// Hyphenated capitalized words are almost always adjectives.
	if strings.Contains(word, "-") && word[0] >= 'A' && word[0] <= 'Z' {
		return []chunky.Tag{chunky.TagADJ}, "morph:hyphen+caps"
	}

	// Non-sentence-initial capitalized word → proper noun or adjective.
	if !isFirst && word[0] >= 'A' && word[0] <= 'Z' {
		add(chunky.TagPROPN, chunky.TagADJ)
		rules = append(rules, "morph:caps")
	}

	// Normalize away plural endings before suffix matching so that e.g.
	// "morphisms" → "morphism", "ethnicities" → "ethnicity" match their rules.
	norm := lower
	switch {
	case strings.HasSuffix(lower, "ies") && len(lower) > 4:
		norm = lower[:len(lower)-3] + "y"
	case strings.HasSuffix(lower, "s") && !strings.HasSuffix(lower, "ss") &&
		!strings.HasSuffix(lower, "ous") && !strings.HasSuffix(lower, "us") &&
		!strings.HasSuffix(lower, "is") && len(lower) > 3:
		norm = lower[:len(lower)-1]
	}

	// Suffix rules (longest match first within each group).
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

	// Prefix rules (additive on top of suffix results).
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

	// -ies: "flies" → "fly"
	if strings.HasSuffix(lower, "ies") && len(lower) > 3 {
		try("-ies", lower[:len(lower)-3]+"y")
	}
	// -s: "accelerates" → "accelerate", "cats" → "cat"
	if strings.HasSuffix(lower, "s") && len(lower) > 2 {
		try("-s", lower[:len(lower)-1])
	}
	// -es: "boxes" → "box"
	if strings.HasSuffix(lower, "es") && len(lower) > 3 {
		try("-es", lower[:len(lower)-2])
	}

	// -ing: "walking" → "walk", "hoping" → "hope", "running" → "run"
	if strings.HasSuffix(lower, "ing") && len(lower) > 4 {
		stem := lower[:len(lower)-3]
		try("-ing", stem)
		try("-ing+e", stem+"e")
		tryDoubled("-ing+double", stem)
	}

	// -ed: "walked" → "walk", "hoped" → "hope", "stopped" → "stop"
	if strings.HasSuffix(lower, "ed") && len(lower) > 3 {
		stem := lower[:len(lower)-2]
		try("-ed", stem)
		try("-ed+e", stem+"e")
		tryDoubled("-ed+double", stem)
	}

	// -er: "faster" → "fast", "nicer" → "nice", "bigger" → "big"
	if strings.HasSuffix(lower, "er") && len(lower) > 3 {
		stem := lower[:len(lower)-2]
		try("-er", stem)
		try("-er+e", stem+"e")
		tryDoubled("-er+double", stem)
	}

	// -est: "fastest" → "fast", "nicest" → "nice"
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

// hyphenAdjSuffixes are suffixes that always produce ADJ in a hyphenated
// compound regardless of how the suffix is tagged as a standalone word.
var hyphenAdjSuffixes = map[string]bool{
	"like": true, // flu-like, war-like, life-like
	"free": true, // sugar-free, tax-free
	"wide": true, // nationwide, industry-wide
}

// HyphenCandidates handles lowercase hyphenated compounds by looking up the
// last component in the lexicon and applying morph rules to it as a fallback.
func HyphenCandidates(word string) ([]chunky.Tag, string) {
	i := strings.LastIndex(word, "-")
	if i < 0 || i == len(word)-1 {
		return nil, ""
	}
	last := strings.ToLower(word[i+1:])

	// Known adjectival suffixes override the lexicon.
	if hyphenAdjSuffixes[last] {
		return []chunky.Tag{chunky.TagADJ}, "hyphen:adj-suffix:" + last
	}

	// Try lexicon lookup on the last component.
	if tags, ok := wordtagmap[last]; ok {
		return tags, "hyphen:lexicon:" + last
	}

	// Try inflection on the last component.
	if tags, rule := InflectionCandidates(last); tags != nil {
		return tags, "hyphen:" + rule
	}

	// Try morph suffix rules on the last component.
	if tags, rule := MorphCandidates(last, false); tags != nil {
		return tags, "hyphen:" + rule
	}

	return nil, ""
}

// Unk1 applies context-based rules for unknown words using neighboring tags.
// isAlpha returns true if every rune in s is a Unicode letter.
// Digits, symbols, and mixed-script words return false and stay untagged.
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

func TagUnknowns(tokens []Token) []Token {
	for i, t := range tokens {
		if !t.IsUnknownTag() {
			continue
		}
		if candidates, rule := InflectionCandidates(t.Word); candidates != nil {
			tokens[i].Canidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := HyphenCandidates(t.Word); candidates != nil {
			tokens[i].Canidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if candidates, rule := MorphCandidates(t.Word, i == 0); candidates != nil {
			tokens[i].Canidates = candidates
			tokens[i].Rule = rule
			continue
		}
		if isAlpha(t.Word) {
			tokens[i].Canidates = []chunky.Tag{chunky.TagNOUN}
			tokens[i].Rule = "unk:word"
		}
	}
	return tokens
}
