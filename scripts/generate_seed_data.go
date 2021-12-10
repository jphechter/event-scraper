// Generate Seed Data creates new database and populates with seed values
package main

import (
	"fmt"
	"time"

	"github.com/event-scraper/database"
	"github.com/event-scraper/venue"
	"gorm.io/gorm"
)

func main() {
	db := database.NewDB()
	// fmt.Println("DROP DATABASE")
	// db.Exec("DROP DATABASE event_scraper")
	// fmt.Println("CREATE DATABASE")
	// db.Exec("CREATE DATABASE event_scraper")
	// db = createDBConnection() // Reconnect to new database

	generateVenues(db)
	generateEvents(db)
	generateVenueEventRules(db)

	// Prove that a venue was populated
	var v venue.Venue
	db.First(&v)
	fmt.Println("\ncreated: ", v.CreatedAt)
	fmt.Println("id: ", v.ID)
	fmt.Println("name: ", v.Name)
	fmt.Println("website: ", v.Website)
}

func generateVenues(db *gorm.DB) {
	fmt.Println("Generating Venues")
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
	fmt.Println("Generating Events")
	rh := venue.GetVenueByID(db, 1) // Ram's Head
	t, _ := time.Parse("2006-01-02", "2021-11-22")
	db.Create(&venue.Event{
		Name:      "Event 1",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-1",
		VenueID:   int(rh.ID), // There are 2 valid ways to establish this key
	})

	db.Create(&venue.Event{
		Name:      "Event 2",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-2",
		Venue:     rh,
	})

	bss := venue.GetVenueByID(db, 2) // Baltimore Soundstage
	db.Create(&venue.Event{
		Name:      "Event 3",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-3",
		VenueID:   int(bss.ID),
	})

	db.Create(&venue.Event{
		Name:      "Event 4",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-4",
		Venue:     bss,
	})
}

func generateVenueEventRules(db *gorm.DB) {
	fmt.Println("Generating Venue Event Rules")
	rh := venue.GetVenueByID(db, 1) // Ram's Head
	rh.SetVenueEventRule(db, venue.CalendarURL, "https://www.ramsheadlive.com/events/all")
	rh.SetVenueEventRule(db, venue.EventSelector, ".entry.ramsheadlive")
	rh.SetVenueEventRule(db, venue.Date, ".date")
	rh.SetVenueEventRule(db, venue.DateFormat, "Mon, Jan 2, 2006")
	rh.SetVenueEventRule(db, venue.Time, ".time")
	rh.SetVenueEventRule(db, venue.TimeFormat, "3:04 PM")
	rh.SetVenueEventRule(db, venue.EventName, ".carousel_item_title_small")
	rh.SetVenueEventRule(db, venue.EventPage, ".carousel_item_title_small a")

	bss := venue.GetVenueByID(db, 2) // Baltimore Soundstage
	bss.SetVenueEventRule(db, venue.CalendarURL, "https://www.baltimoresoundstage.com/events-feed/page/%d")
	bss.SetVenueEventRule(db, venue.EventSelector, ".event")
	bss.SetVenueEventRule(db, venue.Date, ".event-date")
	bss.SetVenueEventRule(db, venue.DateFormat, "Monday, January 2, 2006")
	bss.SetVenueEventRule(db, venue.Time, ".event-time")
	bss.SetVenueEventRule(db, venue.TimeFormat, "3:04 PM")
	bss.SetVenueEventRule(db, venue.EventName, ".span.title")
	bss.SetVenueEventRule(db, venue.EventPage, "h2 a")
}
