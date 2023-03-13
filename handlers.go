package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (app *application) home(w http.ResponseWriter, r *http.Request) {
  fmt.Println("test2 ----")
  app.parseJson(0, 0, "", "", "", "", "", "") // get all data
  jsonStr, err := json.Marshal(*app.data)
  check(err)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  writeto(w, jsonStr)
}

func (app *application) filter(w http.ResponseWriter, r *http.Request) {
  var err error
  params := struct {
    id         int
    catigorie  int
    salaire    string
    contrat    string
    domaine    string
    fonction   string
    niveau     string
    place      string
  }{}
  if r.URL.Query().Get("id") != "" {
    params.id, err = strconv.Atoi(r.URL.Query().Get("id"))
    check(err) // TODO: send http response instead
  }

  if r.URL.Query().Get("catigorie") != "" {
    params.catigorie, err = strconv.Atoi(r.URL.Query().Get("catigorie"))
    check(err) // TODO: send http response instead
  }

  params.salaire = r.URL.Query().Get("salaire")
  params.contrat = r.URL.Query().Get("contrat")
  params.domaine = r.URL.Query().Get("domaine")
  params.fonction = r.URL.Query().Get("fonction")
  params.niveau = r.URL.Query().Get("niveau")
  params.place = r.URL.Query().Get("place")

  app.parseJson(params.catigorie, params.id, params.salaire, params.contrat, params.domaine, params.fonction, params.niveau, params.place)

  // TODO: change this into a function
  jsonStr, err := json.Marshal(*app.data)
  check(err)
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusOK)
  writeto(w, jsonStr)
}
