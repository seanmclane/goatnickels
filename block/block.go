package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"os/user"
	"sort"
	"strconv"
	"time"
)

//define config structure
//does this belong here?
type Config struct {
	Directory string   `json:"directory"`
	Nodes     []string `json:"nodes"`
	Account   string   `json:"account"`
}

//load config
func LoadConfig() (config Config) {
	user, _ := user.Current()
	dir := user.HomeDir
	c, err := os.Open(dir + "/.goatnickels/config.json")
	if err != nil {
		panic(err)
	}

	//fix this to json unmarshal
	decoder := json.NewDecoder(c)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println("error:", err)
	}

	return config
}

type KeyStore struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

func LoadKeyStore() (keystore KeyStore) {
	user, _ := user.Current()
	dir := user.HomeDir
	k, err := os.Open(dir + "/.goatnickels/keystore.json")
	if err != nil {
		panic(err)
	}

	//fix this to json unmarshal
	decoder := json.NewDecoder(k)
	err = decoder.Decode(&keystore)
	if err != nil {
		fmt.Println("error:", err)
	}

	return keystore
}

//initializing blockchain objects here for now
//need to have last validated block (index at minimum)
var LastGoatBlock Block

//need to have last state
var Accounts map[string]Account

//create candidate set of transactions
var CandidateSet []Transaction

//create a staging set of transactions during voting
var StagingCandidateSet []Transaction

//create vote set to store votes from network
var VoteSet []Vote

//create bool to capture whether voting
var Voting bool

// create http client for all network requests
var client = &http.Client{
	Timeout: time.Second * 10,
}

//removing blockchain since it's not needed to store whole chain in memory
// type Blockchain []Block
// var GoatChain Blockchain

type Block struct {
	Index     int    `json:"index"`
	Timestamp int    `json:"timestamp"`
	Data      Data   `json:"data"`
	LastHash  []byte `json:"last_hash"`
	Hash      []byte `json:"hash"`
}

type StoredBlock struct {
	Index     int    `json:"index"`
	Timestamp int    `json:"timestamp"`
	Data      Data   `json:"data"`
	LastHash  string `json:"last_hash"`
	Hash      string `json:"hash"`
}

type Data struct {
	State        map[string]Account `json:"state"`
	Transactions []Transaction      `json:"transactions"`
}

type Account struct {
	Balance  int `json:"balance"`
	Sequence int `json:"sequence"`
}

type Transaction struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Amount   int    `json:"amount"`
	Sequence int    `json:"sequence"`
	R        string `json:"r"`
	S        string `json:"s"`
}

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
}

type Vote struct {
	Account   string    `json:"account"`
	Hash      string    `json:"hash"`
	Signature Signature `json:"signature"`
}

func AsciiGoat() {
	a := "\x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2d \x2e \x5f \x2c \x2d \x2d \x2e \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2c \x2d \x60 \x2e \x0a \x28 \x60 \x2d \x2e \x5f \x5f \x20 \x20 \x20 \x20 \x2f \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2f \x0a \x20 \x60 \x2e \x20 \x20 \x20 \x60 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x5f \x5f \x2c \x2d \x2d \x27 \x2d \x2e \x0a \x20 \x20 \x20 \x60 \x2d \x2d \x2f \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2e \x20 \x20 \x5f \x5f \x5f \x5f \x5f \x5f \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x2d \x2e \x20 \x20 \x20 \x20 \x20 \x2c \x6f \x2d \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x2e \x20 \x3b \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x7c \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x60 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x20 \x6f \x20 \x2c \x20 \x20 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x2d \x2d \x27 \x2c \x27 \x2e \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x3b \x3b \x20 \x20 \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x27 \x20 \x20 \x3b \x20 \x20 \x2c \x27 \x20 \x2c \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x2c \x27 \x20 \x20 \x3a \x20 \x20 \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x20 \x5c \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60"
	fmt.Println(a)
}

