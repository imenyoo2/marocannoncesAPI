package main

import (
  "net/http"
  "encoding/json"
)

func writeto(w http.ResponseWriter, data []byte) {
  count, err := w.Write(data)
  check(err)
  if count != len(data) {
    for {
      count, err = w.Write(data[count:])
      if err != nil {
        check(err)
      }
      if count == 0 {
        return
      }
    }
  }
}

func (app *application) home (w http.ResponseWriter, r *http.Request) {
  jsonStr, err := json.Marshal(*app.data)
  check(err)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  writeto(w, jsonStr)
}
