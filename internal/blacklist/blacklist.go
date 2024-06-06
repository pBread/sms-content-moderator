package blacklist

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/pBread/sms-content-moderator/internal/logger"
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
	if err := verifyPolicyDocuments(blacklist); err != nil {
		logger.Fatal("Policy verification failed: " + err.Error())
	}

	logger.Info("Initalized blacklist: " + absoluteFilePath)
}

func buildBlacklist(absoluteFilePath string) map[string][]*regexp.Regexp {
	csv := readCSV(absoluteFilePath)
	blacklistEntries := csvToEntries(csv)

	regexMap := make(map[string][]*regexp.Regexp)
	stringMap := make(map[string][]string) // temporary map to hold strings for each key

	for _, entry := range blacklistEntries {
		key := fmt.Sprintf("%d-%s", entry.Tier, strings.TrimSpace(entry.Policy))

		contentType := strings.TrimSpace(entry.ContentType)
		if contentType == "regex" {
			re, err := regexp.Compile(entry.Content)
			if err != nil {
				logger.Fatal("Invalid regex: " + entry.Content)
			}
			regexMap[key] = append(regexMap[key], re)
		} else if contentType == "string" {
			// collect strings to compile into a single regex later
			stringMap[key] = append(stringMap[key], regexp.QuoteMeta(entry.Content))
		}
	}

	// compile all strings into a single regex for each key
	for key, strs := range stringMap {
		pattern := "(?i)(" + strings.Join(strs, "|") + ")"
		re, err := regexp.Compile(pattern)
		if err != nil {
			logger.Fatal("Invalid regex from strings: " + strings.Join(strs, ", "))
		}
		regexMap[key] = append(regexMap[key], re)
	}

	return regexMap
}

func readCSV(filePath string) [][]string {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal("Unable to open file: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	data, err := reader.ReadAll()
	if err != nil {
		logger.Fatal("Unable to read CSV data: " + err.Error())
	}

	return data
}

func csvToEntries(csv [][]string) []CSVBlacklistEntry {
	var entries []CSVBlacklistEntry

	for i, row := range csv {
		if i == 0 {
			continue
		}
		if len(row) != 4 {
			logger.Fatal("CSV row does not contain exactly 4 columns.")
		}

		tier, err := strconv.Atoi(row[3])
		if err != nil {
			logger.Fatal("Invalid tier value: " + row[3])
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

func verifyPolicyDocuments(blacklist map[string][]*regexp.Regexp) error {
	missingDocs := []string{}

	for policy := range blacklist {
		policyName := strings.SplitN(policy, "-", 2)[1] // Extract policy name from the key
		policyFilePath := fmt.Sprintf("config/policies/%s.md", strings.TrimSpace(policyName))
		if _, err := os.Stat(policyFilePath); os.IsNotExist(err) {
			missingDocs = append(missingDocs, policyName)
		}
	}

	if len(missingDocs) > 0 {
		return fmt.Errorf("missing policy documents for: %v", missingDocs)
	}
	return nil
}
