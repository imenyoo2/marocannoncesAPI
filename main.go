package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
 	"database/sql"
 	_ "github.com/go-sql-driver/mysql"
)

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

func openDB(dnst string) (*sql.DB, error) {
  db, err := sql.Open("mysql", dnst)
  if err != nil {
    return nil, err
  }
  err = db.Ping()
  if err != nil {
    return nil, err
  }
  return db, nil
}

type application struct {
  data        *map[string]map[string]interface{}
  DB          *sql.DB
  depth       int
  stopCollect bool
  DupRecords  int 
  NewRecords  int
}

func main() {
  addr := flag.String("addr", ":4000", "HTTP network address")
  dnst := flag.String("dnst", "posts:1234@/marocannonces?parseTime=true", "MySQL data source name")
  startHTTP := flag.Bool("nohttp", false, "to stop the http server from running")
  depth := flag.Int("depth", 1, "the number of pages the collector should collect")
  flag.Parse()
  data := map[string]map[string]interface{}{}

  db, err := openDB(*dnst)
  check(err)

  app := &application {
    data:        &data,
    DB:          db,
    depth:       *depth,
    stopCollect: false,
    DupRecords:  0,
    NewRecords:  0,
  }
  app.marocAnnonesCollect()
  app.printSumarry()

  srv := &http.Server{
    Addr: *addr,
    Handler: app.routes(),
    IdleTimeout: time.Minute,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
  }
  if !*startHTTP {
    fmt.Printf("starting server %s\n", *addr)
    err = srv.ListenAndServe()
  }
  check(err)
}
