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


//more for ease of testing than for true use, since sending your private key to an unknown server is bad...
type SignatureRequest struct {
  Transaction block.Transaction `json:"transaction"`
  Private_key string `json:"private_key"`
}

func SignTxion(w http.ResponseWriter, req *http.Request) {
  
  body, err := ioutil.ReadAll(req.Body)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  var sig_req SignatureRequest
  err = json.Unmarshal(body, &sig_req)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  r, s := block.SignTransaction(&sig_req.Transaction, sig_req.Private_key)

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