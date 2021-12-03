# Event Scraper

## Backstory
The general idea for this app came a number of years ago when I found that I was missing interesting events in my city. I found it frustrating to check museum pages, music venues, and bars individually to find the things I was most interested in. There are lots of aveunes to find a more complete list, no one really has everything or provides a meaningful way to sort the information.

## Setup
👷🏻‍♂️ 🚧 
Instructions to come... If someone actually needs help running this you probably know me and you can just ping me.

## DB Migrations
Generating and running migrations can be handled through a CLI. The root command to invoke the migration CLI is: `go run main.go migrate`

### Generate New Migration
All migrations are based on a `template.txt` and are stored in ./database/migrations/

`go run main.main.go migrate create -n migration_name`

### Run Migration Up/Down
Run all migrations up or down.
`go run main.go migrate up`
`go run main.go migrate down`

Use the step flag to limit the number of migrations run
`go run main.go migrate up --step 1`
`go run main.go migrate down --step 1`

### Check Current DB Version
Returns the DB Version # which corresponds to the last migration run.
`go run main.go migrate status`