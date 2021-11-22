package event

import (
	"time"

	"github.com/event-scraper/venue"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name      string
	Date      time.Time
	EventPage string
	VenueID   int
	Venue     venue.Venue
}

func GenerateSeedEvents(db *gorm.DB) {
	db.AutoMigrate(&Event{})

	var v venue.Venue
	db.First(&v, 1) // Ram's Head
	t, _ := time.Parse("2006-01-02", "2021-11-22")
	db.Create(&Event{
		Name:      "Event 1",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-1",
		VenueID:   int(v.ID), // There are 2 valid ways to establish this key
	})

	db.Create(&Event{
		Name:      "Event 2",
		Date:      t,
		EventPage: "https://www.ramsheadlive.com/event-2",
		Venue:     v,
	})

	var w venue.Venue
	db.First(&w, 2) // Soundstage
	db.Create(&Event{
		Name:      "Event 3",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-3",
		VenueID:   int(w.ID),
	})

	db.Create(&Event{
		Name:      "Event 4",
		Date:      t,
		EventPage: "https://www.baltimoresoundstage.com/event-4",
		Venue:     w,
	})
}
