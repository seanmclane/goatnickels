package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/seanmclane/goatnickels/rpc"
	"golang.org/x/crypto/sha3"
	"math/big"
	"net/http"
	"sort"
	"strconv"
)

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   int    `json:"amount"`
	Sequence int    `json:"sequence"`
	R        string `json:"r"`
	S        string `json:"s"`
}

//create candidate set of transactions
var CandidateSet []Transaction

//create a staging set of transactions during voting
var StagingCandidateSet []Transaction

func HashCandidateSet(cs *[]Transaction) (h []byte) {
	var sum string
	for _, txion := range *cs {
		sum += hex.EncodeToString(txion.HashTransaction())
	}
	fixedHash := sha3.Sum512([]byte(sum))
	h = fixedHash[:]
	return h
}

func ResetCandidateSet() {
	CandidateSet = []Transaction{}
	for _, t := range StagingCandidateSet {
		CandidateSet = append(CandidateSet, t)
	}
	StagingCandidateSet = []Transaction{}
}

func (t *Transaction) AddTransaction() (ok bool) {
	ok = t.VerifyTransaction()
	if ok != true {
		return false
	}
	//check if transaction exists in candidate set or staging set
	//if so, return and don't broadcast
	for _, c := range CandidateSet {
		if *t == c {
			return false
		}
	}
	for _, c := range StagingCandidateSet {
		if *t == c {
			return false
		}
	}
	//add to candidate set if not currently voting on transactions
	//otherwise add to staging set
	if Voting {
		StagingCandidateSet = append(StagingCandidateSet, *t)
	} else {
		CandidateSet = append(CandidateSet, *t)
	}

	t.Broadcast()
	return ok
}

func (d *Data) ApplyTransactions() {

	SortTransactions()

	//add and subtract from accounts
	for _, txion := range CandidateSet {
		if txion.From == "mine" {
			tnb := d.State[txion.To].Balance + txion.Amount
			d.State[txion.To] = Account{
				Balance:  tnb,
				Sequence: d.State[txion.To].Sequence,
			}
			continue
		}
		//check if sequence is incremented by one
		ok := txion.VerifySequence()
		if ok == true {
			//increment account sequence and debit balance
			fnb := d.State[txion.From].Balance - txion.Amount
			d.State[txion.From] = Account{
				Balance:  fnb,
				Sequence: txion.Sequence,
			}
			//credit balance
			tnb := d.State[txion.To].Balance + txion.Amount
			d.State[txion.To] = Account{
				Balance:  tnb,
				Sequence: d.State[txion.To].Sequence,
			}
			//add valid transaction to transaction list
			d.Transactions = append(d.Transactions, txion)
		} else {
			//TODO: else mark transaction as failed and do something with it? broadcast to network?

		}
	}

	//reset candidate transactions in goatnickels.go

}

func (t *Transaction) SignTransaction(privateKey string) (r, s string) {
	hash := t.HashTransaction()
	//recreate ecdsa.PrivateKey from private_key
	byteKey, err := hex.DecodeString(privateKey)
	if err != nil {
		fmt.Println("error:", err)
	}
	bigintKey := new(big.Int).SetBytes(byteKey)
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = elliptic.P384()
	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(byteKey)
	priv.D = bigintKey
	rInt, sInt, err := ecdsa.Sign(rand.Reader, priv, hash)
	if err != nil {
		fmt.Println("error:", err)
	}
	r = hex.EncodeToString(rInt.Bytes())
	s = hex.EncodeToString(sInt.Bytes())
	return r, s
}

func (t *Transaction) HashTransaction() (h []byte) {
	hashString := t.To + t.From + strconv.Itoa(t.Amount) + strconv.Itoa(t.Sequence)
	fixedHash := sha512.Sum512([]byte(hashString))
	h = fixedHash[:]
	return h
}

