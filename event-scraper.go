// Event Scraper is a periodic task that looks for upcoming events
// and stores the relevant information in a database.
package main

import (
	"fmt"
	"os"

	"github.com/event-scraper/event"
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
	if db == nil {
		panic("failed to connect database")
	}

	venue.GenerateSeedVenues(db)
	event.GenerateSeedEvents(db)
	var v venue.Venue
	db.First(&v)

	fmt.Println("\n\nresult: ", v)
	fmt.Println("\ncreated: ", v.CreatedAt)
	fmt.Println("\nid: ", v.ID)
	fmt.Println("\nname: ", v.Name)
	fmt.Println("\nwebsite: ", v.Website)

	// Storing to file in the interim
	// file, err := os.Create("data.csv")
	// if err != nil {
	// 	log.Fatalf("Could not create file, err :%q", err)
	// 	return
	// }
	// writer := csv.NewWriter(file)

	// venues := []func(*csv.Writer, *sync.WaitGroup){venue.ScrapeRH, venue.ScrapeBSS}

	// var wg *sync.WaitGroup = new(sync.WaitGroup)
	// for _, venue := range venues {
	// 	wg.Add(1)
	// 	go venue(writer, wg)
	// }

	// wg.Wait()
	// file.Close()
	// writer.Flush()
}
