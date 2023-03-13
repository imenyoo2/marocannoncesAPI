package main

import (
 	"fmt"
 	"log"
 	"strconv"
 	"strings"
 	"time"
  "errors"

 	"github.com/gocolly/colly"
  mysql "github.com/go-sql-driver/mysql"
)

func (app *application) marocAnnonesCollect() {

  c := colly.NewCollector()
  // setting callback functions
  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
  })
  c.OnError(func(_ *colly.Response, err error) {
    log.Println("Something went wrong:", err)
  })

  // matching premium posts
  c.OnHTML("article.listing > a:nth-child(1)", func(e *colly.HTMLElement) {
    result := collectPage(collectPageParams{c: c, e: e, url: e.Attr("href")})
    result.premium[0] = 1
    result.time = "00:00:00"
    result.date = "2001-10-10"
    place := e.ChildText("div:nth-child(3) > span:nth-child(4)")
    result.place = place
    app.Insert(result)
  })

  // matching daily posts
  c.OnHTML(".cars-list > li", func(e *colly.HTMLElement) {
    url := e.ChildAttr("a:nth-child(1)", "href")
    if url != "" {
      time := (e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(1)") + 
                          " " + 
                          e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(3)"))
      place := e.ChildText("a:nth-child(1) > div:nth-child(2) > span:nth-child(2)")
      result := collectPage(collectPageParams{c: c, e: e, url: url , time: time})
      result.premium[0] = 0
      result.place = place
      var err error
      result.date, result.time, err = getTime(time)
      if err == nil {
        app.Insert(result)
      } else {
        // setting stop collect, this is true when 
        app.stopCollect = true
      }
    }
  })


  // uncomment to scrape the whole website
  c.OnHTML(".pagina_suivant > a:nth-child(1)", func(e *colly.HTMLElement) {
    if app.depth > 1 && !app.stopCollect {
      app.depth -= 1
      fmt.Printf("app.stopCollect = %t", app.stopCollect)
      e.Request.Visit(e.Attr("href"))
    }
  })

  c.SetRequestTimeout(1 * time.Minute)

  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

}

type collectPageParams struct {
  c *colly.Collector
  e *colly.HTMLElement
  url string
  //dataStore *map[string]map[string]interface{}
  time string
}


func toInt(arr []byte) int {
  multiplier := 1
  result := 0
  for i := 0; i < len(arr); i++ {
    multiplier *= 10
  }
  for _, v := range arr {
    multiplier /= 10
    result += (int(v) - 48) * multiplier
  }
  return result
}


func extractIdAndCatigorie(url string) (int, int, error) {
  base1 := 10
  if base1 > len(url) {
    return 0, 0, fmt.Errorf("expected a valid url found: %s", url)
  }
  var id []byte
  var cat []byte
  for i := base1; i < len(url); i++ {
    if url[i] >= byte('0') && url[i] <= byte('9') {
      cat = append(cat, url[i])
    } else if url[i] == '/' {
      break;
    } else {
      fmt.Printf("error in the url: %s\n", url)
      return 0, 0, fmt.Errorf("expected / found: %d", url[i])
    }
  }
  base2 := base1 + len(cat) + 23

  for i := base2; i < len(url); i++ {
    if url[i] >= '0' && url[i] <= '9' {
      id = append(id, url[i])
    } else if url[i] == '/' {
      break;
    } else {
      fmt.Printf("error in the url: %s\n", url)
      return 0, 0, fmt.Errorf("expected / found: %d", url[i])
    }
  }

  return toInt(cat),toInt(id), nil
}

func getTime(t string) (string, string, error) {
  parts :=  strings.Split(t, " ")
  if parts[0] == "Aujourd'hui" {
    y, m, d := time.Now().Date()
    return strings.Join([]string{strconv.Itoa(y), strconv.Itoa(int(m)), strconv.Itoa(d)}, "-"),parts[1] + ":00", nil
  } else if parts[0] == "Hier" {
    y, m, d := time.Now().Date()
    return strings.Join([]string{strconv.Itoa(y), strconv.Itoa(int(m)), strconv.Itoa(d - 1)}, "-"),parts[1] + ":00", nil
  } else {
    return "", "", fmt.Errorf("want Aujourd'hui, got %s\n", parts[0])
  }

}

// TODO: get rid of the .onhtml handlers (on every call the onhtml handlers redefined)
func collectPage(params collectPageParams) DBvalues{

  result := DBvalues{}
  var err error
  result.catigorie, result.id, err = extractIdAndCatigorie(params.url)
  check(err)
  result.url = params.url


  // matching the title
  params.c.OnHTML("#content > div.used-cars > div.description.desccatemploi > h1", func(e *colly.HTMLElement) {
    title := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "  ", "")
    e.Request.Ctx.Put("title", title)
    result.title = title
    // adding the url field

    // adding time if exist
    if params.time != "" {
      result.time = params.time
      // TODO add time field to DBvalues
    }
  })

  // matching Annonceur
  params.c.OnHTML(".infoannonce > dl:nth-child(1) > dd:nth-child(2)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Annonceur = e.Text
  })

  // matching Domaine
  params.c.OnHTML("#extraQuestionName > li:nth-child(1) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Domaine = e.Text
  })

  // matching Fonction
  params.c.OnHTML("#extraQuestionName > li:nth-child(2) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Fonction = e.Text
  })

  // matching Entreprise
  params.c.OnHTML("#extraQuestionName > li:nth-child(4) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Entreprise = e.Text
  })

  // matching Contrat
  params.c.OnHTML("#extraQuestionName > li:nth-child(3) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Contrat = e.Text
  })

  // matching Niveau d'études
  params.c.OnHTML("#extraQuestionName > li:nth-child(6) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Niveau = e.Text
  })

  // matching Salaire
  params.c.OnHTML("#extraQuestionName > li:nth-child(5) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    result.Salaire = e.Text
  })

  // start scraping the page
  params.e.Request.Visit(params.url)

  return result
}

func (app *application) Insert(values DBvalues) {
  stmt := `INSERT INTO posts (id, catigorie, url, title, Annonceur, Contrat, Domaine, Entreprise, Fonction, Niveau, Salaire, premium, date, time, place) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
  _, err := app.DB.Exec(stmt, 
                        values.id, 
                        values.catigorie, 
                        values.url, 
                        values.title, 
                        values.Annonceur, 
                        values.Contrat, 
                        values.Domaine,
                        values.Entreprise,
                        values.Fonction,
                        values.Niveau,
                        values.Salaire,
                        values.premium,
                        values.date,
                        values.time,
                        values.place,
                      )

  var mySQLError *mysql.MySQLError
  if err == nil {
    app.NewRecords += 1
  } else if errors.As(err, &mySQLError) { // finding a sqlError in err, and set it to mySQLError
    if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "posts_uc_id") {
      app.stopCollect = true
      app.DupRecords += 1
    } else {
      log.Fatal(err)
    }
  } else {
    log.Fatal(err)
  }
}


