// Event Scraper is a periodic task that looks for upcoming events
// and stores the relevant information in a database.
package main

import (
	"os"
	"sync"

	"github.com/event-scraper/cmd"
	"github.com/event-scraper/venue"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	if len(os.Args[1:]) > 0 { // Currently only handles migrations
		cmd.Execute()
		os.Exit(0)
	}

	// NOTE: the args following db_name are required to establish the connection correctly
	// dsn = "user_name:password@tcp(localhost:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Scrape venues
	var wg *sync.WaitGroup = new(sync.WaitGroup)
	// venues := getAllVenues(db)
	// for _, v := range venues {
	rh := venue.GetVenueByID(db, 1)
	go rh.ScrapeVenue(db, wg)
	// go v.ScrapeVenue(db, wg)
	wg.Add(1)
	// }

	wg.Wait()
}

func getAllVenues(db *gorm.DB) []venue.Venue {
	var venues []venue.Venue
	db.Find(&venues)
	return venues
}
