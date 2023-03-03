package main

import (
	"log"

)

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}

type application struct {
  data      map[string]map[string]interface{}
}

func main() {
  // initializing a new collecotor
  data := map[string]map[string]interface{}{}

  app := &application {
    data: data,
  }
  app.marocAnnonesCollect()



}
