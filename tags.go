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
