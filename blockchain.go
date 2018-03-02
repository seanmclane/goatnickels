package main

import(
  "fmt"
  "time"
  "strconv"
  "golang.org/x/crypto/sha3"
  "encoding/hex"
)
//make a type blockchain and make functions methods on that

type Block struct {
  Index int
  Timestamp int
  Data []byte
  LastHash [64]byte
  Hash [64]byte
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

func NextBlock(last_block Block) (next_block Block) {
  next_block = Block {
    Index: last_block.Index+1,
    Timestamp: int(time.Now().UTC().Unix()),
    Data: last_block.Data,
    LastHash: last_block.Hash,
  }
  
  //create a hash of all values in the block
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


func main() {
  b := CreateGenesisBlock()

  AsciiGoat()
  DescribeBlock(b)
  
  //initialize a blockchain with the genesis block
  blockchain := []Block{b}

  for i := 1; i < 5; i++ {
    //add blocks to the chain at the right index
    fmt.Println("Iteration",i)
    next_block := NextBlock(blockchain[i-1])
    blockchain = append(blockchain, next_block)

  }
}