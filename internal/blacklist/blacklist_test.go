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
