package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type AnnouncesCollector struct {
  collector  *colly.Collector
  data       *DBvalues
  prolinks   *[]string
  links      *[]string
}

func (app *application) marocAnnonesCollect() {

  c := AnnouncesCollector{}
  c.links = &[]string{}
  c.prolinks = &[]string{}

  c.data = &DBvalues{}

  c.newCollector().collectDailyLinks().collectPremiumLinks().collectPage()

  c.collector.SetRequestTimeout(1 * time.Minute)

  c.collector.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

  app.visitLinks(&c)

}

// collectPageParams used to pass parameters to collectPage
// time is an optional parameter
type collectPageParams struct {
  c *colly.Collector
  e *colly.HTMLElement
  url string
  time string
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
// TODO: move onhtml somewhere where it only get called once
// collectPage adds the callbacks that handle collecting page data, and placing it in *AnnouncesCollector.data field
func (ac *AnnouncesCollector) collectPage() {

	ac.collector.OnHTML(".description", func(e *colly.HTMLElement) {
    title := strings.TrimSpace(e.ChildText("h1"))
    ac.data.title = title

		e.ForEach("li", func(i int, e *colly.HTMLElement) { // intresting
      switch i { // isn't this a loop over all children that are of type li, TODO: make a print statment to test this
			case 0:
				ac.data.place = e.ChildText("a")
			case 1:
				ac.data.time, ac.data.date = extractDateAndTime(e.Text)
        fmt.Printf("time = %s, date = %s", ac.data.time, ac.data.date)
      case 2:
        ac.data.vue = e.Text
			}
		})

		e.ForEach(".extraQuestionName", func(_ int, e *colly.HTMLElement) {
			e.ForEach("li", func(idx int, e *colly.HTMLElement) {
        val := strings.TrimSpace(e.ChildText("a")) // TODO: use trimspace in ur code
				switch idx {
				case 0:
					ac.data.Domaine = val
				case 1:
					ac.data.Fonction = val
				case 2:
					ac.data.Contrat = val
				case 3:
					ac.data.Entreprise = val
				case 4:
					ac.data.Salaire = val
				case 5:
					ac.data.Niveau = val
				}
			})
		})
		ac.data.Annonceur = e.ChildText(".infoannonce > dl:nth-child(1) > dd:nth-child(2)")
	})
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
  check(err)
//  var mySQLError *mysql.MySQLError
//  if err == nil {
//    app.NewRecords += 1
//  } else if errors.As(err, &mySQLError) { // finding a sqlError in err, and set it to mySQLError
//    if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "posts_uc_id") {
//      app.stopCollect = true
//      app.DupRecords += 1
//    } else {
//      log.Fatal(err)
//    }
//  } else {
//    log.Fatal(err)
//  }
}

func (ac *AnnouncesCollector) newCollector() *AnnouncesCollector{
  c := colly.NewCollector()
  // setting callback functions
  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
  })
  c.OnError(func(_ *colly.Response, err error) {
    log.Println("Something went wrong:", err)
  })

  ac.collector = c
  return ac
}


func (ac *AnnouncesCollector) collectPremiumLinks() *AnnouncesCollector {

  // matching premium posts
  ac.collector.OnHTML("article.listing > a:nth-child(1)", func(e *colly.HTMLElement) {
    // added the link
    *ac.prolinks = append(*ac.prolinks, e.Attr("href"))
    //result := collectPage(collectPageParams{c: c, e: e, url: e.Attr("href")})
    //result.premium = []byte{1}
    //result.time = "00:00:00"
    //result.date = "2001-10-10"
    //place := e.ChildText("div:nth-child(3) > span:nth-child(4)")
    //result.place = place
    //app.Insert(result)
  })
  return ac
}

func (ac *AnnouncesCollector) collectDailyLinks() *AnnouncesCollector{

  // matching daily posts
  ac.collector.OnHTML(".cars-list > li", func(e *colly.HTMLElement) {
    url := e.ChildAttr("a:nth-child(1)", "href")
    *ac.links = append(*ac.links, url)
//    if url != "" {
//      time := (e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(1)") + 
//                          " " + 
//                          e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(3)"))
//      place := e.ChildText("a:nth-child(1) > div:nth-child(2) > span:nth-child(2)")
//      result := collectPage(collectPageParams{c: c, e: e, url: url , time: time})
//      result.premium = []byte{0}
//      result.place = place
//      var err error
//      result.date, result.time, err = getTime(time)
//      if err == nil {
//        app.Insert(result)
//      } else {
//        // setting stop collect, this is true when 
//        app.stopCollect = true
//      }
//    }
  })
  return ac 
}

func (ac *AnnouncesCollector) collectNext(depth *int, stopCollect bool) {

  // the 'Suivant' button
  ac.collector.OnHTML(".pagina_suivant > a:nth-child(1)", func(e *colly.HTMLElement) {
    if *depth > 1 && !stopCollect {
      *depth -= 1
      fmt.Printf("app.stopCollect = %t", stopCollect)
      e.Request.Visit(e.Attr("href"))
    }
  })
}

func (app *application) visitLinks(ac *AnnouncesCollector) {
  base := "https://www.marocannonces.com/"
  for _, v := range *ac.links {
    if v == "" {
      continue
    }
    ac.collector.Visit(base + v)
    cat, id, err := extractIdAndCatigorie(v)
    check(err)
    ac.data.id = id
    ac.data.catigorie = cat
    ac.data.premium = []byte{0}
    fmt.Println("----- inserting values -----")
    if !ac.checkExistRecord(app.DB, id) {
      app.Insert(*ac.data)
      app.NewRecords += 1
    } else {
      app.DupRecords += 1
    }
  }
  for _, v := range *ac.prolinks {
    ac.collector.Visit(base + v)
    cat, id, err := extractIdAndCatigorie(v)
    check(err)
    ac.data.id = id
    ac.data.catigorie = cat
    ac.data.premium = []byte{1}
    fmt.Println("----- inserting values -----")
    if !ac.checkExistRecord(app.DB, id) {
      app.Insert(*ac.data)
      app.NewRecords += 1
    } else {
      app.DupRecords += 1
    }
  }
}

func (ac *AnnouncesCollector) checkExistRecord (db *sql.DB, id int) bool {
  stmt := `SELECT id FROM posts WHERE id = ?`


  row := db.QueryRow(stmt, id)

  var check int
  err := row.Scan(&check)
  if errors.Is(err, sql.ErrNoRows) {
    return false
  } else {
    return true
  }
}
