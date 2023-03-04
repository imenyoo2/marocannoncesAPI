package main

import (
  "net/http"
  "encoding/json"
)


func (app *application) home (w http.ResponseWriter, r *http.Request) {
  jsonStr, err := json.Marshal(*app.data)
  check(err)
  w.Header().Set("Content-Type", "application/json")
  w.Write(jsonStr)
}
