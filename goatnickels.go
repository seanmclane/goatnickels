package main

import(
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
  serve_flag := flag.Int("serve", 0, "1=yes, 0=no")
  acct_flag := flag.String("create-acct", "no", "do you want a new account?")
  sign_flag := flag.String("sign", "no", "do you want to sign a transaction?")
  hash_flag := flag.String("hash", "no", "do you want to hash a transaction?")
  test_flag := flag.String("test", "yes", "do you want to test whatever you're working on now?")

  flag.Parse()

  if *serve_flag == 1 {
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
  
  if *acct_flag == "yes" {
    block.CreateNewAccount()
  }

  if *sign_flag == "yes" {
    //create transaction with flag values
    //sign
  }

  if *test_flag == "yes" {
    t := Transaction{
      To: "kate",
      From: "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f",
      Amount: 1000,
      Nonce: 1,
      R: ,
      S: ,
    }
    r, s := t.block.SignTransaction()
    t.R = r
    t.S = s
    t.block.VerifyTransaction()
  }
}