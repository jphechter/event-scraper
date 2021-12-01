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
