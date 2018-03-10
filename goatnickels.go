package main

import(
  "fmt"
  "flag"
  "net/http"
  "time"
  "github.com/seanmclane/goatnickels/block"
  "github.com/seanmclane/goatnickels/handler"
  "github.com/gorilla/mux"
  "log"
)

func mine() {
  for {
    time.Sleep(15 * time.Second)
    block.NextBlock()
  }
}

func main() {
  serve_flag := flag.String("serve", "n", "y or n")
  acct_flag := flag.String("create-acct", "n", "do you want a new account?")
  sign_flag := flag.String("sign", "n", "do you want to sign a transaction?")
  //hash_flag := flag.String("hash", "no", "do you want to hash a transaction?")
  test_flag := flag.String("test", "n", "do you want to test whatever you're working on now?")

  flag.Parse()

  if *serve_flag == "y" {
    block.InitializeState()
    block.AsciiGoat()
    block.DescribeBlock(block.LastGoatBlock)
    
    //define server and routes for blockchain node
    //will remove the block chaining loop above
    r := mux.NewRouter().StrictSlash(true)
    s := r.PathPrefix("/api/v1").Subrouter()

    s.HandleFunc("/", handler.Index)
    s.HandleFunc("/txion", handler.AddTxion)

    go mine()
    log.Fatal(http.ListenAndServe(":3000", r))
  }
  
  if *acct_flag == "y" {
    block.CreateNewAccount()
  }

  if *sign_flag == "y" {
    //create transaction with flag values
    //sign
    t := block.Transaction{
      To: "kate",
      From: "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f",
      Amount: 1000,
      Nonce: 1,
    }
    r, s := block.SignTransaction(&t)
    fmt.Println("r:",r)
    fmt.Println("s:",s)
  }

  if *test_flag == "y" {
    t := block.Transaction{
      To: "kate",
      From: "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f",
      Amount: 1000,
      Nonce: 1,
    }
    r, s := block.SignTransaction(&t)
    t.R = r
    t.S = s
    ok := block.AddTransaction(&t)
    fmt.Println(ok)
  }
}