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
    app.collectPage(c, e, e.Attr("href"))
  })


  c.SetRequestTimeout(1 * time.Minute)

  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

}

func (app *application) collectPage(c *colly.Collector, e *colly.HTMLElement, url string) {

  e.Request.Visit(url)

  // matching the title
  c.OnHTML("#content > div.used-cars > div.description.desccatemploi > h1", func(e *colly.HTMLElement) {
    title := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "  ", "")
    e.Request.Ctx.Put("title", title)
    (*app.data)[title] = map[string]interface{}{"title": title}
  })

  // matching Annonceur
  c.OnHTML(".infoannonce > dl:nth-child(1) > dd:nth-child(2)", func(e *colly.HTMLElement) {
    // adding annonceur feild to data
    (*app.data)[e.Request.Ctx.Get("title")]["Annonceur"] = e.Text
  })


}

