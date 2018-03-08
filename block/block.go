package block

import(
  "fmt"
  "io/ioutil"
  "os"
  "time"
  "encoding/json"
  "strconv"
  "golang.org/x/crypto/sha3"
  "crypto/elliptic"
  "crypto/ecdsa"
  "crypto/rand"
  "math/big"
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
//need to have last validated block (index at minimum)
var LastGoatBlock Block

//need to have last state
var Accounts map[string]Account

//create candidate set of transactions
var CandidateSet []Transaction

//specify json lowercase values with `json:"test"`

//removing blockchain since it's not needed to store whole chain in memory
// type Blockchain []Block
// var GoatChain Blockchain

type Block struct {
  Index int `json:"index"`
  Timestamp int `json:"timestamp"`
  Data Data `json:"data"`
  LastHash [64]byte `json:"last_hash"`
  Hash [64]byte `json:"hash"`
}

type StoredBlock struct {
  Index int `json:"index"`
  Timestamp int `json:"timestamp"`
  Data Data `json:"data"`
  LastHash string `json:"last_hash"`
  Hash string `json:"hash"`
}

type Data struct {
  State map[string]Account `json:"state"`
  Transactions []Transaction `json:"transactions"`
}

type Account struct {
  Balance int `json:"balance"`
}

type Transaction struct {
  From string `json:"from"`
  To string `json:"to"`
  Amount int `json:"amount"`
  Nonce int `json:"nonce"`
  R string `json:"r"`
  S string `json:"s"`
}

func AsciiGoat() {
  a := "\x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2d \x2e \x5f \x2c \x2d \x2d \x2e \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2c \x2d \x60 \x2e \x0a \x28 \x60 \x2d \x2e \x5f \x5f \x20 \x20 \x20 \x20 \x2f \x20 \x20 \x2c \x27 \x20 \x20 \x20 \x2f \x0a \x20 \x60 \x2e \x20 \x20 \x20 \x60 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x5f \x5f \x2c \x2d \x2d \x27 \x2d \x2e \x0a \x20 \x20 \x20 \x60 \x2d \x2d \x2f \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x2d \x2e \x20 \x20 \x5f \x5f \x5f \x5f \x5f \x5f \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x2d \x2e \x20 \x20 \x20 \x20 \x20 \x2c \x6f \x2d \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x2e \x20 \x3b \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x7c \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x60 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x20 \x20 \x20 \x5c \x0a \x20 \x20 \x20 \x20 \x20 \x28 \x6f \x20 \x6f \x20 \x2c \x20 \x20 \x2d \x2d \x27 \x20 \x20 \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x2d \x2d \x27 \x2c \x27 \x2e \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60 \x3b \x3b \x20 \x20 \x3a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2f \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x3b \x27 \x20 \x20 \x3b \x20 \x20 \x2c \x27 \x20 \x2c \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x2c \x27 \x2c \x27 \x20 \x20 \x3a \x20 \x20 \x27 \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x5c \x20 \x5c \x20 \x20 \x20 \x3a \x0a \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x20 \x60"
  fmt.Println(a, "\n")
}

func (b *Block) HashBlock() {
  //create a hash of all values in the block
  //TODO: handle error
  hash_data, _ := json.Marshal(b.Data)
  block_string := strconv.Itoa(b.Index)+strconv.Itoa(b.Timestamp)+string(hash_data)+hex.EncodeToString(b.LastHash[:])
  b.Hash = sha3.Sum512([]byte(block_string))
}

func CreateGenesisBlock() {
  //temporary
  //manually adding accounts for now
  Accounts = make(map[string]Account)

  Accounts["sean"] = Account{
    Balance: 50884325,
  }
  Accounts["kate"] = Account{
    Balance: 94043214,
  }

  //set arbitrary data
  data := Data{
    State: Accounts,
    Transactions: nil,
  }

  //genesis block for now
  b := Block {
    Index: 1,
    Timestamp: 0, //TODO: convert this to the birthdate of GoatNickels
    Data: data,
    LastHash: sha3.Sum512([]byte("Goatnickels baby!")),
  }

  b.HashBlock()

  b.WriteBlockToLocalStorage()

}

func InitializeState() {
  //check last block
  //if no blockchain, start a new one
  //TODO: if no blockchain, get it from the network instead
  config := LoadConfig()
  if config.LastBlock < 1 {
    CreateGenesisBlock()
    config.LastBlock = 1
  }

  b, err := ioutil.ReadFile(string(config.Directory)+strconv.Itoa(config.LastBlock))
  if err != nil {
    panic(err)
  }

  //make bytestring to Block
  err = json.Unmarshal(b, &LastGoatBlock)
  if err != nil {
    fmt.Println("error:", err)
  }

}

func CreateBlockData() (data Data){
  
  data = Data{
    State: LastGoatBlock.Data.State,
    Transactions: CandidateSet,
  }

  data.ApplyTransactionsToState()
  
  return data
}

func NextBlock() {

  next_block := Block {
    Index: LastGoatBlock.Index+1,
    Timestamp: int(time.Now().UTC().Unix()),
    Data: CreateBlockData(),
    LastHash: LastGoatBlock.Hash,
  }
  
  next_block.HashBlock()

  DescribeBlock(next_block)

  LastGoatBlock = next_block

  next_block.WriteBlockToLocalStorage()

}

func DescribeBlock(b Block) {
  fmt.Println("Block ID:", b.Index)
  fmt.Println("Block State:", b.Data.State)
  fmt.Println("Block Transactions:", b.Data.Transactions)
  fmt.Println("Last Hash:", hex.EncodeToString(b.LastHash[:]))
  fmt.Println("Block Hash:", hex.EncodeToString(b.Hash[:]))
  fmt.Println("Block Time:", time.Unix(int64(b.Timestamp),0))
  fmt.Println("\n\n")
}

func (d *Data) ApplyTransactionsToState() {
  //add and subtract from accounts
  for _, txion := range CandidateSet {
    //TODO: check if transaction results in negative balance
    //TODO: check if nonce is being reused


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

  //convert data to plain json
  out, err := json.Marshal(b)
  if err != nil {
    fmt.Println("error:", err)
  }

  //write json to file at config directory
  //TODO: check if file exists and don't overwrite
  err = ioutil.WriteFile(string(config.Directory)+strconv.Itoa(b.Index), out, 0644)
  if err != nil {
      panic(err)
  }

  //make sure config has record of last block mined
  //should this be from checking the file names locally for now?
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

type AccountResponse struct {
  PrivateKey string `json:"private_key"`
  PublicKey string `json:"public_key"`
}

func CreateNewAccount() {
  //create the keypair
  priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
  if err != nil {
    fmt.Println("error:", err)
  }

  //create the address from the public key variables
  pub := priv.PublicKey
  pubkey := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)

  response := AccountResponse{
    PrivateKey: hex.EncodeToString(priv.D.Bytes()),
    PublicKey: "goat_"+hex.EncodeToString(pubkey),
  }

  bytes, err := json.Marshal(response)
  if err != nil {
    fmt.Println("error:", err)
  }

  fmt.Println(string(bytes))
}

//TODO: make this real and not a test of some hardcoded values
func SignTransaction(t *Transaction) (r, s string) {
  hash := sha3.Sum512([]byte("Goatnickels baby!"))
  private_key := "8b63849798d4633fe16553d428fdd50a1214296f0e02e5ebd0a7c78040a84775153a4dcacfc9dc7f4aeab9cc981fbb78"
  //recreate ecdsa.PrivateKey from priv
  byte_key, _ := hex.DecodeString(private_key)
  bigint_key := new(big.Int).SetBytes(byte_key)
  priv := new(ecdsa.PrivateKey)
  priv.PublicKey.Curve = elliptic.P384()
  priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(byte_key)

  //TODO: handle error
  r_int, s_int, _ := ecdsa.Sign(rand.Reader, priv, hash)
  r := hex.EncodeToString(r_int)
  s := hex.EncodeToString(s_int.Bytes())
  return r, s
}

func HashTransaction(t *Transaction) (h string) {
  h := hex.EncodeToString(sha3.Sum512([]byte("Goatnickels baby!")))
  //TODO: hash the actual transaction
}

func (t *Transaction) VerifySignature() {
  //check if t.Signature ok with public key?
  //what is being signed exactly? hash of transaction nonce, to, from, and amount
  hash := sha3.Sum512([]byte("Goatnickels baby!"))
  pub := "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f"
  //remove goat_ from key
  pub = pub[4:]
  
  //TODO: convert r and s back


  //verify signature
  ecdsa.Verify(*pub, hash, t.R, t.S)

}