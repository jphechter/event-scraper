package migrations

import (
	"github.com/event-scraper/venue"
	"gorm.io/gorm"
)

func init() {
	migrator.AddMigration(&Migration{
		Version: "000001",
		Up:      mig_000001_init_schema_up,
		Down:    mig_000001_init_schema_down,
	})
}

func mig_000001_init_schema_up(db *gorm.DB) error {
	db.Migrator().CreateTable(&venue.Venue{})
	db.Migrator().CreateTable(&venue.Event{})
	db.Migrator().CreateTable(&venue.VenueEventRule{})
	return nil
}

func mig_000001_init_schema_down(db *gorm.DB) error {
	db.Migrator().DropTable(&venue.VenueEventRule{})
	db.Migrator().DropTable(&venue.Event{})
	db.Migrator().DropTable(&venue.Venue{})
	return nil
}
