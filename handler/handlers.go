package handler

import (
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/seanmclane/goatnickels/block"
//  "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {

  bytes, err := json.Marshal(block.CreateGenesisBlock())
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  w.Write(bytes)

}

func NewBlock(w http.ResponseWriter, r *http.Request) {

  bytes, err := json.Marshal(block.GoatChain)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
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