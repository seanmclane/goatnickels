package main

import(
  "fmt"
  "github.com/seanmclane/goatnickels/block"
)

func main() {
  b := block.CreateGenesisBlock()

  block.AsciiGoat()
  block.DescribeBlock(b)
  
  //initialize a blockchain with the genesis block
  blockchain := []block.Block{b}

  for i := 1; i < 5; i++ {
    //add blocks to the chain at the right index
    fmt.Println("Iteration",i)
    next_block := block.NextBlock(blockchain[i-1])
    blockchain = append(blockchain, next_block)

  }
}