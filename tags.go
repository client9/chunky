package chunky

import (
	"fmt"
	"strings"
)

type Tag uint32

const (
	TagUNK Tag = 1 << iota
	TagADJ
	TagADP
	TagADV
	TagAUX
	TagDET
	TagCCONJ
	TagINTJ
	TagNOUN
	TagNUM
	TagPART
	TagPRON
	TagPROPN
	TagPUNCT
	TagSCONJ
	TagSYM
	TagVERB
	TagX
)

func (t Tag) String() string {
	switch t {
	case 0:
		return "<NONE>"
	case TagUNK:
		return "<UNK>"
	case TagADJ:
		return "ADJ"
	case TagADP:
		return "ADP"
	case TagADV:
		return "ADV"
	case TagAUX:
		return "AUX"
	case TagCCONJ:
		return "CCONJ"
	case TagDET:
		return "DET"
	case TagINTJ:
		return "INTJ"
	case TagNOUN:
		return "NOUN"
	case TagNUM:
		return "NUM"
	case TagPART:
		return "PART"
	case TagPRON:
		return "PRON"
	case TagPROPN:
		return "PROPN"
	case TagPUNCT:
		return "PUNCT"
	case TagSCONJ:
		return "SCONJ"
	case TagSYM:
		return "SYM"
	case TagVERB:
		return "VERB"
	case TagX:
		return "X"
	default:
		panic("unknown tag")
	}
}

// AllTags lists every tag in a stable order for display/iteration.
var AllTags = []Tag{
	TagADJ, TagADP, TagADV, TagAUX, TagCCONJ, TagDET, TagINTJ,
	TagNOUN, TagNUM, TagPART, TagPRON, TagPROPN, TagPUNCT,
	TagSCONJ, TagSYM, TagVERB, TagX, TagUNK,
}

var tagNames = map[string]Tag{
	"<UNK>": TagUNK,
	"ADJ":   TagADJ,
	"ADP":   TagADP,
	"ADV":   TagADV,
	"AUX":   TagAUX,
	"CCONJ": TagCCONJ,
	"DET":   TagDET,
	"INTJ":  TagINTJ,
	"NOUN":  TagNOUN,
	"NUM":   TagNUM,
	"PART":  TagPART,
	"PRON":  TagPRON,
	"PROPN": TagPROPN,
	"PUNCT": TagPUNCT,
	"SCONJ": TagSCONJ,
	"SYM":   TagSYM,
	"VERB":  TagVERB,
	"X":     TagX,
}

// ParseTag converts a string to a Tag, returning an error if unrecognized.
func ParseTag(s string) (Tag, error) {
	if t, ok := tagNames[s]; ok {
		return t, nil
	}
	return TagUNK, fmt.Errorf("pos: unknown tag %q", s)
}

