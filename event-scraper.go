// Event Scraper is a periodic task that looks for upcoming events
// and stores the relevant information in a database.
package main

import (
	"encoding/csv"
	"log"
	"os"
	"sync"

	"github.com/event-scraper/venue"
)

func main() {
	// Storing to file in the interim
	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatalf("Could not create file, err :%q", err)
		return
	}
	writer := csv.NewWriter(file)

	venues := []func(*csv.Writer, *sync.WaitGroup){venue.ScrapeRH, venue.ScrapeBSS}

	var wg *sync.WaitGroup = new(sync.WaitGroup)
	for _, venue := range venues {
		wg.Add(1)
		go venue(writer, wg)
	}

	wg.Wait()
	file.Close()
	writer.Flush()
}