func (b *Block) HashBlock() {
	//create a hash of all values in the block
	//TODO: handle error
	hashData, _ := json.Marshal(b.Data)
	blockString := strconv.Itoa(b.Index) + string(hashData) + hex.EncodeToString(b.LastHash[:])
	fixedHash := sha3.Sum512([]byte(blockString))
	b.Hash = fixedHash[:]
}

func CreateGenesisBlock() {
	//temporary
	//manually adding accounts for now
	Accounts = make(map[string]Account)

	Accounts["goat_04dbb67ae9650ca3258071909f74be5400fe53fc2e5dcc82103020f3aeefeee5f9980c4c05bb8696215458dfa7ddaa1505d2826cab3d246b8930b0694f766a22f8bb63932368c0b12bf80cfaee8a18db1d7ce19df0a84215d20b0bbfbd30d95c25"] = Account{
		Balance:  50884323425,
		Sequence: 0,
	}
	Accounts["goat_04ab1594a3b65e440653b1a54952aee3cb7f5c41cb476f7ecd3ce58dc23cef0923beb45fc275ff4149cd9f0417f8ca885e882b3b68d00bab2988b22f2eaf7f6683ba3e672abd668e5788a8ecb4d055cd024f004ff03db06158f18e5bd02914685a"] = Account{
		Balance:  94043534214,
		Sequence: 0,
	}
	Accounts["goat_04c7cb2cef7da5cda83333f34fba7f07b3d1a7572ca909487c7ed20d147706b731e26983c18659bc1caf260a4fd4fc390d9bec208c92d123498faad57ae365ba3aebcd4a93e74802adee03cfbac8f71ed7f5d00824de59bf292c20b2b73bd3228d"] = Account{
		Balance:  38763423645,
		Sequence: 0,
	}
	Accounts["goat_045b4dfabe49048ef6fb6e47fc4e2b33dd54e46b3ed4ab008f8dce7457f588f7a6975690328db4bd48eb874ff909c579fe37ae4f39e9b9b10ac1f2f49083c7d2d8fe91ff5314b2742d58e894681d55682876417f33f851e8091f9c00045a7a9ebc"] = Account{
		Balance:  76457654265,
		Sequence: 0,
	}
	//set arbitrary data
	data := Data{
		State:        Accounts,
		Transactions: nil,
	}

	//convert [64]byte to []byte
	fixedHash := sha3.Sum512([]byte("Goatnickels baby!"))
	hash := fixedHash[:]

	//genesis block for now
	b := Block{
		Index:     1,
		Timestamp: 0, //TODO: convert this to the birthdate of GoatNickels
		Data:      data,
		LastHash:  hash,
	}

	b.HashBlock()

	b.WriteBlockToLocalStorage()

}

//TODO: figure out where to keep data structures and have one way imports
type MaxBlockResponse struct {
	MaxBlock int `json:"max_block"`
}

func InitializeState() {

	//initialize voting to whether a voting round on the candidate set is happening
	Voting = false

	maxBlock, nodes := GetMaxBlockNumberFromNetwork()

	//genesis block only created by calling function manually
	//always check the network for max block, then start

	localMax := FindMaxBlock()

	//check if different and get blocks from network if behind
	//TODO: loop through nodes? or retry on failure
	//TODO: loop to recheck that max_block was not updated while catching up to network
	if localMax < maxBlock {
		GetBlockChainFromNetwork(localMax, maxBlock, nodes[0])
		localMax = maxBlock
	}

	LastGoatBlock = ReadBlockFromLocalStorage(strconv.Itoa(localMax))

}

