package main

import (
	"fmt"
	"log"
  "time"
  "strings"
  "encoding/json"
  //"os"

	"github.com/gocolly/colly"
)

func check(err error) {
  if err != nil {
    log.Fatal(err)
  }
}


func main() {
  // initializing a new collecotor
  data := map[string]map[string]interface{}{}
  c := colly.NewCollector()
  // setting callback functions
  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
  })
  c.OnError(func(_ *colly.Response, err error) {
    log.Println("Something went wrong:", err)
  })
  c.OnResponse(func(r *colly.Response){
    //fmt.Println("Visited", r.Request.URL)
  })

  c.OnHTML("article.listing > a:nth-child(1)", func(e *colly.HTMLElement) {
    e.Request.Visit(e.Attr("href"))
  })

  // matching the title
  c.OnHTML("#content > div.used-cars > div.description.desccatemploi > h1", func(e *colly.HTMLElement) {
    //fmt.Println("getting title")
    title := strings.ReplaceAll(strings.ReplaceAll(e.Text, "\n", ""), "  ", "")
    e.Request.Ctx.Put("title", title)
    //fmt.Println("title added to context")
    data[title] = map[string]interface{}{"title": title}
    //fmt.Println("title added to data")
  })
  // matching Annonceur
  c.OnHTML(".infoannonce > dl:nth-child(1) > dd:nth-child(2)", func(e *colly.HTMLElement) {
    //fmt.Println("getting annonceur")
    data[e.Request.Ctx.Get("title")]["Annonceur"] = e.Text
    //fmt.Println("annonceur added to data")
  })
  c.SetRequestTimeout(1 * time.Minute)
  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")
  jsonStr, err := json.Marshal(data)
  check(err)
  fmt.Println(string(jsonStr))

}
