package blacklist

import (
	"encoding/csv"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
	tier0 []*regexp.Regexp
	tier1 []*regexp.Regexp
}

var (
	instance *Blacklist
	once     sync.Once
)

// blacklist is CSV style array w/columns [phrase, tier (0 | 1), type ("regex" | "string")]
func GetBlackList() *Blacklist {
	once.Do(func() {
		entries, _ := readCSV("config/blacklist.csv")
		instance = buildBlacklist(entries)
	})

	return instance
}

func buildBlacklist(entries [][]string) *Blacklist {
	blacklist := Blacklist{
		tier0: []*regexp.Regexp{},
		tier1: []*regexp.Regexp{},
	}

	// string entries are evaluated as one large regular expression
	tier0Strings := []string{}
	tier1Strings := []string{}

	for i, entry := range entries {
		checkEntry(i, entry)

		content := entry[0]
		tier, _ := strconv.Atoi(entry[1])
		mtype := MatchType(entry[2])

		if mtype == RegexType {
			re := regexp.MustCompile(entry[0])
			if tier == 0 {
				blacklist.tier0 = append(blacklist.tier0, re)
			} else if tier == 1 {
				blacklist.tier1 = append(blacklist.tier1, re)
			}
		} else if mtype == StringType {
			if tier == 0 {
				tier0Strings = append(tier0Strings, content)
			} else if tier == 1 {
				tier1Strings = append(tier1Strings, content)
			}
		}
	}

	tier0StringsReg := regexp.MustCompile("(?i)\\b(" + strings.Join(tier0Strings, "|") + ")\\b")
	tier1StringsReg := regexp.MustCompile("(?i)\\b(" + strings.Join(tier1Strings, "|") + ")\\b")

	blacklist.tier0 = append([]*regexp.Regexp{tier0StringsReg}, blacklist.tier0...)
	blacklist.tier1 = append([]*regexp.Regexp{tier1StringsReg}, blacklist.tier1...)

	return &blacklist
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

func readCSV(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	entries, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return entries, nil
}
