package main

import (
  "fmt"
)

func (app *application) printSumarry() {
  fmt.Println("-------------------------------")
  fmt.Printf("Duplicated records: %d\n", app.DupRecords)
  fmt.Printf("New records: %d\n", app.NewRecords)
  fmt.Println("-------------------------------")
}
