package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: Rewrite for GORM
// Establish a new database connection
func NewDB() *sql.DB {
	fmt.Println("Connecting to MySQL database...")

	dsn := os.Getenv("ES_DATABASE_URL")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Unable to connect to database", err.Error())
		return nil
	}

	if err := db.Ping(); err != nil {
		fmt.Println("Unable to connect to database", err.Error())
		return nil
	}

	fmt.Println("Database connected!")

	return db
}
