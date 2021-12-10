package venue

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"gorm.io/gorm"
)

// Events are the reason we're here
type Event struct {
	ID        int
	Name      string
	Date      time.Time
	EventPage string
	VenueID   int
	Venue     Venue
	gorm.Model
}

// A Venue is where an Event takes place
type Venue struct {
	ID      int
	Name    string
	Website string
	Address string
	gorm.Model
}

// Rules for collecting data about an Event from a Venue website
type VenueEventRule struct {
	ID      int
	VenueID int
	Venue   Venue
	Type    string
	Info    string
	gorm.Model
}

// Venue Event Rule types
const (
	CalendarURL   = "calendar base url"
	EventSelector = "event selector"
	Date          = "date selector"
	DateFormat    = "date format"
	Time          = "time selector"
	TimeFormat    = "time format"
	EventName     = "event name selector"
	EventPage     = "event page url selector"
)

func GetVenueByID(db *gorm.DB, id int) Venue {
	var v Venue
	db.First(&v, id)
	return v
}

func GetVenueByName(db *gorm.DB, name string) Venue {
	var v Venue
	db.Where("name = ?", name).First(&v)
	return v
}

func (v Venue) GetVenueEventRule(db *gorm.DB, prop string) VenueEventRule {
	var ver VenueEventRule
	db.Where("type = ?", prop).Where("venue_id = ?", v.ID).First(&ver)
	return ver
}

func (v Venue) SetVenueEventRule(db *gorm.DB, prop string, val string) {
	var ver VenueEventRule
	// Delete any old rule before creating a new one
	err := db.Where("type = ?", prop).Where("venue_id = ?", v.ID).First(&ver).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Print("Existing Rule: ", ver)
		db.Delete(&ver)
	}

	db.Create(&VenueEventRule{
		Venue: v,
		Type:  prop,
		Info:  val,
	})
}

func (v Venue) ScrapeVenue(db *gorm.DB, wg *sync.WaitGroup) {
	// Initialize Colly
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) { fmt.Println("Visiting", r.URL) })
	c.OnError(func(_ *colly.Response, err error) { log.Println("Something went wrong:", err) })

	// Find relevant event information
	c.OnHTML(
		v.GetVenueEventRule(db, EventSelector).Info,
		func(e *colly.HTMLElement) {

			dt, err := v.scrapeEventDate(db, e)
			newEvent := Event{
				Name:      e.ChildText(v.GetVenueEventRule(db, EventName).Info),         // event name
				Date:      dt,                                                           // date
				EventPage: e.ChildAttr(v.GetVenueEventRule(db, EventPage).Info, "href"), // event link
				VenueID:   int(v.ID),
			}
			if err == nil {
				db.Create(&newEvent)
			}
		})

	// TODO: Account for pagination
	c.Visit(v.GetVenueEventRule(db, CalendarURL).Info)
	wg.Done()
}

func (v Venue) scrapeEventDate(db *gorm.DB, e *colly.HTMLElement) (time.Time, error) {
	// Format Date & Time
	d, err := time.Parse(
		v.GetVenueEventRule(db, DateFormat).Info,
		e.ChildText(v.GetVenueEventRule(db, Date).Info),
	)
	if err != nil {
		log.Printf("\u001b[31mERROR:\u001b[0m Could not parse date, err :%q", err)
	}

	// Check if separate Time Rules exist
	if timeRule := v.GetVenueEventRule(db, Time); timeRule != (VenueEventRule{}) {
		// Clean time
		re := regexp.MustCompile(`((1[0-2]|0?[1-9]):([0-5][0-9]) ?([AaPp][Mm]))`)
		tString := re.FindString(e.ChildText(".time")) // Clean time
		t, err := time.Parse("3:04 PM", tString)
		if err != nil {
			log.Printf("\u001b[31mERROR:\u001b[0m Could not parse time, err :%q", err)
		}

		// Consolidate
		loc, _ := time.LoadLocation("EST")
		dt := time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), 0, 0, loc)
		return dt, err
	}
	// TODO: Returning here doesn't set the correct Time Zone
	// Need to account for cases where Day Hour Min not set
	return d, err
}
