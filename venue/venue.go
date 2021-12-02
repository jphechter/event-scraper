package venue

import (
	"time"

	"gorm.io/gorm"
)

type Venue struct {
	gorm.Model
	Name    string
	Website string
	Address string
}

type Event struct {
	gorm.Model
	Name      string
	Date      time.Time
	EventPage string
	VenueID   int
	Venue     Venue
}