// TagFromPennTag converts a Penn Treebank POS tag to UD candidate tags.
// The returned Tag may have multiple bits set when the Penn tag is ambiguous
// in UD (e.g. IN → ADP|SCONJ, VB* → VERB|AUX). Returns 0 for unrecognized tags.
//
// Key Penn→UD ambiguities:
//   - VB/VBD/VBG/VBN/VBP/VBZ → VERB|AUX  (Penn has no AUX tag)
//   - IN → ADP|SCONJ           (Penn conflates prepositions and subordinators)
//   - TO → PART|ADP            (infinitive marker vs preposition)
//   - RB → ADV|PART            (Penn uses RB for "not/n't" which is PART in UD)
//   - WRB → ADV|SCONJ          (when/where used as subordinators are SCONJ in UD)
//   - JJ → ADJ|NOUN|DET|ADV   (Penn uses JJ for prenominal nouns, quantifiers, degree words)
//   - WDT → DET|PRON           (relative "which/that")
//   - RBS → ADV|DET            (Penn uses RBS for quantifier "most" which is DET in UD)
func TagFromPennTag(s string) Tag {
	switch s {
	case "NN", "NNS":
		return TagNOUN
	case "NNP", "NNPS":
		return TagPROPN
	case "VB", "VBD", "VBG", "VBN", "VBP", "VBZ":
		return TagVERB | TagAUX
	case "MD":
		return TagAUX
	case "JJ":
		return TagADJ | TagNOUN | TagDET | TagADV
	case "JJR":
		return TagADJ | TagADV
	case "JJS":
		return TagADJ | TagADV
	case "RB":
		return TagADV | TagPART
	case "RBR":
		return TagADV
	case "RBS":
		return TagADV | TagDET
	case "WRB":
		return TagADV | TagSCONJ
	case "DT", "PDT":
		return TagDET | TagPRON
	case "WDT":
		return TagDET | TagPRON
	case "IN":
		return TagADP | TagSCONJ
	case "TO":
		return TagPART | TagADP
	case "PRP", "PRP$", "WP", "WP$":
		return TagPRON
	case "EX":
		return TagPRON
	case "CC":
		return TagCCONJ
	case "RP":
		return TagPART
	case "POS":
		return TagPART
	case "CD":
		return TagNUM
	case "UH":
		return TagINTJ
	case ".", ",", ":", "''", "``", "-LRB-", "-RRB-":
		return TagPUNCT
	case "$", "#":
		return TagSYM
	case "FW", "LS":
		return TagX
	default:
		return 0
	}
}

// TagFromBrownTag converts a single Brown tag to the reduce tag set
func TagFromBrownTag(s string) Tag {

	// if it's a foriegn word, just don't use part of speech
	if strings.HasPrefix(s, "FW-") {
		s = "FW"
	}
	// some have multiple suffixes... just remove it twice.
	s = strings.TrimSuffix(s, "-TL")
	s = strings.TrimSuffix(s, "-HL")
	s = strings.TrimSuffix(s, "-HC")
	s = strings.TrimSuffix(s, "-NC")

	s = strings.TrimSuffix(s, "-TL")
	s = strings.TrimSuffix(s, "-HL")
	s = strings.TrimSuffix(s, "-HC")
	s = strings.TrimSuffix(s, "-NC")

	switch s {
	case "NN", "NNS", "NR", "NR$", "NRS", "NN$", "NNS$":
		return TagNOUN
	case "NP", "NPS", "NP$", "NPS$":
		return TagPROPN
	case "VB", "VBD", "VBG", "VBN", "VBZ":
		return TagVERB
	case "BE", "BED", "BED*", "BEDZ", "BEDZ*", "BEG", "BEM", "BEN", "BER", "BER*", "BEZ", "BEZ*",
		"DO", "DO*", "DOD", "DOD*", "DOZ", "DOZ*",
		"HV", "HV*", "HVD", "HVD*", "HVG", "HVN", "HVZ", "HVZ*",
		"MD", "MD*":
		return TagAUX
	case "JJ", "JJ$", "JJR", "JJS", "JJT":
		return TagADJ
	case "RB", "RB$", "RBR", "RBT", "RN", "WRB",
		"QL", "QLP", "WQL":
		return TagADV
	case "PN", "PN$", "PPSS", "PPO", "PP", "PP$", "PPS", "PPL", "PPLS", "WP$", "WP", "WPO", "WPS", "PP$$":
		return TagPRON
	case "AT", "DT", "DT$", "DTI", "DTS", "DTX", "WDT", "PDT",
		"ABX", "AP", "AP$":
		return TagDET
	case "CD", "OD", "CD$":
		return TagNUM
	case "IN", "TO":
		return TagADP
	case "CC":
		return TagCCONJ
	case "CS":
		return TagSCONJ
	case "ABL", "ABN", "RP":
		return TagPART
	case "EX":
		return TagPART
	case "FW":
		return TagX
	case "UH":
		return TagINTJ
	case "(", ")", ".", ":", ",", "'", "--":
		return TagPUNCT
	case "*":
		return TagSYM
	case "", "NIL", "''", "``":
		return TagUNK // SKIP
	default:
		panic("Unknown tag: " + s)
	}
}
