package block

import (
	"encoding/hex"
	"encoding/json"
	//	"fmt"
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

func TestVerifyBlock(t *testing.T) {
	//based on testnet genesis block
	j1 := []byte("{\"index\":2,\"timestamp\":1522011630,\"data\":{\"state\":{\"goat_045b4dfabe49048ef6fb6e47fc4e2b33dd54e46b3ed4ab008f8dce7457f588f7a6975690328db4bd48eb874ff909c579fe37ae4f39e9b9b10ac1f2f49083c7d2d8fe91ff5314b2742d58e894681d55682876417f33f851e8091f9c00045a7a9ebc\":{\"balance\":76457654265,\"sequence\":0},\"goat_04ab1594a3b65e440653b1a54952aee3cb7f5c41cb476f7ecd3ce58dc23cef0923beb45fc275ff4149cd9f0417f8ca885e882b3b68d00bab2988b22f2eaf7f6683ba3e672abd668e5788a8ecb4d055cd024f004ff03db06158f18e5bd02914685a\":{\"balance\":94043534214,\"sequence\":0},\"goat_04c7cb2cef7da5cda83333f34fba7f07b3d1a7572ca909487c7ed20d147706b731e26983c18659bc1caf260a4fd4fc390d9bec208c92d123498faad57ae365ba3aebcd4a93e74802adee03cfbac8f71ed7f5d00824de59bf292c20b2b73bd3228d\":{\"balance\":38763423645,\"sequence\":0},\"goat_04dbb67ae9650ca3258071909f74be5400fe53fc2e5dcc82103020f3aeefeee5f9980c4c05bb8696215458dfa7ddaa1505d2826cab3d246b8930b0694f766a22f8bb63932368c0b12bf80cfaee8a18db1d7ce19df0a84215d20b0bbfbd30d95c25\":{\"balance\":50984323425,\"sequence\":0}},\"transactions\":null},\"last_hash\":\"fe1e2c67bceaf9e71c11f1816e8693873a1647998954f91465fc02aba9003812793e2e2eeefa1696e46e18501bbd07ada4dbdc0f51e622d18cb6aaa92b031332\",\"hash\":\"5ee36adee27eb5d600d464a7e50829e439c015f777e1619fd7ad36a633d7deafd9fdff5d18e21fa575027d56f018af12d0ede8470027ecb04c53b5755686758c\"}")

	var s1 StoredBlock
	err := json.Unmarshal(j1, &s1)
	if err != nil {
		t.Errorf("Failed to read block")
	}

	lastHash, _ := hex.DecodeString(s1.LastHash)
	hash, _ := hex.DecodeString(s1.Hash)

	b1 := Block{
		Index:     s1.Index,
		Timestamp: s1.Timestamp,
		Data:      s1.Data,
		LastHash:  lastHash,
		Hash:      hash,
	}

	b2 := b1
	b2.Data.State = make(map[string]Account)

	var blocks = []Block{b1, b2}

	var expected = []bool{true, false}

	for key, b := range blocks {
		ok := b.VerifyBlock()
		if ok != expected[key] {
			t.Errorf("Failed verification of block hashes: %v was %v", hex.EncodeToString(b.LastHash), ok)
		}
	}
}
