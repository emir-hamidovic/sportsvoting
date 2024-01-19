package main

import (
	"log"
	"os"
	"sportsvoting/database"
	"sportsvoting/http"
	"sportsvoting/migrate"
	"sportsvoting/syncer"
	"time"
)

var db database.Database

func main() {
	db_addr := os.Getenv("DBADDRESS")
	if db_addr == "" {
		db_addr = "localhost:3306"
	}

	var err error
	retries := 10
	delayTime := 5 * time.Second

	for i := 0; i < retries; i++ {
		log.Println("Attempting to connect to database")

		db, err = database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: db_addr})
		if err == nil {
			break
		}

		log.Printf("Couldn't connect to database due to: %v\n", err)

		if i < retries {
			log.Printf("Retrying in %v ... \n", delayTime)
			time.Sleep(delayTime)
		}
	}

	if db == nil {
		log.Fatalf("Couldn't connect to database after %d retries, exiting!\n", retries)
	}

	defer db.CloseConnection()

	if err := migrate.RunMigrations(db.GetDB()); err != nil {
		log.Fatalf("Error running migrations: %v", err)
	}

	err = db.CreateAdminUser()

	if err != nil {
		log.Fatalf("Error creating admin user: %v", err)
	}

	syncer.SyncRegular(db)
	syncer.SyncGOAT(db)
	syncer.SetupSyncSchedules(db)

	http.StartServer(db)
}
