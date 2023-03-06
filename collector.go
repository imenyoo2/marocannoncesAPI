package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly"
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
    collectPage(collectPageParams{c: c, e: e, url: e.Attr("href"), dataStore: app.data})
  })

  // matching daily posts
  c.OnHTML(".cars-list > li", func(e *colly.HTMLElement) {
    time := e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(1)") + " " + e.ChildText("div:nth-child(2) > em:nth-child(1) > span:nth-child(3)")
    collectPage(collectPageParams{c: c, e: e, url: e.ChildAttr("a:nth-child(1)", "href"), dataStore: app.data, time: time})
  })

  c.SetRequestTimeout(1 * time.Minute)

  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

}

type collectPageParams struct {
  c *colly.Collector
  e *colly.HTMLElement
  url string
  dataStore *map[string]map[string]interface{}
  time string
}
func collectPage(params collectPageParams) {

  params.e.Request.Visit(params.url)

  // matching the title
  params.c.OnHTML("#content > div.used-cars > div.description.desccatemploi > h1", func(e *colly.HTMLElement) {
    title := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "  ", "")
    e.Request.Ctx.Put("title", title)
    (*params.dataStore)[title] = map[string]interface{}{"title": title}
    // adding the url field
    (*params.dataStore)[e.Request.Ctx.Get("title")]["URL"] = params.url

    // adding time if exist
    if params.time != "" {
      (*params.dataStore)[e.Request.Ctx.Get("title")]["time"] = params.time
    }
  })


  // matching Annonceur
  params.c.OnHTML(".infoannonce > dl:nth-child(1) > dd:nth-child(2)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Annonceur"] = e.Text
  })

  // matching Domaine
  params.c.OnHTML("#extraQuestionName > li:nth-child(1) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Domaine"] = e.Text
  })

  // matching Fonction
  params.c.OnHTML("#extraQuestionName > li:nth-child(2) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Fonction"] = e.Text
  })

  // matching Entreprise
  params.c.OnHTML("#extraQuestionName > li:nth-child(4) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Entreprise"] = e.Text
  })

  // matching Contrat
  params.c.OnHTML("#extraQuestionName > li:nth-child(3) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Contrat"] = e.Text
  })

  // matching Niveau d'études
  params.c.OnHTML("#extraQuestionName > li:nth-child(6) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Niveau d'études"] = e.Text
  })

  // matching Salaire
  params.c.OnHTML("#extraQuestionName > li:nth-child(5) > a:nth-child(1)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*params.dataStore)[e.Request.Ctx.Get("title")]["Salaire"] = e.Text
  })
}
