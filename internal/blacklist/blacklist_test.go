package blacklist

import (
	"os"
	"testing"
)

const (
	csvFilePath = "/tmp/mock_blacklist.csv"
)

func TestMain(m *testing.M) {
	setupMockCSV()
	defer tearDownMockCSV() // Clean up after all tests have run

	os.Exit(m.Run()) // Run the tests
}

func TestReadCSV(t *testing.T) {
	result, err := readCSV(csvFilePath)
	if err != nil {
		t.Fatalf("readCSV failed with error: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("Expected 2 rows, got %d", len(result))
	}
}

func TestCsvToEntries(t *testing.T) {
	cases := []struct {
		name     string
		input    [][]string
		expected int
		wantErr  bool
	}{
		{
			name: "valid input",
			input: [][]string{
				{"content", "contentType", "policy", "tier"},
				{"test", "regex", "profanity", "0"},
				{"hello", "string", "lead-generation", "1"},
			},
			expected: 2,
			wantErr:  false,
		},
		{
			name:     "empty csv",
			input:    [][]string{},
			expected: 0,
			wantErr:  true,
		},
		{
			name: "incorrect columns",
			input: [][]string{
				{"content", "contentType", "policy", "tier"},
				{"invalid", "row"},
			},
			expected: 0,
			wantErr:  true,
		},
		{
			name: "invalid tier value",
			input: [][]string{
				{"content", "contentType", "policy", "tier"},
				{"test", "regex", "profanity", "invalid"},
			},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := csvToEntries(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("csvToEntries() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if len(result) != tc.expected {
				t.Errorf("Expected %d entries, got %d", tc.expected, len(result))
			}
		})
	}
}

func TestMakeBlacklist(t *testing.T) {
	// Define test cases
	cases := []struct {
		name         string
		entries      []CSVBlacklistEntry
		expectedKeys []string
		wantErr      bool
	}{
		{
			name: "valid entries",
			entries: []CSVBlacklistEntry{
				{Content: "badword", ContentType: "string", Policy: "profanity", Tier: 0},
				{Content: "^badword$", ContentType: "regex", Policy: "profanity", Tier: 0},
				{Content: "example", ContentType: "string", Policy: "advertising", Tier: 1},
			},
			expectedKeys: []string{"0-profanity", "1-advertising"},
			wantErr:      false,
		},
		{
			name: "invalid regex",
			entries: []CSVBlacklistEntry{
				{Content: "invalid(regex", ContentType: "regex", Policy: "profanity", Tier: 0},
			},
			expectedKeys: nil,
			wantErr:      true,
		},
	}

	// Execute each test case
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := makeBlacklist(tc.entries)

			// Check for error handling
			if (err != nil) != tc.wantErr {
				t.Errorf("makeBlacklist() error = %v, wantErr %v", err, tc.wantErr)
			}

			// If no error expected, validate the content
			if !tc.wantErr {
				if len(result) != len(tc.expectedKeys) {
					t.Errorf("Expected map size %d, got %d", len(tc.expectedKeys), len(result))
				}

				// Check if all expected keys are in the result
				for _, key := range tc.expectedKeys {
					if _, exists := result[key]; !exists {
						t.Errorf("Expected key %s to exist", key)
					}
				}

				// Optionally, test the regex matching functionality
				if key, exists := result["0-profanity"]; exists {
					testString := "badword"
					if !key[0].MatchString(testString) {
						t.Errorf("Regex should match the string %s", testString)
					}
				}
			}
		})
	}
}

func setupMockCSV() {
	csvContent := "content,contentType,policy,tier\n" +
		"test,regex,profanity,0\n" +
		"hello,string,lead-generation,1"

	err := os.WriteFile(csvFilePath, []byte(csvContent), 0644)

	if err != nil {
		panic("Failed to create mock CSV file: " + err.Error())
	}

}
func tearDownMockCSV() {
	os.Remove(csvFilePath)
}
