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
