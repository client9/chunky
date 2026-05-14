package chunky

// DottedAbbreviations is the set of forms that the surface tokenizer must keep
// atomic — the trailing dot is part of the abbreviation, not sentence-ending
// punctuation. Entries are lowercase with the dot included.
// Add domain-specific forms here (e.g. medical, legal) at startup.
var DottedAbbreviations = map[string]bool{
	// Titles
	"mr.": true, "mrs.": true, "ms.": true, "dr.": true,
	"prof.": true, "rev.": true, "jr.": true, "sr.": true, "st.": true,

	// Discourse / Latin
	"etc.": true, "vs.": true, "e.g.": true, "i.e.": true,

	// Time
	"a.m.": true, "p.m.": true,

	// Geographic
	"u.s.": true, "u.k.": true, "u.s.a.": true,

	// Organizational
	"inc.": true, "ltd.": true, "corp.": true, "co.": true,

	// Reference / bibliographic
	"p.": true, "pp.": true, "vol.": true, "fig.": true, "no.": true,
}

// AbbreviationTags maps common abbreviations (dotted and undotted) to their
// UD tags. Used at tagging time; the tokenizer uses DottedAbbreviations to
// decide whether to keep a trailing dot attached.
var AbbreviationTags = map[string][]Tag{
	// Contraction suffixes produced by the surface tokenizer split.
	// 's is ambiguous: AUX (copula: "it's") or PART (possessive marker: "John's").
	"'ll": {TagAUX},
	"'re": {TagAUX},
	"'ve": {TagAUX},
	"'m":  {TagAUX},
	"'d":  {TagAUX},
	"'s":  {TagAUX, TagPART},
	"n't": {TagADV},
	"'t":  {TagADV},

	// Irregular contractions kept whole (ContractionNorm handles won't/shan't split).
	"won't":  {TagAUX},
	"wont":   {TagAUX},
	"ain't":  {TagAUX},
	"aint":   {TagAUX},
	"shan't": {TagAUX},
	"shant":  {TagAUX},
	"gonna":  {TagVERB},
	"wanna":  {TagVERB},
	"gotta":  {TagVERB},

	// Discourse / Latin
	"e.g":  {TagADV},
	"e.g.": {TagADV},
	"i.e":  {TagADV},
	"i.e.": {TagADV},
	"etc":  {TagADV},
	"etc.": {TagADV},
	"vs":   {TagADP},
	"vs.":  {TagADP},

	// Time
	"a.m":  {TagADV},
	"a.m.": {TagADV},
	"am":   {TagADV},
	"p.m":  {TagADV},
	"p.m.": {TagADV},
	"pm":   {TagADV},

	// Titles (part of a proper name → PROPN)
	"mr":    {TagPROPN},
	"mr.":   {TagPROPN},
	"mrs":   {TagPROPN},
	"mrs.":  {TagPROPN},
	"ms":    {TagPROPN},
	"ms.":   {TagPROPN},
	"dr":    {TagPROPN},
	"dr.":   {TagPROPN},
	"prof":  {TagPROPN},
	"prof.": {TagPROPN},
	"rev":   {TagPROPN},
	"rev.":  {TagPROPN},
	"jr":    {TagPROPN},
	"jr.":   {TagPROPN},
	"sr":    {TagPROPN},
	"sr.":   {TagPROPN},
	"st":    {TagPROPN},
	"st.":   {TagPROPN},

	// Geographic
	"u.s":    {TagPROPN},
	"u.s.":   {TagPROPN},
	"u.k":    {TagPROPN},
	"u.k.":   {TagPROPN},
	"u.s.a":  {TagPROPN},
	"u.s.a.": {TagPROPN},

	// Organizational suffixes
	"inc":   {TagPROPN},
	"inc.":  {TagPROPN},
	"ltd":   {TagPROPN},
	"ltd.":  {TagPROPN},
	"corp":  {TagPROPN},
	"corp.": {TagPROPN},
	"co":    {TagPROPN},
	"co.":   {TagPROPN},

	// Miscellaneous
	"fig":  {TagNOUN},
	"fig.": {TagNOUN},
	"no":   {TagNOUN},
	"no.":  {TagNOUN},
	"p":    {TagNOUN},
	"p.":   {TagNOUN},
	"vol":  {TagNOUN},
	"vol.": {TagNOUN},
	"pp":   {TagNOUN},
	"pp.":  {TagNOUN},
}
