package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Establish a new database connection
func NewDB() *gorm.DB {
	fmt.Println("Connecting to MySQL database...")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			IgnoreRecordNotFoundError: true, // Mute Record Not Found error
		})

	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		fmt.Println("Unable to connect to database", err.Error())
		return nil
	}

	fmt.Println("Database connected!")

	return db
}
