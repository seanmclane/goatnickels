package block

import (
	"testing"
)

func BenchmarkCreateGenesisBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateGenesisBlock()
	}
}
