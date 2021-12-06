package venue

import (
	"time"

	"gorm.io/gorm"
)

type Venue struct {
	Name    string
	Website string
	Address string
	gorm.Model
}

type Event struct {
	Name      string
	Date      time.Time
	EventPage string
	VenueID   int
	Venue     Venue
	gorm.Model
}
