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
  //"encoding/hex"
)

func mine() {
  block.InitializeState()
  block.AsciiGoat()
  block.DescribeBlock(block.LastGoatBlock)

  for {
    if time.Now().Second() == 20 || time.Now().Second() == 50 {
      fmt.Println("Tally percent of vote/stake for each transaction")
      fmt.Println("Move any transactions not reaching 80 percent to the staging set")
      block.Voting = true
      fmt.Println("Broadcast candidate set vote to nodes")
      block.SendVoteToNetwork()
    }
    if time.Now().Second() == 30 || time.Now().Second() == 0 {
      block.CheckConsensus()
      block.Voting = false
      //add staging transactions back into candidate set
      //for now I will overwrite the candidate set with the staging set
      //TODO: ensure transactions that were not in the applied candidate set stay in the new candidate set with all the staging transactions
      block.ResetCandidateSet()
    }
    time.Sleep(1 * time.Second)
  }
}

func main() {
  genesis_flag := flag.String("genesis", "n", "create the genesis block?")
  serve_flag := flag.String("serve", "n", "y or n")
  acct_flag := flag.String("generate-acct", "n", "do you want generate a keypair for a new account?")
  sign_flag := flag.String("sign", "n", "do you want to sign a transaction?")
  //hash_flag := flag.String("hash", "no", "do you want to hash a transaction?")
  test_flag := flag.String("test", "n", "do you want to test whatever you're working on now?")

  //transaction flags
  to_flag := flag.String("to", "", "what account should be credited in this transaction?")
  from_flag := flag.String("from", "", "what account should be debited in this transaction?")
  amount_flag := flag.Int("amount", 0, "how much do you want to send?")
  //TODO: add sequence or figure out where to source sequence
  private_key_flag := flag.String("private-key", "", "needed to sign and send a transaction")

  flag.Parse()

  if *genesis_flag == "y" {
    block.CreateGenesisBlock()
  }

  if *serve_flag == "y" {    
    //define server and routes for blockchain node
    r := mux.NewRouter().StrictSlash(true)
    s := r.PathPrefix("/api/v1").Subrouter()

    s.HandleFunc("/", handler.Index)
    s.HandleFunc("/txion", handler.AddTxion).Methods("POST")
    s.HandleFunc("/txion", handler.GetTxions).Methods("GET")
    s.HandleFunc("/acct/{key}", handler.GetAcct).Methods("GET")
    s.HandleFunc("/block/{index}", handler.GetBlock).Methods("GET")
    s.HandleFunc("/maxblock", handler.GetMaxBlock).Methods("GET")
    s.HandleFunc("/sign", handler.SignTxion).Methods("POST")
    s.HandleFunc("/vote", handler.Vote).Methods("POST")

    srv := &http.Server{
      Addr: ":3000",
      Handler: r,
      ReadTimeout: 5 * time.Second,
      WriteTimeout: 10 * time.Second,
    }

    go mine()
    log.Fatal(srv.ListenAndServe())

  }
  
  if *acct_flag == "y" {
    block.GenerateAccount()
  }

  if *sign_flag == "y" {
    //create transaction with flag values
    t := block.Transaction{
      To: *to_flag,
      From: *from_flag,
      //From: "goat_04c12951412edfc215fe6d288491eb1251e2d8d99375c01049588dd228c6346f068246353d84702418f797d672af512d89742f6842b32f43541ea703f08170a67687f75fe0c6f15bd518764dee5476c86f9ba33f28036a76d018c1d7c8b14c307f",
      Amount: *amount_flag,
      Sequence: 1,
    }
    //private_key := "8b63849798d4633fe16553d428fdd50a1214296f0e02e5ebd0a7c78040a84775153a4dcacfc9dc7f4aeab9cc981fbb78"
    r, s := t.SignTransaction(*private_key_flag)
    fmt.Println("r:",r)
    fmt.Println("s:",s)
  }

  if *test_flag == "y" {
    fmt.Println("test")
  }
}