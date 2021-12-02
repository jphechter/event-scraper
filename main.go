// Event Scraper is a periodic task that looks for upcoming events
// and stores the relevant information in a database.
package main

import (
	"os"
	"sync"

	"github.com/event-scraper/venue"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// NOTE: the args following db_name are required to establish the connection correctly
	// dsn = "user_name:password@tcp(localhost:3306)/db_name?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	venues := []func(*gorm.DB, *sync.WaitGroup){venue.ScrapeRH}
	// venues := []func(*csv.Writer, *sync.WaitGroup){venue.ScrapeRH, venue.ScrapeBSS}

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	for _, venue := range venues {
		wg.Add(1)
		go venue(db, wg)
	}

	wg.Wait()
}
