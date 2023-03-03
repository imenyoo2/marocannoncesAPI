package main

import (
	"fmt"
	"log"
  "time"
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
  c := colly.NewCollector()
  // setting callback functions
  c.OnRequest(func(r *colly.Request) {
    fmt.Println("Visiting", r.URL)
  })
  c.OnError(func(_ *colly.Response, err error) {
    log.Println("Something went wrong:", err)
  })
  c.OnResponse(func(r *colly.Response){
    fmt.Println("Visited", r.Request.URL)
  })
  c.OnHTML(".listing .browsing_result_table_body_even a[href]", func(e *colly.HTMLElement) {
    e.Request.Visit(e.Attr("href"))
  })
  c.OnHTML(".infoannonce dl dd", func(e *colly.HTMLElement) {
    fmt.Println(e.Text)
  })
  c.SetRequestTimeout(1 * time.Minute)
  c.Visit("https://www.marocannonces.com/categorie/309/Emploi/Offres-emploi.html")

}
