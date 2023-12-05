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
	Tier0 []*regexp.Regexp
	Tier1 []*regexp.Regexp
}

var (
	instance *Blacklist
	once     sync.Once
)

/****************************************************
 Get Blacklist
****************************************************/

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
				blacklist.Tier0 = append(blacklist.Tier0, re)
			} else if tier == 1 {
				blacklist.Tier1 = append(blacklist.Tier1, re)
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

	blacklist.Tier0 = append([]*regexp.Regexp{tier0StringsReg}, blacklist.Tier0...)
	blacklist.Tier1 = append([]*regexp.Regexp{tier1StringsReg}, blacklist.Tier1...)

	return &blacklist
}

// func buildBlacklist0(entries [][]string) *Blacklist {
// 	blacklist := Blacklist{
// 		Tier0: []*regexp.Regexp{},
// 		Tier1: []*regexp.Regexp{},
// 	}

// 	// string entries are evaluated as one large regular expression
// 	var tier0ExactBuilder, tier1ExactBuilder strings.Builder

// 	tier0Strings := []string{}
// 	tier1Strings := []string{}

// 	for i, entry := range entries {
// 		// skip header row
// 		if i == 0 {
// 			continue
// 		}
// 		checkEntry(i, entry)

// 		content := entry[0]
// 		tier, _ := strconv.Atoi(entry[1])
// 		mtype := MatchType(entry[2])

// 		if mtype == RegexType {
// 			re := regexp.MustCompile(entry[0])
// 			if tier == 0 {
// 				blacklist.Tier0 = append(blacklist.Tier0, re)
// 			} else if tier == 1 {
// 				blacklist.Tier1 = append(blacklist.Tier1, re)
// 			}
// 		} else if mtype == StringType {
// 			if tier == 0 {
// 				tier0Strings = append(tier0Strings, content)
// 			} else if tier == 1 {
// 				tier1Strings = append(tier1Strings, content)
// 			}
// 		}
// 	}

// 	tier0StringsReg := regexp.MustCompile("(?i)\\b(" + strings.Join(tier0Strings, "|") + ")\\b")
// 	tier1StringsReg := regexp.MustCompile("(?i)\\b(" + strings.Join(tier1Strings, "|") + ")\\b")

// 	blacklist.Tier0 = append([]*regexp.Regexp{tier0StringsReg}, blacklist.Tier0...)
// 	blacklist.Tier1 = append([]*regexp.Regexp{tier1StringsReg}, blacklist.Tier1...)

// 	return &blacklist
// }

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

/** Review Message ******************************/
func (bl *Blacklist) SyncCheckTier0(msg string) bool {
	isOK := true

	for _, re := range bl.Tier0 {
		if re.MatchString(msg) {
			isOK = false
			break
		}
	}

	return isOK
}

func (bl *Blacklist) CheckTier1(msg string) bool {
	isOK := false

	for _, re := range bl.Tier1 {
		if re.MatchString(msg) {
			isOK = true
			break
		}
	}

	return isOK
}
