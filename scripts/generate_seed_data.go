// Generate Seed Data creates new database and populates with seed values
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/event-scraper/event"
	"github.com/event-scraper/venue"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db := createDBConnection()
	fmt.Println("DROP DATABASE")
	db.Exec("DROP DATABASE event_scraper")
	fmt.Println("CREATE DATABASE")
	db.Exec("CREATE DATABASE event_scraper")
	db = createDBConnection() // Reconnect to new database

	generateVenues(db)
	generateEvents(db)

	// Prove that a venue was populated
	var v venue.Venue
	db.First(&v)
	fmt.Println("\ncreated: ", v.CreatedAt)
	fmt.Println("id: ", v.ID)
	fmt.Println("name: ", v.Name)
	fmt.Println("website: ", v.Website)
}

func createDBConnection() *gorm.DB {
	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func generateVenues(db *gorm.DB) {
	db.AutoMigrate(&venue.Venue{})

	db.Create(&venue.Venue{
		Name:    "Ram's Head Live",
		Website: "https://www.ramsheadlive.com/",
		Address: "20 Market Pl Baltimore, MD 21202",
	})

	db.Create(&venue.Venue{
		Name:    "Baltimore Sound Stage",
		Website: "https://www.baltimoresoundstage.com/",
		Address: "124 Market Pl Baltimore, MD 21202",
	})
}

func generateEvents(db *gorm.DB) {
	db.AutoMigrate(&event.Event{})

	var v venue.Venue
	db.First(&v, 1) // Ram's Head
	t, _ := time.Parse("2006-01-02", "2021-11-22")
	db.Create(&event.Event{
		Name:      "Event 1",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-1",
		VenueID:   int(v.ID), // There are 2 valid ways to establish this key
	})

	db.Create(&event.Event{
		Name:      "Event 2",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-2",
		Venue:     v,
	})

	var w venue.Venue
	db.First(&w, 2) // Soundstage
	db.Create(&event.Event{
		Name:      "Event 3",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-3",
		VenueID:   int(w.ID),
	})

	db.Create(&event.Event{
		Name:      "Event 4",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-4",
		Venue:     w,
	})
}
