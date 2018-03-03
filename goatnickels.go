package main

import(
  "fmt"
  "net/http"
  "github.com/seanmclane/goatnickels/block"
  "github.com/seanmclane/goatnickels/handler"
  "github.com/gorilla/mux"
  "log"
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

  //define server and routes for blockchain node
  //will remove the block chaining loop above
  r := mux.NewRouter().StrictSlash(true)
  s := r.PathPrefix("/api/v1").Subrouter()

  s.HandleFunc("/", handler.Index)
  s.HandleFunc("/txion/{test}", handler.Txion)

  log.Fatal(http.ListenAndServe(":3000", r))
}