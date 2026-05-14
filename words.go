package chunky

// WordTags maps known open-category words that are missing from the
// generated lexicon. Unlike ClosedFormTags, these are content words
// added by hand as gaps are discovered.
var WordTags = map[string][]Tag{
	// punctuation overrides — Brown corpus has noise tags on these
	// (e.g. "," tagged FW-RB-TL, ":" tagged NP/IN)
	",":      {TagPUNCT},
	".":      {TagPUNCT},
	":":      {TagPUNCT},
	";":      {TagPUNCT},
	"!":      {TagPUNCT},
	"?":      {TagPUNCT},
	"(":      {TagPUNCT},
	")":      {TagPUNCT},
	"-":      {TagPUNCT},
	"--":     {TagPUNCT},
	"\u2013": {TagPUNCT}, // en dash
	"\u2014": {TagPUNCT}, // em dash

	// -ede verbs: too few to justify a suffix rule, all clearly VERB
	"accede":    {TagVERB},
	"concede":   {TagVERB},
	"impede":    {TagVERB},
	"precede":   {TagVERB},
	"recede":    {TagVERB},
	"secede":    {TagVERB},
	"stampede":  {TagNOUN, TagVERB},
	"supersede": {TagVERB},
}
