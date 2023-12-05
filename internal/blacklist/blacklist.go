package blacklist

import (
	"regexp"
	"strconv"
	"strings"
)

type MatchType string
type TierType int

const (
	RegexType  MatchType = "regex"
	StringType MatchType = "string"

	Tier0 TierType = 0
	Tier1 TierType = 1
)

type Blacklist struct {
	Tier0 []*regexp.Regexp
	Tier1 []*regexp.Regexp
}

// blacklist is CSV style array w/columns [phrase, tier (0 | 1), type ("regex" | "string")]
func MakeBlacklist(entries [][]string) *Blacklist {
	bl := Blacklist{
		Tier0: []*regexp.Regexp{},
		Tier1: []*regexp.Regexp{},
	}

	// string entries are evaluated as one large regular expression
	var tier0Builder, tier1Builder strings.Builder

	for i, entry := range entries {
		// skip header row
		if i == 0 {
			continue
		}
		checkEntry(i, entry)

		content := entry[0]
		tier, _ := strconv.Atoi(entry[1])
		mtype := MatchType(entry[2])

		if mtype == RegexType {
			re := regexp.MustCompile(entry[0])
			if tier == 0 {
				bl.Tier0 = append(bl.Tier0, re)
			} else if tier == 1 {
				bl.Tier1 = append(bl.Tier1, re)
			}
		} else if mtype == StringType {
			if tier == 0 {
				if tier0Builder.Len() > 0 {
					tier0Builder.WriteString("|")
				}
				tier0Builder.WriteString(regexp.QuoteMeta(content))
			} else if tier == 1 {
				if tier1Builder.Len() > 0 {
					tier1Builder.WriteString("|")
				}
				tier1Builder.WriteString(regexp.QuoteMeta(content))
			}
		}
	}

	tier0StringsReg := regexp.MustCompile("(?i)\\b(" + tier0Builder.String() + ")\\b")
	tier1StringsReg := regexp.MustCompile("(?i)\\b(" + tier1Builder.String() + ")\\b")

	bl.Tier0 = append([]*regexp.Regexp{tier0StringsReg}, bl.Tier0...)
	bl.Tier1 = append([]*regexp.Regexp{tier1StringsReg}, bl.Tier1...)

	return &bl
}

func checkEntry(row int, entry []string) {
	if len(entry) != 3 {
		rowPanic(row, "Does not contain 3 columns")
	}
	if len(entry[0]) < 1 {
		rowPanic(row, "Content is empty")
	}
	tier, err := strconv.Atoi(entry[1])
	if err != nil || (tier != 0 && tier != 1) {
		rowPanic(row, "Tier must be 0 or 1, received: "+entry[1])
	}
	mtype := MatchType(entry[2])
	if mtype != RegexType && mtype != StringType {
		rowPanic(row, "Type must be \"string\" or \"regex\", received: "+entry[2])
	}
}

func rowPanic(row int, context string) {
	panic("Blacklist malformed at row " + string(rune(row+1)) + ". " + context)
}

// true if message matches any tier 0 blacklist entry
func (bl *Blacklist) EvalTier0(msg string) bool {
	isMatched := false

	for _, re := range bl.Tier0 {
		if re.MatchString(msg) {
			isMatched = true
			break
		}
	}

	return isMatched
}

// true if message matches any tier 1 blacklist entry
func (bl *Blacklist) EvalTier1(msg string) bool {
	isMatched := false

	for _, re := range bl.Tier1 {
		if re.MatchString(msg) {
			isMatched = true
			break
		}
	}

	return isMatched
}
