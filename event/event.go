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
