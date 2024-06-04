package blacklist

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type CSVBlacklistEntry struct {
	Content     string // the text or regex pattern to match against
	ContentType string // "regex" for regular expressions, "string" for direct string matches
	Policy      string // a descriptor of the policy the entry enforces
	Tier        int    // importance level: 0 for auto-rejection, 1 for LLM evaluation
}

// blacklist is a map where each key is a combination of tier and policy (e.g. "0-profanity"),
// and each value is a list of compiled regex patterns for that category.
var blacklist map[string][]*regexp.Regexp

// CheckContent checks the specified content against the compiled blacklist and returns an array
// of policies matched by the content. Policy matches are formatted like so: tier-policy,
// e.g. ["0-profanity"].
func CheckContent(content string) []string {
	var violations []string

	for category, regexList := range blacklist {
		for _, re := range regexList {
			if re.MatchString(content) {
				violations = append(violations, category)
				break
			}
		}
	}

	return violations
}

// Init initializes the blacklist from a CSV file at the specified absolute path.
func Init(absoluteFilePath string) {
	blacklist = buildBlacklist(absoluteFilePath)
	log.Println("Initalized blacklist: " + absoluteFilePath)
}

func buildBlacklist(absoluteFilePath string) map[string][]*regexp.Regexp {
	csv := readCSV(absoluteFilePath)
	blacklistEntries := csvToEntries(csv)

	regexMap := make(map[string][]*regexp.Regexp)
	stringMap := make(map[string][]string) // temporary map to hold strings for each key

	for _, entry := range blacklistEntries {
		key := fmt.Sprintf("%d-%s", entry.Tier, entry.Policy)

		if entry.ContentType == "regex" {
			re, err := regexp.Compile(entry.Content)
			if err != nil {
				panic("Invalid regex: " + entry.Content)
			}
			regexMap[key] = append(regexMap[key], re)
		} else if entry.ContentType == "string" {
			// collect strings to compile into a single regex later
			stringMap[key] = append(stringMap[key], regexp.QuoteMeta(entry.Content))
		}
	}

	// compile all strings into a single regex for each key
	for key, strs := range stringMap {
		pattern := "(?i)(" + strings.Join(strs, "|") + ")"
		re, err := regexp.Compile(pattern)
		if err != nil {
			panic("Invalid regex from strings: " + strings.Join(strs, ", "))
		}
		regexMap[key] = append(regexMap[key], re)
	}

	return regexMap
}

func readCSV(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Unable to open file: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	data, err := reader.ReadAll()
	if err != nil {
		panic("Unable to read CSV data: " + err.Error())
	}

	return data
}

func csvToEntries(csv [][]string) []CSVBlacklistEntry {
	var entries []CSVBlacklistEntry

	for i, row := range csv {
		if i == 0 { // Skip header
			continue
		}
		if len(row) != 4 {
			panic("CSV row does not contain exactly 4 columns.")
		}

		tier, err := strconv.Atoi(row[3])
		if err != nil {
			panic("Invalid tier value: " + row[3])
		}

		entry := CSVBlacklistEntry{
			Content:     row[0],
			ContentType: row[1],
			Policy:      strings.ToLower(row[2]),
			Tier:        tier,
		}
		entries = append(entries, entry)
	}

	return entries
}
