package main

import (
	"errors"
	"log"
	"scraper/database"
	"scraper/parser/players"
	"scraper/parser/teams"
)

func Parse(db database.Database, parse string) error {
	if parse == "teams" {
		return teams.ParseTeams(db)
	} else if parse == "players" {
		return players.ParsePlayers(db)
	}

	return errors.New("incorrect type")
}

func main() {
	db, err := database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: "localhost:3306"})
	if err != nil {
		log.Fatal(err)
	}

	err = Parse(db, "teams")
	if err != nil {
		log.Fatalf("%v/n", err)
	}

	err = Parse(db, "players")
	if err != nil {
		log.Fatalf("%v/n", err)
	}
}
