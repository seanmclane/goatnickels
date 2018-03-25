package block

import (
	"testing"
)

func BenchmarkCreateGenesisBlock(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CreateGenesisBlock()
	}
}

func TestSortTransactions(t *testing.T) {
	CandidateSet = []Transaction{
		{"a", "b", 100, 2, "", ""},
		{"b", "a", 50, 1, "", ""},
		{"a", "b", 100, 1, "", ""},
		{"a", "b", 101, 1, "", ""},
		{"b", "a", 100, 2, "", ""},
	}

	SortTransactions()

	var expected = []Transaction{
		{"a", "b", 101, 1, "", ""},
		{"a", "b", 100, 1, "", ""},
		{"a", "b", 100, 2, "", ""},
		{"b", "a", 50, 1, "", ""},
		{"b", "a", 100, 2, "", ""},
	}

	for key, cs := range CandidateSet {
		if cs.From != expected[key].From {
			t.Errorf("From sort failed: expected %v, actual %v", expected[key].From, cs.From)
		}
		if cs.Amount != expected[key].Amount {
			t.Errorf("Amount sort failed: expected %v, actual %v", expected[key].Amount, cs.Amount)
		}
		if cs.Sequence != expected[key].Sequence {
			t.Errorf("Sequence sort failed: expected %v, actual %v", expected[key].Sequence, cs.Sequence)
		}
	}

}
