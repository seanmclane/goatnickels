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

		go handler.ConnectToNodes(hub)
		go func() {
			err := srv.ListenAndServe()
			if err != nil {
				log.Panic(err)
			}
		}()
		block.Run()
		block.InitializeState()
		block.AsciiGoat()
		block.DescribeBlock(block.LastGoatBlock)
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
