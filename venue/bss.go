package venue

import (
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gocolly/colly"
)

// Scrape Baltimore Sound Stage
func ScrapeBSS(w *csv.Writer, wg *sync.WaitGroup) {
	c := colly.NewCollector()

	// Find relevant information
	c.OnHTML(".event", func(e *colly.HTMLElement) {
		w.Write([]string{
			e.ChildText("span.title"),   // event name
			"",                          // date (with time)
			e.ChildText(".event-time"),  // time
			e.ChildAttr("h2 a", "href"), // event link
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	for i := 1; i <= 3; i++ {
		c.Visit("https://www.baltimoresoundstage.com/events-feed/page/" + strconv.Itoa(i))
	}
	wg.Done()
}
