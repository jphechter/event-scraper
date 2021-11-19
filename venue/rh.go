package venue

import (
	"encoding/csv"
	"fmt"
	"log"
	"sync"

	"github.com/gocolly/colly"
)

// Scrape Ram's Head Live
func ScrapeRH(w *csv.Writer, wg *sync.WaitGroup) {
	c := colly.NewCollector()

	// Find relevant information
	c.OnHTML(".entry.ramsheadlive", func(e *colly.HTMLElement) {
		w.Write([]string{
			e.ChildText(".carousel_item_title_small"),           // event name
			e.ChildText(".date"),                                // date
			e.ChildText(".time"),                                // time
			e.ChildAttr(".carousel_item_title_small a", "href"), // event link
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	c.Visit("https://www.ramsheadlive.com/events/all")
	wg.Done()
}
