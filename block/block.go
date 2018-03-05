package block

import(
  "fmt"
  "io/ioutil"
  "os"
  "time"
  "encoding/json"
  "strconv"
  "golang.org/x/crypto/sha3"
  "encoding/hex"
)

//define config structure
//does this belong here?
type Config struct {
  Directory string `json:"directory"`
  LastBlock int `json:"last_block"`
}

//load config
func LoadConfig() (config Config) {
  c, err := os.Open("config.json")
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

//initializing blockchain objects here for now
//create candidate set of transactions
var CandidateSet []Transaction

//make a type blockchain and make functions methods on that
//specify json lowercase values with `json:"test"`

type Blockchain []Block

var GoatChain Blockchain

type Block struct {
  Index int
  Timestamp int
  Data []byte
  LastHash [64]byte
  Hash [64]byte
}

type Data struct {
  State map[string]Account
  Transactions []Transaction
}

//initialize state
var Accounts map[string]Account

type Account struct {
  Balance int
}

type Transaction struct {
  From string
  To string
  Amount int
}

func AsciiGoat() {
  a := "\x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2d \x2e \x5f \x2c \x2d \x2d \x2e \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2c \x2d \x60 \x2e \x0a \x28 \x60 \x2d \x2e \x5f \x5f \x20 \x20 \x20 \x20 \x2f \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2f \x0a \x20 \x60 \x2e \x20 \x20 \x20 \x60 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x5f \x5f \x2c \x2d \x2d \x27 \x2d \x2e \x0a \x20 \x20 \x20 \x60 \x2d \x2d \x2f \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2e \x20 \x20 \x5f \x5f \x5f \x5f \x5f \x5f \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x2d \x2e \x20 \x20 \x20 \x20 \x20 \x2c \x6f \x2d \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x2e \x20 \x3b \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x7c \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x60 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x20 \x6f \x20 \x2c \x20 \x20 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x2d \x2d \x27 \x2c \x27 \x2e \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x3b \x3b \x20 \x20 \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x27 \x20 \x20 \x3b \x20 \x20 \x2c \x27 \x20 \x2c \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x2c \x27 \x20 \x20 \x3a \x20 \x20 \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x20 \x5c \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60"
  fmt.Printf("Goatnickels baby!\n")
  fmt.Println(a, "\n")
}

func (b *Block) HashBlock() {
  //create a hash of all values in the block
  block_string := strconv.Itoa(b.Index)+strconv.Itoa(b.Timestamp)+string(b.Data)+hex.EncodeToString(b.LastHash[:])
  b.Hash = sha3.Sum512([]byte(block_string))
}

func (bc *Blockchain) CreateGenesisBlock() {
  //set arbitrary data
  data := []byte("0")

  //genesis block for now
  b := Block {
    Index: 0,
    Timestamp: 0, //TODO: convert this to the birthdate of GoatNickels
    Data: data,
    LastHash: sha3.Sum512(data),
  }

  b.HashBlock()

  (*bc) = append((*bc), b)

  b.WriteBlockToLocalStorage()

}

func (bc *Blockchain) InitializeState() {
  //temporary
  //manually adding accounts and transactions for now
  Accounts = make(map[string]Account)

  Accounts["sean"] = Account{
    Balance: 50884325,
  }
  Accounts["kate"] = Account{
    Balance: 94043214,
  }

  //check last block
  //if no blockchain, start a new one
  //TODO: if no blockchain, get it from the network instead
  config := LoadConfig()
  if config.LastBlock < 0 {
    GoatChain.CreateGenesisBlock()
    config.LastBlock = 0
  }

  b, err := ioutil.ReadFile(string(config.Directory)+strconv.Itoa(config.LastBlock)+".json")
  if err != nil {
    panic(err)
  }

  //make bytestring to Block
  last_block := Block{}
  err = json.Unmarshal(b, &last_block)
  if err != nil {
    fmt.Println("error:", err)
  }

  //PROBLEM: can't just append, needs to be in the right position in array for NextBlock to work...
  for i := 0; i < last_block.Index; i++ {
    (*bc) = append((*bc), Block{})
  }

  (*bc) = append((*bc), last_block)

}

func CreateBlockData() (byte_data []byte){

  state := Accounts
  transactions := CandidateSet
  
  data := Data{
    State: state,
    Transactions: transactions,
  }

  data.ApplyTransactionsToState()
  
  d, err := json.Marshal(data)
  if err != nil {
    fmt.Println("error:", err)
  }

  fmt.Println("\n")

  byte_data = []byte(d)

  return byte_data
}

func (bc *Blockchain) NextBlock() {

  config := LoadConfig()

  next_block := Block {
    Index: (*bc)[config.LastBlock].Index+1,
    Timestamp: int(time.Now().UTC().Unix()),
    Data: CreateBlockData(),
    LastHash: (*bc)[config.LastBlock].Hash,
  }
  
  next_block.HashBlock()

  DescribeBlock(next_block)

  (*bc) = append((*bc), next_block)

  next_block.WriteBlockToLocalStorage()

}

func DescribeBlock(b Block) {
  fmt.Println("Block Data:", string(b.Data[:]))
  fmt.Println("Last Hash:", hex.EncodeToString(b.LastHash[:]))
  fmt.Println("Block Hash:", hex.EncodeToString(b.Hash[:]))
  fmt.Println("Block Time:", time.Unix(int64(b.Timestamp),0))
  fmt.Println("\n\n")
}

func (d *Data) ApplyTransactionsToState() {
  //add and subtract from accounts
  for _, txion := range CandidateSet {
    fnb := d.State[txion.From].Balance - txion.Amount
    d.State[txion.From] = Account{Balance: fnb} 
    tnb := d.State[txion.To].Balance + txion.Amount
    d.State[txion.To] = Account{Balance: tnb} 
  }
  //reset candidate transactions to apply
  CandidateSet = nil

}

func (b *Block) WriteBlockToLocalStorage() {
  config := LoadConfig()

  //convert block to json
  //TODO: convert hashes to hex encoding
  //TODO: convert data to plain json
  out, err := json.Marshal(b)
  if err != nil {
    fmt.Println("error:", err)
  }

  //write json to file at config directory
  //TODO: check if file exists and don't overwrite
  err = ioutil.WriteFile(string(config.Directory)+strconv.Itoa(b.Index)+".json", out, 0644)
  if err != nil {
      panic(err)
  }

  //make sure config has record of last block mined
  //TODO: refactor this when setting last block mined from checking network
  config.LastBlock = b.Index
  cout, err := json.Marshal(config)
  if err != nil {
    fmt.Println("error:", err)
  }
  
  err = ioutil.WriteFile("config.json", cout, 0644)
  if err != nil {
      panic(err)
  }
  fmt.Println("Block written successfully!")
}

