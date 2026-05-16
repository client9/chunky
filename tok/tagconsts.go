package tok

import "github.com/client9/chunky"

// Tag is an alias for chunky.Tag so tok-package files need not import chunky directly.
type Tag = chunky.Tag

// Tag constants re-exported from chunky so tok-package files need not import chunky directly.
const (
	TagUNK   = chunky.TagUNK
	TagADJ   = chunky.TagADJ
	TagADP   = chunky.TagADP
	TagADV   = chunky.TagADV
	TagAUX   = chunky.TagAUX
	TagCCONJ = chunky.TagCCONJ
	TagDET   = chunky.TagDET
	TagINTJ  = chunky.TagINTJ
	TagNOUN  = chunky.TagNOUN
	TagNUM   = chunky.TagNUM
	TagPART  = chunky.TagPART
	TagPRON  = chunky.TagPRON
	TagPROPN = chunky.TagPROPN
	TagPUNCT = chunky.TagPUNCT
	TagSCONJ = chunky.TagSCONJ
	TagSYM   = chunky.TagSYM
	TagVERB  = chunky.TagVERB
	TagX     = chunky.TagX
)
