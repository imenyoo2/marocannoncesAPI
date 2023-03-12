package main

import (
  "net/http"
)

func (app *application) routes() http.Handler {
  mux := http.NewServeMux()
  mux.HandleFunc("/", app.home)
  mux.HandleFunc("/filter", app.filter)
  return mux
}

