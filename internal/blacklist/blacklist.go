package blacklist

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
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

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}

	csvPath := getCsvPath()

	csv, err := readCSV(csvPath)
	if err != nil {
		logger.Fatal("Unable to read Blacklist CSV: ", err.Error())
	}

	blacklistEntries, err := csvToEntries(csv)
	if err != nil {
		logger.Fatal("Malformed blacklist CSV: ", err.Error())
	}

	blacklist, err = makeBlacklist(blacklistEntries)
	if err != nil {
		logger.Fatal("Error making blacklist: ", err.Error())
	}

	if err := verifyPolicyDocuments(blacklist); err != nil {
		logger.Fatal("Error varifying policy documents: ", err.Error())
	}

	logger.Info("Successfully initialized blacklist: " + csvPath)
}

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

	sort.Strings(violations)
	return violations
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		logger.Fatal("Unable to open file: " + err.Error())
		return nil, fmt.Errorf("unable to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	data, err := reader.ReadAll()
	if err != nil {
		logger.Fatal("Unable to read CSV data: " + err.Error())
		return nil, fmt.Errorf("unable to open file: %w", err)
	}

	return data, nil
}

func csvToEntries(csv [][]string) ([]CSVBlacklistEntry, error) {
	var entries []CSVBlacklistEntry

	if len(csv) == 0 {
		return nil, fmt.Errorf("csv is empty")
	}

	for i, row := range csv {
		if i == 0 { // Skip header
			continue
		}
		if len(row) != 4 {
			return nil, fmt.Errorf("CSV row %d does not contain exactly 4 columns", i)
		}

		tier, err := strconv.Atoi(row[3])
		if err != nil {
			return nil, fmt.Errorf("invalid tier value at row %d: %w", i, err)
		}

		entry := CSVBlacklistEntry{
			Content:     row[0],
			ContentType: row[1],
			Policy:      strings.ToLower(row[2]),
			Tier:        tier,
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func makeBlacklist(blacklistEntries []CSVBlacklistEntry) (map[string][]*regexp.Regexp, error) {
	regexMap := make(map[string][]*regexp.Regexp)
	stringMap := make(map[string][]string) // temporary map to hold strings for each key

	for _, entry := range blacklistEntries {
		key := fmt.Sprintf("%d-%s", entry.Tier, strings.TrimSpace(entry.Policy))

		contentType := strings.TrimSpace(entry.ContentType)
		if contentType == "regex" {
			re, err := regexp.Compile(entry.Content)
			if err != nil {
				return nil, fmt.Errorf("invalid regex in entry %s: %w", entry.Content, err)
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
			return nil, fmt.Errorf("invalid regex from strings at key %s: %w", key, err)
		}
		regexMap[key] = append(regexMap[key], re)
	}

	return regexMap, nil
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
		return fmt.Errorf("unable to initialize blacklist due to missing policy documents: %s", strings.Join(missingDocs, ", "))
	}

	return nil
}

func getCsvPath() string {
	csvRelPath := os.Getenv("BLACKLIST_CSV_PATH")
	if csvRelPath == "" {
		csvRelPath = "/config/blacklist.csv"
	}

	// retrieve the runtime file path
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		logger.Fatal("Cannot retrieve runtime information")
	}

	// navigate up to the project root from current file (`internal/blacklist/blacklist.go`)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(b)))

	return filepath.Join(projectRoot, csvRelPath)
}
