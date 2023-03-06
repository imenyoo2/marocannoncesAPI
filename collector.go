package main

import (
  "fmt"
  "log"
  "strings"
	"github.com/gocolly/colly"
  "time"
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

  c.OnHTML("article.listing > a:nth-child(1)", func(e *colly.HTMLElement) {
    collectPage(c, e, e.Attr("href"), app.data)
  })


  c.SetRequestTimeout(1 * time.Minute)

  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

}

func collectPage(c *colly.Collector, e *colly.HTMLElement, url string, buffer *map[string]map[string]interface{}) {

  e.Request.Visit(url)

  // matching the title
  c.OnHTML("#content > div.used-cars > div.description.desccatemploi > h1", func(e *colly.HTMLElement) {
    title := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "  ", "")
    e.Request.Ctx.Put("title", title)
    (*buffer)[title] = map[string]interface{}{"title": title}
    // adding the url field
    (*buffer)[e.Request.Ctx.Get("title")]["URL"] = url
  })


  // matching Annonceur
  c.OnHTML(".infoannonce > dl:nth-child(1) > dd:nth-child(2)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Annonceur"] = e.Text
  })

  // matching Domaine
  c.OnHTML("#extraQuestionName > li:nth-child(1) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Domaine"] = e.Text
  })

  // matching Fonction
  c.OnHTML("#extraQuestionName > li:nth-child(2) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Fonction"] = e.Text
  })

  // matching Entreprise
  c.OnHTML("#extraQuestionName > li:nth-child(4) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Entreprise"] = e.Text
  })

  // matching Contrat
  c.OnHTML("#extraQuestionName > li:nth-child(3) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Contrat"] = e.Text
  })

  // matching Niveau d'études
  c.OnHTML("#extraQuestionName > li:nth-child(6) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Niveau d'études"] = e.Text
  })

  // matching Salaire
  c.OnHTML("#extraQuestionName > li:nth-child(5) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*buffer)[e.Request.Ctx.Get("title")]["Salaire"] = e.Text
  })
}

