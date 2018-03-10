package handler

import (
  "strconv"
  "encoding/json"
  "net/http"
  "io/ioutil"
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

  var txion block.Transaction
  err = json.Unmarshal(body, &txion)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  ok := block.AddTransaction(&txion)

  response := []byte(`{"success":`+strconv.FormatBool(ok)+"}")

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  w.Write(response)

}