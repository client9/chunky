// mkrules reads the pipe-delimited output of rules.sh from stdin and writes
// a Go source file containing a sorted ContextRule table for tok/rules_gen.go.
//
// Usage:
//
//	bash rules.sh | go run ./cmd/mkrules -tag1 NOUN -tag2 VERB > tok/rules_gen.go
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/client9/chunky"
)

// featureSlots maps a feature template name to the ordered list of context
// slot names that its key segments fill.
var featureSlots = map[string][]string{
	"nexttag":                           {"next"},
	"prevtag":                           {"prev"},
	"prevtag+nexttag":                   {"prev", "next"},
	"prev2tag+prevtag":                  {"prev2", "prev"},
	"prevtag+nexttag+next2tag":          {"prev", "next", "next2"},
	"prev2tag+prevtag+nexttag":          {"prev2", "prev", "next"},
	"prev2tag+prevtag+nexttag+next2tag": {"prev2", "prev", "next", "next2"},
}

type rule struct {
	tags        chunky.Tag
	prev2, prev chunky.Tag
	next, next2 chunky.Tag
	mask        uint8
	resolve     chunky.Tag
	specificity int
}

const (
	bitPrev2 uint8 = 1 << 3
	bitPrev  uint8 = 1 << 2
	bitNext  uint8 = 1 << 1
	bitNext2 uint8 = 1 << 0
)

func slotMask(slot string) uint8 {
	switch slot {
	case "prev2":
		return bitPrev2
	case "prev":
		return bitPrev
	case "next":
		return bitNext
	case "next2":
		return bitNext2
	}
	return 0
}

// singleTagName returns the Go identifier for a single-bit Tag constant.
func singleTagName(tag chunky.Tag) string {
	if tag == chunky.TagUNK {
		return "TagUNK"
	}
	return "Tag" + tag.String()
}

// tagGoName emits the Go expression for a Tag value, handling SuperTags (ORed bits).
func tagGoName(t chunky.Tag) string {
	if t == 0 {
		return "0"
	}
	var parts []string
	for _, tag := range chunky.AllTags {
		if t&tag != 0 {
			parts = append(parts, singleTagName(tag))
		}
	}
	if len(parts) == 0 {
		return "0"
	}
	return strings.Join(parts, " | ")
}

// canMerge reports whether two rules are identical except in exactly one slot.
// Rules that differ in a slot where either value is TagUNK cannot be merged:
// TagUNK has special boundary semantics in matchSlot that don't compose with |.
// 1-slot rules are not merged: their corpus-frequency ordering is load-bearing.
func canMerge(a, b rule) bool {
	if a.tags != b.tags {
		return false
	}
	if a.mask != b.mask || a.resolve != b.resolve {
		return false
	}
	if bits(a.mask) < 2 {
		return false
	}
	diff := 0
	if a.prev2 != b.prev2 {
		if a.prev2 == chunky.TagUNK || b.prev2 == chunky.TagUNK {
			return false
		}
		diff++
	}
	if a.prev != b.prev {
		if a.prev == chunky.TagUNK || b.prev == chunky.TagUNK {
			return false
		}
		diff++
	}
	if a.next != b.next {
		if a.next == chunky.TagUNK || b.next == chunky.TagUNK {
			return false
		}
		diff++
	}
	if a.next2 != b.next2 {
		if a.next2 == chunky.TagUNK || b.next2 == chunky.TagUNK {
			return false
		}
		diff++
	}
	return diff == 1
}

// mergeRules combines b into a by ORing differing slot values.
func mergeRules(a, b rule) rule {
	a.prev2 |= b.prev2
	a.prev |= b.prev
	a.next |= b.next
	a.next2 |= b.next2
	return a
}

// consolidate iterates until no two rules can be merged, collapsing rules that
// differ in exactly one slot into a single rule with a SuperTag in that slot.
func consolidate(rules []rule) []rule {
	for {
		merged := make([]bool, len(rules))
		var result []rule
		changed := false
		for i := 0; i < len(rules); i++ {
			if merged[i] {
				continue
			}
			for j := i + 1; j < len(rules); j++ {
				if merged[j] {
					continue
				}
				if canMerge(rules[i], rules[j]) {
					rules[i] = mergeRules(rules[i], rules[j])
					merged[j] = true
					changed = true
				}
			}
			result = append(result, rules[i])
		}
		rules = result
		if !changed {
			break
		}
	}
	return rules
}

