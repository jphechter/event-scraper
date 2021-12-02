package venue

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

// Scrape Ram's Head Live
func ScrapeRH(db *gorm.DB, wg *sync.WaitGroup) {
	c := colly.NewCollector()

	var v Venue
	db.First(&v, 1) // Ram's Head

	// Find relevant information
	c.OnHTML(".entry.ramsheadlive", func(e *colly.HTMLElement) {
		// Format Date & Time
		d, err := time.Parse("Mon, Jan 2, 2006", e.ChildText(".date"))

		// Clean time
		re := regexp.MustCompile(`((1[0-2]|0?[1-9]):([0-5][0-9]) ?([AaPp][Mm]))`)
		tString := re.FindString(e.ChildText(".time")) // Clean time
		t, _ := time.Parse("3:04 PM", tString)

		// Consolidate
		loc, _ := time.LoadLocation("EST")
		dt := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, loc)

		if err != nil {
			log.Printf("\u001b[31mERROR:\u001b[0m Could not parse date, err :%q", err)
		} else {
			db.Create(&Event{
				Name:      e.ChildText(".carousel_item_title_small"),           // event name
				Date:      dt,                                                  // date
				EventPage: e.ChildAttr(".carousel_item_title_small a", "href"), // event link
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

	c.Visit("https://www.ramsheadlive.com/events/all")
	wg.Done()
}
