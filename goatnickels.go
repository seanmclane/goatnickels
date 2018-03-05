package main

import(
  "net/http"
  "time"
  "github.com/seanmclane/goatnickels/block"
  "github.com/seanmclane/goatnickels/handler"
  "github.com/gorilla/mux"
  "log"
)

func mine() {
  for {
    time.Sleep(5 * time.Second)
    block.GoatChain.NextBlock()
  }
}

func main() {
  block.LoadConfig()
  
  block.InitializeState()
  block.GoatChain.CreateGenesisBlock()

  block.AsciiGoat()
  block.DescribeBlock(block.GoatChain[0])
  
  //define server and routes for blockchain node
  //will remove the block chaining loop above
  r := mux.NewRouter().StrictSlash(true)
  s := r.PathPrefix("/api/v1").Subrouter()

  s.HandleFunc("/", handler.Index)
  s.HandleFunc("/txion", handler.AddTxion)

  go mine()
  log.Fatal(http.ListenAndServe(":3000", r))
}