package chunky

// WordTags maps known open-category words that are missing from the
// generated lexicon. Unlike ClosedFormTags, these are content words
// added by hand as gaps are discovered.
var WordTags = map[string]Tag{
	// punctuation overrides — Brown corpus has noise tags on these
	// (e.g. "," tagged FW-RB-TL, ":" tagged NP/IN)
	",":  TagPUNCT,
	".":  TagPUNCT,
	":":  TagPUNCT,
	";":  TagPUNCT,
	"!":  TagPUNCT,
	"?":  TagPUNCT,
	"(":  TagPUNCT,
	")":  TagPUNCT,
	"-":  TagPUNCT,
	"--": TagPUNCT,
	"{":  TagPUNCT, // TBD - probably strip brace blocks
	"}":  TagPUNCT, // TBD -
	"$":  TagSYM,   // as per spaCy
	"%":  TagNOUN,  // as per spaCy
	"_":  TagPUNCT,
	"\"": TagPUNCT,
	"''": TagPUNCT, // closing double quote (Penn Treebank convention)
	"``": TagPUNCT, // opening double quote (Penn Treebank convention)
	"–":  TagPUNCT, // en dash
	"—":  TagPUNCT, // em dash

	// Brown corpus assigns DET to these but spaCy/UD always tags them ADJ.
	// They are gradable prenominal modifiers, never true determiners.
	"various":    TagADJ,
	"fewer":      TagADJ,
	"lesser":     TagADJ,
	"certain":    TagADJ,
	"particular": TagADJ,

	// Brown corpus over-tagged these with JJ; spaCy/UD always ADV (n≥1k).
	"just":       TagADV,
	"even":       TagADV,
	"rather":     TagADV,
	"indeed":     TagADV,
	"newly":      TagADV,
	"similarly":  TagADV,
	"elsewhere":  TagADV,
	"everywhere": TagADV,

	// Brown JJ but always ADJ in practice (n≥500, spaCy ≥96%).
	"strong":      TagADJ,
	"subsequent":  TagADJ,
	"widespread":  TagADJ,
	"previous":    TagADJ,
	"appropriate": TagADJ,
	"ready":       TagADJ,
	"usual":       TagADJ,
	"false":       TagADJ,
	"large":       TagADJ,
	"full":        TagADJ,

	// Brown ADJ|NOUN but spaCy ≥93% ADJ; these words are essentially never
	// used as standalone noun heads in prose — force to ADJ.
	// Note: "main" omitted — "Main Street" (proper name) relies on the NOUN
	// candidate before RetagCapitalized runs to correctly disambiguate "down Main St."
	"small":        TagADJ,
	"common":       TagADJ,
	"local":        TagADJ,
	"original":     TagADJ,
	"annual":       TagADJ,
	"commercial":   TagADJ,
	"active":       TagADJ,
	"private":      TagADJ,
	"professional": TagADJ,
	"musical":      TagADJ,

	// Brown tagged UH (interjection) alongside VB; overwhelmingly VERB in prose.
	// "et" is a genuine Latin foreign word → TagX.
	"see":    TagVERB,
	"please": TagVERB,

	// Brown tagged as ADV|PART (RP); always ADV as a degree modifier.
	"quite": TagADV,

	// {ADJ,ADV,DET}: always ADV in practice (degree intensifier).
	"very": TagADV,

	// {ADV,NOUN}: always ADV in prose.
	"somewhere": TagADV,
	"outdoors":  TagADV,
	"nowhere":   TagADV,

	// {ADJ,ADV,NOUN}: overwhelmingly ADJ in prose (spaCy ≥97%).
	"true": TagADJ,

	// {ADJ,ADV,NOUN}: always ADV in prose — Brown noise added ADJ/NOUN.
	"posthumously":    TagADV,
	"interchangeably": TagADV,

	// {ADV,NOUN}: used as NOUN modifier in titles ("vice president", "vice versa").
	"vice": TagNOUN,

	// {ADV,NOUN}: always NOUN in prose (Brown noise added ADV).
	"offs":     TagNOUN,
	"branding": TagNOUN,
	"ante":     TagNOUN,
	"meantime": TagNOUN,

	// Brown tagged as NOUN|NUM; NUM in UD for cardinal/quantity use (dominant).
	// Inflected forms listed explicitly because InflectionCandidates reads wordtagmap
	// (raw lexicon) and wouldn't see the "zero" override.
	"zero":    TagNUM,
	"zeros":   TagNUM,
	"zeroes":  TagNUM,
	"zeroed":  TagVERB,
	"zeroing": TagVERB,

	// Brown UH noise gave these spurious AUX/PRON/X tags; all pure interjections.
	"ah":  TagINTJ,
	"hey": TagINTJ,
	"oh":  TagINTJ,
	"uh":  TagINTJ,
	"um":  TagINTJ,

	// Brown tagged as multiple including VERB; spaCy always VERB (n≥500, ≥99%).
	"born":        TagVERB,
	"using":       TagVERB,
	"deleted":     TagVERB,
	"fell":        TagVERB,
	"broke":       TagVERB,
	"trying":      TagVERB,
	"suggested":   TagVERB,
	"done":        TagVERB,
	"consisting":  TagVERB,
	"depending":   TagVERB,
	"considering": TagVERB,
	"excluding":   TagVERB,
	"providing":   TagVERB,
	"seeing":      TagVERB,

	// Brown tagged as ADV/DET/PART; spaCy/UD always ADJ (n=19k, 99.7%).
	"many": TagADJ,

	// Brown tagged "i"/"me" as PROPN in "CHAPTER I" / title contexts; always PRON in prose.
	"i":  TagPRON,
	"me": TagPRON,

	// Abbreviations used as common nouns — not proper names — so capitalized
	// forms (TV, CDs, ADRs) are not promoted to PROPN by RetagCapitalized.
	"tv":     TagNOUN,
	"cds":    TagNOUN,
	"cd":     TagNOUN,
	"adrs":   TagNOUN,
	"adr":    TagNOUN,
	"dvd":    TagNOUN,
	"dvds":   TagNOUN,
	"pc":     TagNOUN,
	"pcs":    TagNOUN,
	"atm":    TagNOUN,
	"remics": TagNOUN,
	"remic":  TagNOUN,

	// Brown noise on symbols and conjunctions.
	"&":      TagCCONJ,
	"and/or": TagCCONJ,

	// Brown tagged as DET; spaCy always DET (n=5k, 95.7%).
	"every": TagDET,

	// Brown added spurious ADV/ADP tag; spaCy always NOUN (n≥100, ≥95%).
	"lots":  TagNOUN,
	"spite": TagNOUN,

	// Brown tagged as NOUN only; should be DET (prenominal: "no reason", "no one").
	"no": TagDET,

	// gerund-prepositions: these words introduce PP chunks and are never
	// used as verbs in the prose target; force to ADP to resolve ambiguity.
	"including": TagADP,
	"involving": TagADP,
	"regarding": TagADP,

	// gerund-nouns: Brown corpus tagged only as VERB but these are frequently
	// used as noun heads or noun modifiers ("trading volume", "operating profit",
	// "banking sector"). Adding NOUN candidate lets DisambiguateByChunk resolve
	// via chunk position rather than always forcing VP.
	"trading":    TagNOUN | TagVERB,
	"operating":  TagNOUN | TagVERB,
	"banking":    TagNOUN | TagVERB,
	"filing":     TagNOUN | TagVERB,
	"funding":    TagNOUN | TagVERB,
	"developing": TagNOUN | TagVERB,
	"smoking":    TagNOUN | TagVERB,
	"warming":    TagNOUN | TagVERB,
	"fuels":      TagNOUN | TagVERB,

	// -ede verbs: too few to justify a suffix rule, all clearly VERB
	"accede":    TagVERB,
	"concede":   TagVERB,
	"impede":    TagVERB,
	"precede":   TagVERB,
	"recede":    TagVERB,
	"secede":    TagVERB,
	"stampede":  TagNOUN | TagVERB,
	"supersede": TagVERB,
}
