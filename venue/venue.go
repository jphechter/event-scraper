package venue

import (
	"gorm.io/gorm"
)

type Venue struct {
	gorm.Model
	Name    string
	Website string
	Address string
}

func GenerateSeedVenues(db *gorm.DB) {
	db.AutoMigrate(&Venue{})

	db.Create(&Venue{
		Name:    "Ram's Head Live",
		Website: "https://www.ramsheadlive.com/",
		Address: "20 Market Pl Baltimore, MD 21202",
	})

	db.Create(&Venue{
		Name:    "Baltimore Sound Stage",
		Website: "https://www.baltimoresoundstage.com/",
		Address: "124 Market Pl Baltimore, MD 21202",
	})
}
