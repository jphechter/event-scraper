package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// TODO: Rewrite for GORM
// Establish a new database connection
func NewDB() *gorm.DB {
	fmt.Println("Connecting to MySQL database...")

	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Unable to connect to database", err.Error())
		return nil
	}

	fmt.Println("Database connected!")

	return db
}