func GetMaxBlockNumberFromNetwork() (maxBlockId int, nodes []string) {

	config := LoadConfig()

	//TODO: remove self node from list of nodes

	//list of nodes with the maximum block id as value
	maxList := make(map[string]int)
	//list of maximum block ids with the count as value
	maxCount := make(map[int]int)

	for _, node := range config.Nodes {
		r, err := client.Get("http://" + node + ":3000/api/v1/maxblock")
		if err != nil {
			fmt.Println("no response from node:", node)
			fmt.Println("error:", err)
		} else {
			defer r.Body.Close()
			var res MaxBlockResponse
			err = json.NewDecoder(r.Body).Decode(&res)
			if err != nil {
				fmt.Println("error:", err)
			}
			maxList[node] = res.MaxBlock
		}
	}

	//count occurences of max block id across nodes
	for _, maxBlockId := range maxList {
		maxCount[maxBlockId] += 1
	}

	//TODO: use large % of network agreement to determine true max block
	//get the max block id with the highest count
	maxBlockCount := 0
	maxBlockId = 0
	for blockId, count := range maxCount {
		if count > maxBlockCount {
			maxBlockCount = count
			maxBlockId = blockId
		}
	}

	//get the nodes that have that block
	for node, blockId := range maxList {
		if blockId == maxBlockId {
			nodes = append(nodes, node)
		}
	}

	return maxBlockId, nodes
}

func GetBlockFromNetwork(blockNumber int, node string) {
	r, err := client.Get("http://" + node + ":3000/api/v1/block/" + strconv.Itoa(int(blockNumber)))
	if err != nil {
		fmt.Println("could not get block from", node)
		fmt.Println("error:", err)
	} else {
		defer r.Body.Close()
		var b Block
		err = json.NewDecoder(r.Body).Decode(&b)
		if err != nil {
			fmt.Println("error:", err)
		}
		if b.VerifyBlock() {
			b.WriteBlockToLocalStorage()
		} else {
			fmt.Println("error: invalid block not written")
			//TODO: return something for the parent to try to get the block again
		}
	}
}

func GetBlockChainFromNetwork(localMax int, networkMax int, node string) {
	for i := localMax + 1; i <= networkMax; i++ {
		GetBlockFromNetwork(i, node)
	}
}

func (b *Block) VerifyBlock() bool {
	prev := ReadBlockFromLocalStorage(strconv.Itoa(b.Index - 1))
	prev.HashBlock()
	//check if last hash is previous block hash
	if !bytes.Equal(prev.Hash, b.LastHash) {
		fmt.Println("not equal to lasthash", hex.EncodeToString(prev.Hash), hex.EncodeToString(b.LastHash))
		return false
	}
	//check if current hash is real hash
	bHash := b.Hash
	b.HashBlock()
	if !bytes.Equal(bHash, b.Hash) {
		fmt.Println("not same block data", hex.EncodeToString(bHash), hex.EncodeToString(b.Hash))
		return false
	}
	return true
}

func ReadBlockFromLocalStorage(index string) (b Block) {
	config := LoadConfig()
	blockJson, err := ioutil.ReadFile(string(config.Directory) + index)
	if err != nil {
		fmt.Println("error:", err)
	}
	//make bytestring to Block
	var s StoredBlock

	err = json.Unmarshal(blockJson, &s)
	if err != nil {
		fmt.Println("error:", err)
	}

	//TODO: handle errors
	lastHash, _ := hex.DecodeString(s.LastHash)
	hash, _ := hex.DecodeString(s.Hash)

	b = Block{
		Index:     s.Index,
		Timestamp: s.Timestamp,
		Data:      s.Data,
		LastHash:  lastHash,
		Hash:      hash,
	}

	return b
}

func FindMaxBlock() (max int) {
	config := LoadConfig()
	files, err := ioutil.ReadDir(config.Directory)
	if err != nil {
		panic(err)
	}

	max = 0
	for _, file := range files {
		cur, err := strconv.ParseInt(file.Name(), 10, 0)
		if err != nil {
			fmt.Println("error:", err)
		}
		current := int(cur)
		if current > max {
			max = current
		}
	}

	return max
}

func MakeNextBlockData() (data Data) {

	var emptyTxions []Transaction

	data = Data{
		State:        LastGoatBlock.Data.State,
		Transactions: emptyTxions,
	}

	data.ApplyTransactions()

	return data
}

func (v *Vote) AddVote() (ok bool) {
	ok = v.VerifyVote()
	if ok != true {
		return false
	}

	VoteSet = append(VoteSet, *v)
	return ok
}

