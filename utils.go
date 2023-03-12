package main

import (
  "fmt"
  "log"
)

func (app *application) printSumarry() {
  fmt.Println("-------------------------------")
  fmt.Printf("Duplicated records: %d\n", app.DupRecords)
  fmt.Printf("New records: %d\n", app.NewRecords)
  fmt.Println("-------------------------------")
}

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}
