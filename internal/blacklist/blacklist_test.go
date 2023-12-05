package blacklist

import (
	"testing"
)

func BenchmarkCheckTier0(b *testing.B) {
	// Setup code (e.g., create an instance of Blacklist with some data)
	bl := GetBlackList()

	// Run the function b.N times
	for n := 0; n < b.N; n++ {
		bl.CheckTier0("some test message")
	}
}
