package venue

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

// Scrape Baltimore Sound Stage
func ScrapeBSS(db *gorm.DB, wg *sync.WaitGroup) {
	c := colly.NewCollector()

	var v Venue
	db.First(&v, 2) // BSS

	// Find relevant information
	c.OnHTML(".event", func(e *colly.HTMLElement) {
		// Format Date & Time
		d, err := time.Parse("Monday, January 2, 2006", e.ChildText(".event-date"))

		// Clean time
		re := regexp.MustCompile(`((1[0-2]|0?[1-9]):([0-5][0-9]) ?([AaPp][Mm]))`)
		tString := re.FindString(e.ChildText(".event-time"))
		t, _ := time.Parse("3:04 PM", tString)

		// Consolidate
		loc, _ := time.LoadLocation("EST")
		dt := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, loc)

		if err != nil {
			log.Printf("\u001b[31mERROR:\u001b[0m Could not parse date, err :%q", err)
		} else {
			db.Create(&Event{
				Name:      e.ChildText("span.title"),   // event name
				Date:      dt,                          // date
				EventPage: e.ChildAttr("h2 a", "href"), // event link
				VenueID:   int(v.ID),
			})
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// TODO: Progromatically check for num of pages
	for i := 1; i <= 3; i++ {
		c.Visit("https://www.baltimoresoundstage.com/events-feed/page/" + strconv.Itoa(i))
	}
	wg.Done()
}
