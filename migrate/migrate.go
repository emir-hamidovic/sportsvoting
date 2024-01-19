package migrate

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	migrateV4 "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Running migrations")

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	filePath := "file://./migrations/"
	isDev := true
	if isdevEnv, exists := os.LookupEnv("IS_DEVELOPMENT"); exists {
		isDev, _ = strconv.ParseBool(isdevEnv)
	}

	if !isDev {
		filePath = "file:///migrations"
	}

	m, err := migrateV4.NewWithDatabaseInstance(
		filePath,
		"mysql", driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrateV4.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully")
	return nil
}
