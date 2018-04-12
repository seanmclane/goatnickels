package block

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/seanmclane/goatnickels/pubsub"
	"github.com/seanmclane/goatnickels/rpc"
	"golang.org/x/crypto/sha3"
	"io/ioutil"
	"net/http"
	"os"
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
	dir := os.Getenv("HOME")
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
	dir := os.Getenv("HOME")
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

// create http client for all network requests
var client = &http.Client{
	Timeout: time.Second * 10,
}

//bring in default hub for broadcasting messages
var hub = pubsub.DefaultHub

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

type Signature struct {
	R string `json:"r"`
	S string `json:"s"`
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

func WriteDefaultConfig() {
	dir := os.Getenv("HOME")
	config := Config{
		Directory: dir + "/.goatnickels/goatchain/",
		Nodes:     []string{"s1.goatnickels.com", "s2.goatnickels.com", "s3.goatnickels.com"},
		Account:   "",
	}

	file, err := json.Marshal(config)
	if err != nil {
		fmt.Println("error:", err)
	}

	os.MkdirAll(dir+"/.goatnickels/goatchain/", 0777)
	ioutil.WriteFile(dir+"/.goatnickels/config.json", file, 0777)
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
	fixedHash := sha3.Sum512([]byte("TESTNET"))
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
	fmt.Println("Got max block", maxBlock)
	fmt.Println("At nodes", nodes)
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

	hub.Broadcast <- rpc.BuildRequest(1, "getblock", []byte(`{"index": 1}`))

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

	//convert data to plain json
	out, err := json.Marshal(nextBlock)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	//broadcast the block as a json rpc message with the method "block"
	hub.Broadcast <- rpc.BuildNotification("block", out)

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
	err = ioutil.WriteFile(string(config.Directory)+strconv.Itoa(b.Index), out, 0777)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Block " + strconv.Itoa(b.Index) + " written successfully!")
	}

}

func GenerateAccount(save string) (k KeyStore) {
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
		PublicKey:  hex.EncodeToString(pubkey),
	}

	if save == "y" {
		dir := os.Getenv("HOME")

		keystore, err := json.Marshal(k)
		if err != nil {
			fmt.Println("error:", err)
		}

		c := LoadConfig()

		c.Account = "goat_" + k.PublicKey

		config, err := json.Marshal(c)
		if err != nil {
			fmt.Println("error:", err)
		}

		ioutil.WriteFile(dir+"/.goatnickels/keystore.json", keystore, 0777)
		ioutil.WriteFile(dir+"/.goatnickels/config.json", config, 0777)
	}

	k.PublicKey = "goat_" + k.PublicKey

	return k
}
