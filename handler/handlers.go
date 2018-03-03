package handler

import (
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
  
  //not sure what should be returned by the index
  //maybe it should be an ascii goat...

  test := struct {
    Test, Testing string
  }{
    Test: "yes",
    Testing: "hi",
  }

  bytes, err := json.Marshal(test)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  writeJsonResponse(w, bytes)

}

func Txion(w http.ResponseWriter, r *http.Request) {
  
  response := mux.Vars(r)

  bytes, err := json.Marshal(response)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }

  writeJsonResponse(w, bytes)

}

func writeJsonResponse(w http.ResponseWriter, bytes []byte) {
  w.Header().Set("Content-Type", "application/json; charset=UTF-8")
  w.Write(bytes)
}