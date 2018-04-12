package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/seanmclane/goatnickels/block"
	"github.com/seanmclane/goatnickels/handler"
	"github.com/seanmclane/goatnickels/pubsub"
	"log"
	"net/http"
	"time"
	//"encoding/hex"
)

func mine() {
	block.InitializeState()
	block.AsciiGoat()
	block.DescribeBlock(block.LastGoatBlock)

	for {
		voteStartList := [6]int{8, 18, 28, 38, 48, 58}
		voteEndList := [6]int{0, 10, 20, 30, 40, 50}

		for _, sec := range voteStartList {
			if time.Now().Second() == sec {
				//TODO: Move any transactions not reaching 80 percent to the staging set
				block.Voting = true
				fmt.Println("---------- Start voting round ---------")
				block.SendVoteToNetwork()
			}
		}
		for _, sec := range voteEndList {
			if time.Now().Second() == sec {
				block.CheckConsensus()
				block.Voting = false
				//add staging transactions back into candidate set
				//for now I will overwrite the candidate set with the staging set
				//TODO: ensure transactions that were not in the applied candidate set stay in the new candidate set with all the staging transactions
				block.ResetCandidateSet()
				block.ResetVoteSet()
				fmt.Println("---------- End voting round ---------")
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	genesisFlag := flag.String("genesis", "n", "create the genesis block?")
	serveFlag := flag.String("serve", "n", "y or n")
	acctFlag := flag.String("generate-acct", "n", "do you want generate a keypair for a new account?")
	saveKeyFlag := flag.String("save", "n", "saves the generated key")
	configFlag := flag.String("init-config", "n", "do you want to write the default config file?")
	testFlag := flag.String("test", "n", "do you want to test whatever you're working on now?")

	flag.Parse()

	if *genesisFlag == "y" {
		block.CreateGenesisBlock()
	}

	if *serveFlag == "y" {
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

		hub := pubsub.DefaultHub
		go hub.Run()

		s.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { handler.HandleConnections(hub, w, r) })

		srv := &http.Server{
			Addr:         ":3000",
			Handler:      r,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}

		go mine()
		go handler.ConnectToNodes(hub)
		log.Fatal(srv.ListenAndServe())
	}

	if *acctFlag == "y" {
		k := block.GenerateAccount(*saveKeyFlag)

		bytes, err := json.Marshal(k)
		if err != nil {
			fmt.Println("error:", err)
		}

		fmt.Println(string(bytes))
	}

	if *configFlag == "y" {
		block.WriteDefaultConfig()
	}

	if *testFlag == "y" {
		fmt.Println("test")
	}
}
