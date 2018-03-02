package block

import(
  "fmt"
//  "os"
  "time"
  "encoding/json"
  "strconv"
  "golang.org/x/crypto/sha3"
  "encoding/hex"
)
//make a type blockchain and make functions methods on that
//specify json lowercase values with `json:"test"`

type Block struct {
  Index int
  Timestamp int
  Data []byte
  LastHash [64]byte
  Hash [64]byte
}

type Data struct {
  State []Account
  Transactions []Transaction
}

//change this to have an account as a key and balance in the object
// "HASH": {"balance": 834729564}
type Account struct {
  Account string
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

func CreateGenesisBlock() (b Block) {
  //set arbitrary data
  data := []byte("0")

  //genesis block for now
  b = Block {
    Index: 1,
    Timestamp: 0, //convert this to the birthdate of GoatNickels
    Data: data,
    LastHash: sha3.Sum512(data),
  }

  b.HashBlock()

  return b
}

func CreateBlockData() (byte_data []byte){
  //manually adding accounts and transactions for now
  sean := Account{
    Account: "sean",
    Balance: 50884325,
  }
  kate := Account{
    Account: "kate",
    Balance: 94043214,
  }
  sean_to_kate := Transaction{
    From: "sean",
    To: "kate",
    Amount: 140543,
  }

  state := []Account{sean, kate}
  transactions := []Transaction{sean_to_kate}
  
  data := Data{
    State: state,
    Transactions: transactions,
  }
  
  d, err := json.Marshal(data)
  if err != nil {
    fmt.Println("error:", err)
  }
//  os.Stdout.Write(d)
  fmt.Println("\n")

  byte_data = []byte(d)

  return byte_data
}

func NextBlock(last_block Block) (next_block Block) {
  next_block = Block {
    Index: last_block.Index+1,
    Timestamp: int(time.Now().UTC().Unix()),
    Data: CreateBlockData(),
    LastHash: last_block.Hash,
  }
  
  next_block.HashBlock()

  DescribeBlock(next_block)

  return next_block
}

func DescribeBlock(b Block) {
  fmt.Println("Block Data:", string(b.Data[:]))
  fmt.Println("Last Hash:", hex.EncodeToString(b.LastHash[:]))
  fmt.Println("Block Hash:", hex.EncodeToString(b.Hash[:]))
  fmt.Println("Block Time:", time.Unix(int64(b.Timestamp),0))
  fmt.Println("\n\n")
}