func (v *Vote) VerifyVote() (ok bool) {

	//check that the vote is not in the voteset already
	for _, vote := range VoteSet {
		if *v == vote {
			return false
		}
		//check that there is not another vote from the same account in this round
		if v.Account == vote.Account && v.Hash != vote.Hash {
			//TODO: penalize double voting for different hashes, slash deposit
			return false
		}
	}

	//verify signature of account sending the vote
	//TODO: abstract key recreation into a function (hash, r, s) (ok bool)
	hash := v.HashVote()
	//check that key is well formed
	if len(v.Account) < 100 {
		return false
	}
	//remove goat_ from key
	publicKey := v.Account[5:]
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
	byteR, err := hex.DecodeString(v.Signature.R)
	if err != nil {
		fmt.Println("error:", err)
	}
	r := new(big.Int).SetBytes(byteR)
	byteS, err := hex.DecodeString(v.Signature.S)
	if err != nil {
		fmt.Println("error:", err)
	}
	s := new(big.Int).SetBytes(byteS)
	//verify signature
	sigOk := ecdsa.Verify(pub, hash, r, s)

	if sigOk {
		fmt.Println("vote verified")
		return true
	} else {
		return false
	}
}

func (v *Vote) SignVote() (r string, s string) {
	keystore := LoadKeyStore()
	privateKey := keystore.PrivateKey

	hash := v.HashVote()
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

func (v *Vote) HashVote() (h []byte) {
	hashString := v.Account + v.Hash
	fixedHash := sha3.Sum512([]byte(hashString))
	h = fixedHash[:]
	return h
}

func SendVoteToNetwork() {
	config := LoadConfig()

	v := Vote{
		Account: config.Account,
		Hash:    hex.EncodeToString(HashCandidateSet(&CandidateSet)),
	}

	r, s := v.SignVote()

	v.Signature = Signature{
		R: r,
		S: s,
	}

	//TODO: remove self node from list of nodes

	for _, node := range config.Nodes {
		json, err := json.Marshal(v)
		if err != nil {
			fmt.Println("error:", err)
		}
		req, err := http.NewRequest("POST", "http://"+node+":3000/api/v1/vote", bytes.NewBuffer(json))
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

func CheckConsensus() {
	//TODO: final criteria for consensus = 2/3 of stakes sign hash of candidate set transaction

	voteCount := make(map[string]int)
	total := 0

	//TODO: get proportions for all, this won't work if this node is the one out of sync
	//changed the consensus to check votes in voteset
	//if agreed-upon hash is same as local hash, apply transactions
	//TODO: if not, wait and request new block

	for _, v := range VoteSet {
		voteCount[v.Hash] += 1
		total += 1
	}

	winningVote := 0
	var winningHash string
	for key, vote := range voteCount {
		if vote > winningVote {
			winningVote = vote
			winningHash = key
		}
	}

	if total < 1 {
		//catch up to network
		time.Sleep(3 * time.Second)
		maxBlock, nodes := GetMaxBlockNumberFromNetwork()
		if FindMaxBlock() < maxBlock {
			GetBlockFromNetwork(maxBlock, nodes[0])
			LastGoatBlock = ReadBlockFromLocalStorage(strconv.Itoa(maxBlock))
			//TODO: loop through to get real max block
		} else {
			fmt.Println("No consensus reached due to lack of votes")
		}
		return
	}

	if winningVote/total >= 2/3 {
		csHash := HashCandidateSet(&CandidateSet)
		vHash, err := hex.DecodeString(winningHash)
		if err != nil {
			fmt.Println("error:", err)
		}

		if bytes.Equal(vHash, csHash) {
			//reward voters on winning hash (add transactions to candidate set)
			RewardVoters(winningHash)
			//build next block
			//TODO: consider whether something needs to be refactored out to make this simpler
			NextBlock()
		} else {
			time.Sleep(3 * time.Second)
			maxBlock, nodes := GetMaxBlockNumberFromNetwork()
			GetBlockFromNetwork(maxBlock, nodes[0])
			LastGoatBlock = ReadBlockFromLocalStorage(strconv.Itoa(maxBlock))
		}
		//TODO: loop through to get real max block
	} else {
		//TODO: tell network to restart consensus round?
		fmt.Println("No consensus reached")
	}
}

func ResetVoteSet() {
	VoteSet = []Vote{}
}

func RewardVoters(hash string) {
	for _, v := range VoteSet {
		if v.Hash == hash {
			CandidateSet = append(CandidateSet, Transaction{
				From:   "mine",
				To:     v.Account,
				Amount: 100000000, //TODO: scale reward to stake
			})
		} else {
			CandidateSet = append(CandidateSet, Transaction{
				From:   "mine",
				To:     v.Account,
				Amount: -100000000, //TODO: scale reward to stake
			})
		}
	}
}

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

func NextBlock() {

	nextBlock := Block{
		Index:     LastGoatBlock.Index + 1,
		Timestamp: int(time.Now().UTC().Unix()),
		Data:      MakeNextBlockData(),
		LastHash:  LastGoatBlock.Hash,
	}

	nextBlock.HashBlock()

	DescribeBlock(nextBlock)

	LastGoatBlock = nextBlock

	nextBlock.WriteBlockToLocalStorage()

}

func DescribeBlock(b Block) {
	fmt.Printf("----------------------------------------------------------------------------------------------\n\n")
	fmt.Println("Block ID:", b.Index)
	fmt.Printf("\n---Block State---\n")
	for key, val := range b.Data.State {
		fmt.Printf("Account: %s\nBalance: %d\nSequence: %d\n", key, val.Balance, val.Sequence)
	}
	fmt.Printf("\n---Block Transactions---\n")
	for _, txion := range b.Data.Transactions {
		fmt.Printf("To: %s\nFrom: %s\nAmount: %d\n", txion.To, txion.From, txion.Amount)
	}
	fmt.Printf("\n---Hashes---\n")
	fmt.Println("Last Hash:", hex.EncodeToString(b.LastHash[:]))
	fmt.Println("Block Hash:", hex.EncodeToString(b.Hash[:]))
	fmt.Println("-----------")
	fmt.Println("Block Time:", time.Unix(int64(b.Timestamp), 0))
	fmt.Println("----------------------------------------------------------------------------------------------")
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

func (b *Block) WriteBlockToLocalStorage() {
	config := LoadConfig()

	storedBlock := StoredBlock{
		Index:     b.Index,
		Timestamp: b.Timestamp,
		Data:      b.Data,
		LastHash:  hex.EncodeToString(b.LastHash),
		Hash:      hex.EncodeToString(b.Hash),
	}

	//convert data to plain json
	out, err := json.Marshal(storedBlock)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//write json to file at config directory
	//TODO: check if file exists and don't overwrite
	err = ioutil.WriteFile(string(config.Directory)+strconv.Itoa(b.Index), out, 0644)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Block " + strconv.Itoa(b.Index) + " written successfully!")
	}

}

func GenerateAccount() (k KeyStore) {
	//create the keypair
	priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Println("error:", err)
	}

	//create the address from the public key variables
	pub := priv.PublicKey
	pubkey := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)

	k = KeyStore{
		PrivateKey: hex.EncodeToString(priv.D.Bytes()),
		PublicKey:  "goat_" + hex.EncodeToString(pubkey),
	}

	return k
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
	fmt.Println(hashString)
	fixedHash := sha3.Sum512([]byte(hashString))
	h = fixedHash[:]
	fmt.Println(string([]byte(hashString)))
	fmt.Println(hex.EncodeToString(h))
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
	spendOk := t.VerifyNegativeSpend()
	//do not check sequence here, so you can have more than one transaction per block
	if sigOk && spendOk {
		return true
	} else {
		return false
	}

}

func (t *Transaction) VerifyNegativeSpend() (ok bool) {
	if LastGoatBlock.Data.State[t.From].Balance < t.Amount {
		return false
	} else {
		return true
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
