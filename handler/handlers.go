package handler

import (
  "encoding/json"
  "net/http"
  "io/ioutil"
  "crypto/elliptic"
  "crypto/ecdsa"
  "crypto/rand"
  "encoding/hex"
  "github.com/seanmclane/goatnickels/block"
//  "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {

  bytes, err := json.Marshal(block.LastGoatBlock)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  w.Write(bytes)

}

func AddTxion(w http.ResponseWriter, r *http.Request) {
  
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var response block.Transaction
  err = json.Unmarshal(body, &response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  block.CandidateSet = append(block.CandidateSet, response)

  bytes, err := json.Marshal(block.CandidateSet)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  w.Write(bytes)

}

type AccountResponse struct {
  PrivateKey string
  PublicKey string
}

func AddAccount(w http.ResponseWriter, r *http.Request) {
  //return a public and private key
  //eventually this should be done in the wallet app

  //create the keypair
  priv, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  //create the address from the public key variables
  pub := priv.PublicKey
  pubkey := elliptic.Marshal(elliptic.P384(), pub.X, pub.Y)

  //make a new account in state with some goatnickels and the public key as the acct name


  response := AccountResponse{
    PrivateKey: hex.EncodeToString(priv.D.Bytes()),
    PublicKey: hex.EncodeToString(pubkey),
  }

  bytes, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  //return public and private keys in response writer body
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  w.Write(bytes)

}