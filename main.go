package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

type application struct {
  data      *map[string]map[string]interface{}
}

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")
  flag.Parse()
  data := map[string]map[string]interface{}{}


  app := &application {
    data: &data,
  }
  app.marocAnnonesCollect()

  srv := &http.Server{
    Addr: *addr,
    Handler: app.routes(),
    IdleTimeout: time.Minute,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }
  fmt.Printf("starting server %s", *addr)
  err := srv.ListenAndServe()
  check(err)
}
