package main

import (
	"fmt"
	"strings"
)

type DBvalues struct {
  id          int
  catigorie   int
  url         string
  title       string
  Annonceur   string
  Contrat     string
  Domaine     string
  Entreprise  string
  Fonction    string
  Niveau      string
  Salaire     string
  premium     int
  date        string
  time        string
  place       string
}


func (app *application) parseJson(catigorie int, id int, Salaire string, contrat string, Domaine string) {
  conditions := []string{}
  values := []interface{}{} // TODO: change from interface{}

  if catigorie != 0 {
    conditions = append(conditions, "catigorie = ?")
    values = append(values, catigorie)
  }
  if id != 0 {
    conditions = append(conditions, "id = ?")
    values = append(values, id)
  }
  if Salaire != "" {
    conditions = append(conditions, "Salaire = ?")
    values = append(values, Salaire)
  }
  if contrat != "" {
    conditions = append(conditions, "contrat = ?")
    values = append(values, contrat)
  }
  if Domaine != "" {
    conditions = append(conditions, "Domaine = ?")
    values = append(values, Domaine)
  }

  stmt := `SELECT id, catigorie, url, title, Annonceur, Contrat, Domaine, Entreprise, Fonction, Niveau, Salaire, premium, date, time, place FROM posts`
  // add conditions to stmt if any exist
  if len(conditions) > 0 {
    stmt +=  " WHERE " + strings.Join(conditions, " AND ")
  }

  rows, err := app.DB.Query(stmt, values...)
  defer rows.Close()
  check(err)

  row := DBvalues{}

  for rows.Next() {
    rows.Scan(&row.id,
              &row.catigorie,
              &row.url,
              &row.title,
              &row.Annonceur,
              &row.Contrat,
              &row.Domaine,
              &row.Entreprise,
              &row.Fonction,
              &row.Niveau,
              &row.Salaire,
              &row.premium,
              &row.date,
              &row.time,
              &row.place,
            )
    fmt.Printf("parsing row id=%d\n", row.id)
    (*app.data)[row.title] = map[string]interface{}{}
    (*app.data)[row.title]["id"] =  row.id
    (*app.data)[row.title]["catigorie"] = row.catigorie
    (*app.data)[row.title]["title"] = row.title
    (*app.data)[row.title]["Annonceur"] = row.Annonceur
    (*app.data)[row.title]["Contrat"] = row.Contrat
    (*app.data)[row.title]["Domaine"] = row.Domaine
    (*app.data)[row.title]["Entreprise"] = row.Entreprise
    (*app.data)[row.title]["Fonction"] = row.Fonction
    (*app.data)[row.title]["Niveau"] = row.Niveau
    (*app.data)[row.title]["Salaire"] = row.Salaire
    (*app.data)[row.title]["premium"] = row.premium
    (*app.data)[row.title]["date"] = row.date
    (*app.data)[row.title]["time"] = row.time
    (*app.data)[row.title]["place"] = row.place
  }



}