func main() {
	tag1Flag := flag.String("tag1", "", "first tag in the ambiguous pair (e.g. NOUN)")
	tag2Flag := flag.String("tag2", "", "second tag in the ambiguous pair (e.g. VERB)")
	pkgFlag := flag.String("pkg", "tok", "Go package name for the output file")
	varFlag := flag.String("var", "contextRules", "name of the generated []ContextRule variable")
	noHeaderFlag := flag.Bool("noheader", false, "suppress the package/generated header (for appending to a combined file)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: rules.sh | mkrules -tag1 TAG -tag2 TAG [-pkg PKG] [-var NAME] [-noheader]")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *tag1Flag == "" || *tag2Flag == "" {
		flag.Usage()
		os.Exit(1)
	}

	tag1, err := chunky.ParseTag(*tag1Flag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tag2, err := chunky.ParseTag(*tag2Flag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var rules []rule
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 5 {
			fmt.Fprintf(os.Stderr, "skipping malformed line: %q\n", line)
			continue
		}
		feature := parts[1]
		key := parts[2]
		resolveStr := parts[3]

		slots, ok := featureSlots[feature]
		if !ok {
			fmt.Fprintf(os.Stderr, "unknown feature %q, skipping\n", feature)
			continue
		}

		resolve, err := chunky.ParseTag(resolveStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad resolve tag %q: %v\n", resolveStr, err)
			continue
		}
		if resolve != tag1 && resolve != tag2 {
			fmt.Fprintf(os.Stderr, "resolve %q is not %s or %s, skipping\n", resolveStr, tag1, tag2)
			continue
		}

		keyParts := strings.Split(key, "+")
		if len(keyParts) != len(slots) {
			fmt.Fprintf(os.Stderr, "key %q has %d parts, feature %q needs %d, skipping\n",
				key, len(keyParts), feature, len(slots))
			continue
		}

		r := rule{tags: tag1 | tag2, resolve: resolve}
		for i, slot := range slots {
			part := keyParts[i]
			r.mask |= slotMask(slot)
			var t chunky.Tag
			if part != "" {
				t, err = chunky.ParseTag(part)
				if err != nil {
					fmt.Fprintf(os.Stderr, "bad tag %q in key %q: %v\n", part, key, err)
					goto skip
				}
			} else {
				t = chunky.TagUNK // empty part = sentence boundary
			}
			switch slot {
			case "prev2":
				r.prev2 = t
			case "prev":
				r.prev = t
			case "next":
				r.next = t
			case "next2":
				r.next2 = t
			}
		}

		r.specificity = bits(r.mask)
		rules = append(rules, r)
		continue
	skip:
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	before := len(rules)
	rules = consolidate(rules)
	after := len(rules)

	// Sort most-specific-first so earlier rules win on first match.
	sort.SliceStable(rules, func(i, j int) bool {
		return rules[i].specificity > rules[j].specificity
	})

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	if !*noHeaderFlag {
		fmt.Fprintf(w, "// Code generated by cmd/mkrules. DO NOT EDIT.\n")
		fmt.Fprintf(w, "package %s\n\n", *pkgFlag)
	}
	fmt.Fprintf(w, "// %s is the compiled context disambiguation rule table for %s vs %s.\n", *varFlag, tag1, tag2)
	fmt.Fprintf(w, "// Sorted most-specific-first (most active context slots first).\n")
	fmt.Fprintf(w, "// %d rules consolidated from %d (%d merged).\n", after, before, before-after)
	fmt.Fprintf(w, "// Generated by: F1=%s F2=%s bash rules.sh | go run ./cmd/mkrules -tag1 %s -tag2 %s -var %s\n", tag1, tag2, tag1, tag2, *varFlag)
	fmt.Fprintf(w, "var %s = []ContextRule{\n", *varFlag)
	for _, r := range rules {
		fmt.Fprintf(w, "\t{Tags: %s, Prev2: %s, Prev: %s, Next: %s, Next2: %s, Mask: 0x%02x, Resolve: %s},\n",
			tagGoName(r.tags),
			tagGoName(r.prev2), tagGoName(r.prev),
			tagGoName(r.next), tagGoName(r.next2),
			r.mask, tagGoName(r.resolve),
		)
	}
	fmt.Fprintf(w, "}\n")
}

func bits(b uint8) int {
	n := 0
	for b != 0 {
		n += int(b & 1)
		b >>= 1
	}
	return n
}