func (t *Transaction) VerifyTransaction() (ok bool) {
	//check if t.R and t.S ok with public key
	//what is being signed exactly? hash of transaction sequence, to, from, and amount
	hash := t.HashTransaction()
	//public_key := "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f"
	//check that key is well formed
	if len(t.From) < 100 {
		return false
	}
	//remove goat_ from key
	publicKey := t.From[5:]
	//recreate ecdsa.PublicKey from pub
	byteKey, err := hex.DecodeString(publicKey)
	if err != nil {
		fmt.Println("error:", err)
	}
	x, y := elliptic.Unmarshal(elliptic.P384(), byteKey)
	if x == nil || y == nil {
		return false
	}
	pub := new(ecdsa.PublicKey)
	pub.Curve = elliptic.P384()
	pub.X, pub.Y = x, y
	//convert r and s back to big ints
	byteR, err := hex.DecodeString(t.R)
	if err != nil {
		fmt.Println("error:", err)
	}
	r := new(big.Int).SetBytes(byteR)
	byteS, err := hex.DecodeString(t.S)
	if err != nil {
		fmt.Println("error:", err)
	}
	s := new(big.Int).SetBytes(byteS)
	//verify signature
	sigOk := ecdsa.Verify(pub, hash, r, s)
	//verify balance is sufficient
	balOk := t.VerifyPositiveBalance()
	//verify transaction not < 0
	spendOk := t.VerifyNegativeSpend()
	//do not check sequence here, so you can have more than one transaction per block
	if sigOk && balOk && spendOk {
		return true
	} else {
		return false
	}
}

func (t *Transaction) VerifyPositiveBalance() (ok bool) {
	if LastGoatBlock.Data.State[t.From].Balance < t.Amount {
		return false
	} else {
		return true
	}
}

func (t *Transaction) VerifyNegativeSpend() (ok bool) {
	if t.Amount > 0 {
		return true
	} else {
		return false
	}
}

func (t *Transaction) VerifySequence() (ok bool) {
	//sequence must be current account sequence number plus one
	if t.Sequence == LastGoatBlock.Data.State[t.From].Sequence+1 {
		return true
	} else {
		return false
	}
}

func (t *Transaction) Broadcast() {

	//convert data to plain json
	out, err := json.Marshal(t)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//broadcast transaction through websocket
	rpc.BroadcastChannel <- rpc.BuildNotification("transaction", out)

	config := LoadConfig()

	//TODO: remove self node from list of nodes

	for _, node := range config.Nodes {
		json, err := json.Marshal(t)
		if err != nil {
			fmt.Println("error:", err)
		}
		req, err := http.NewRequest("POST", "http://"+node+":3000/api/v1/txion", bytes.NewBuffer(json))
		if err != nil {
			fmt.Println("error:", err)
		}
		req.Header.Set("Content-Type", "application/json")
		r, err := client.Do(req)
		if err != nil {
			fmt.Println("error:", err)
		} else {
			fmt.Println("Status:", r.Status)
		}
	}
}

//begin accessory functions and types for multisorting transactions
type lessFunc func(p1, p2 *Transaction) bool
type multiSorter struct {
	transactions []Transaction
	less         []lessFunc
}

func (ms *multiSorter) Sort(transactions []Transaction) {
	ms.transactions = transactions
	sort.Sort(ms)
}
func (ms *multiSorter) Len() int {
	return len(ms.transactions)
}
func (ms *multiSorter) Swap(i, j int) {
	ms.transactions[i], ms.transactions[j] = ms.transactions[j], ms.transactions[i]
}
func (ms *multiSorter) Less(i, j int) bool {
	p, q := &ms.transactions[i], &ms.transactions[j]
	var k int
	for k = 0; k < len(ms.less)-1; k++ {
		less := ms.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return ms.less[k](p, q)
}
func OrderedBy(less ...lessFunc) *multiSorter {
	return &multiSorter{
		less: less,
	}
}

//end accessory functions and types for multisorting transactions

func SortTransactions() {
	//order the candidate set to apply deterministically

	fromSort := func(t1, t2 *Transaction) bool {
		return t1.From < t2.From
	}
	amountSort := func(t1, t2 *Transaction) bool {
		return t1.Amount > t2.Amount
	}
	sequenceSort := func(t1, t2 *Transaction) bool {
		return t1.Sequence < t2.Sequence
	}

	OrderedBy(fromSort, sequenceSort, amountSort).Sort(CandidateSet)

}
