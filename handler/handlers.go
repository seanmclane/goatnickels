package handler

import (
  "strconv"
  "encoding/json"
  "net/http"
  "io/ioutil"
  "github.com/seanmclane/goatnickels/block"
  "github.com/gorilla/mux"
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

  ok := txion.AddTransaction()

  response := []byte(`{"success":`+strconv.FormatBool(ok)+"}")

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  if ok {
    w.WriteHeader(http.StatusCreated)
  } else {
    w.WriteHeader(http.StatusBadRequest)
  }
  w.Write(response)

}

func GetTxions(w http.ResponseWriter, r *http.Request) {
  
  bytes, err := json.Marshal(block.CandidateSet)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  w.Write(bytes)

}

func GetAcct(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  key := vars["key"]

  val := block.LastGoatBlock.Data.State[key]

  empty := block.Account{}

  if val != empty {
    bytes, err := json.Marshal(val)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
    }
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    w.Write(bytes)
  } else {
    bytes := []byte(`{"error":"Account not found"}`)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    w.Write(bytes)
  }  

}

func GetBlock(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  index := vars["index"]

  bytes := block.ReadBlockFromLocalStorage(index)

  if bytes != nil { 
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
    w.Write(bytes)
  } else {
    bytes := []byte(`{"error":"Block not found"}`)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusNotFound)
    w.Write(bytes)
  }  
}

type MaxBlockResponse struct {
  MaxBlock int `json:"max_block"`
}

func GetMaxBlock(w http.ResponseWriter, r *http.Request) {
  max := block.FindMaxBlock()
  response := MaxBlockResponse{
    MaxBlock: max,
  }
  bytes, err := json.Marshal(response)
  if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
  }
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusOK)
  w.Write(bytes)
}

func Vote(w http.ResponseWriter, r *http.Request) {
  //TODO: fix these ReadAlls to use json.NewDecoder instead
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var vote block.Vote
  err = json.Unmarshal(body, &vote)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  ok := vote.AddVote()

  response := []byte(`{"success":`+strconv.FormatBool(ok)+"}")

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  if ok {
    w.WriteHeader(http.StatusCreated)
  } else {
    w.WriteHeader(http.StatusBadRequest)
  }
  w.Write(response)
}






//more for ease of testing than for true use, since sending your private key to an unknown server is bad...
type SignatureRequest struct {
  Transaction block.Transaction `json:"transaction"`
  PrivateKey string `json:"private_key"`
}

func SignTxion(w http.ResponseWriter, req *http.Request) {
  
  body, err := ioutil.ReadAll(req.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var sigReq SignatureRequest
  err = json.Unmarshal(body, &sigReq)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  r, s := sigReq.Transaction.SignTransaction(sigReq.PrivateKey)

  var response []byte
  if r != "" && s != "" {
      sig := block.Signature {
        R: r,
        S: s,
      }
      response, err = json.Marshal(sig)
      if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
      }
    } else {
      response = []byte(`{"success": false}`)
    }

  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.WriteHeader(http.StatusCreated)
  w.Write(response)

}
