package chunky

// AbbreviationTags maps common abbreviations to their UD tags.
// Both dotted and undotted forms are included since the tokenizer
// may strip trailing periods depending on context.
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
	"vol":  {TagNOUN},
	"vol.": {TagNOUN},
	"pp":   {TagNOUN},
	"pp.":  {TagNOUN},
}
