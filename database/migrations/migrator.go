package migrations

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gorm.io/gorm"
)

type Migration struct {
	Version string
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error

	done bool
}

type Migrator struct {
	db         *gorm.DB
	Versions   []string
	Migrations map[string]*Migration
}

type DatabaseVersion struct {
	Version string
}

var migrator = &Migrator{
	Versions:   []string{},
	Migrations: map[string]*Migration{},
}

func (m *Migrator) AddMigration(mg *Migration) {
	// Add the migration to the hash with version as key
	m.Migrations[mg.Version] = mg

	// Insert version into versions array using insertion sort
	index := 0
	for index < len(m.Versions) {
		if m.Versions[index] > mg.Version {
			break
		}
		index++
	}

	m.Versions = append(m.Versions, mg.Version)
	copy(m.Versions[index+1:], m.Versions[index:])
	m.Versions[index] = mg.Version
}

// Create new migration from template
func Create(name string, db *gorm.DB) error {
	var last_version DatabaseVersion
	db.Last(&last_version)

	fmt.Printf("Result: %s\n", last_version.Version)

	lv_int, _ := strconv.Atoi(last_version.Version)
	lv_int++ // increment the version
	version := fmt.Sprintf("%06d", lv_int)

	in := struct {
		Version string
		Name    string
	}{
		Version: version,
		Name:    name,
	}

	var out bytes.Buffer

	t := template.Must(template.ParseFiles("./database/migrations/template.txt"))
	err := t.Execute(&out, in)
	if err != nil {
		return errors.New("Unable to execute template:" + err.Error())
	}

	f, err := os.Create(fmt.Sprintf("./database/migrations/%s_%s.go", version, name))
	if err != nil {
		return errors.New("Unable to create migration file:" + err.Error())
	}
	defer f.Close()

	if _, err := f.WriteString(out.String()); err != nil {
		return errors.New("Unable to write to migration file:" + err.Error())
	}

	fmt.Println("Generated new migration files...", f.Name())
	return nil
}

func Init(db *gorm.DB) (*Migrator, error) {
	if !db.Migrator().HasTable(&DatabaseVersion{}) {
		db.Table("_version").AutoMigrate(&DatabaseVersion{})
	}
	migrator.db = db

	// Find executed migrations
	rows, err := db.Raw("SELECT version FROM `_version`;").Rows()
	if err != nil {
		return migrator, err
	}

	defer rows.Close()

	// Mark the migrations as Done if it is already executed
	for rows.Next() {
		var version string
		err := rows.Scan(&version)
		if err != nil {
			return migrator, err
		}

		if migrator.Migrations[version] != nil {
			migrator.Migrations[version].done = true
		}
	}

	return migrator, err
}

// Run Up migrations in a single SQL transaction
// step sets num of migrations
func (m *Migrator) Up(step int) error {
	m.db.Transaction(func(tx *gorm.DB) error {
		count := 0
		for _, v := range m.Versions {
			if step > 0 && count == step {
				break
			}

			mg := m.Migrations[v]

			if mg.done {
				continue
			}

			fmt.Println("Running migration", mg.Version)
			if err := mg.Up(m.db); err != nil {
				tx.Rollback()
				return err
			}

			if result := tx.Exec("INSERT INTO `_version` VALUES(?)", mg.Version); result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			fmt.Println("Finished running migration", mg.Version)

			count++
		}
		return nil
	})
	return nil
}

// Run Down migrations in a single SQL transaction
// step sets num of migrations
func (m *Migrator) Down(step int) error {
	m.db.Transaction(func(tx *gorm.DB) error {
		count := 0
		for _, v := range reverse(m.Versions) {
			if step > 0 && count == step {
				break
			}

			mg := m.Migrations[v]

			if !mg.done {
				continue
			}

			fmt.Println("Reverting Migration", mg.Version)
			if err := mg.Down(tx); err != nil {
				tx.Rollback()
				return err
			}

			if result := tx.Exec("DELETE FROM `schema_migrations` WHERE version = ?", mg.Version); result.Error != nil {
				tx.Rollback()
				return result.Error
			}
			fmt.Println("Finished reverting migration", mg.Version)

			count++
		}
		return nil
	})
	return nil
}

func reverse(arr []string) []string {
	for i := 0; i < len(arr)/2; i++ {
		j := len(arr) - i - 1
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func (m *Migrator) MigrationStatus() error {
	for _, v := range m.Versions {
		mg := m.Migrations[v]

		if mg.done {
			fmt.Println(fmt.Sprintf("Migration %s... completed", v))
		} else {
			fmt.Println(fmt.Sprintf("Migration %s... pending", v))
		}
	}

	return nil
}
