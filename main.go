package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sportsvoting/database"
	"sportsvoting/players"
	"sportsvoting/teams"
	"time"

	"github.com/gorilla/mux"
)

func InsertTeamAndPlayerInfo(db database.Database) (map[string]players.Player, error) {
	fmt.Println("Parsing teams")
	playerList, err := teams.ParseTeams(db)
	if err != nil {
		return nil, err
	}

	fmt.Println("Inserting players")
	err = players.InsertPlayers(db, playerList)
	if err != nil {
		return nil, err
	}

	return playerList, err
}

func RunUpdate(db database.Database, ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ctx.Done():
			log.Println("Context cancelled. Goroutine exiting")
			return // if the context is cancelled, return without sending the error
		case <-ticker.C:
			now := time.Now().UTC()
			if now.Hour() == 8 && now.Minute() == 0 {
				err := players.UpdatePlayersWhoPlayedAGame(db)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}

var db database.Database

func main() {
	var err error
	db, err = database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: "localhost:3306"})
	if err != nil {
		log.Fatal(err)
	}

	rosters, err := InsertTeamAndPlayerInfo(db)
	if err != nil {
		log.Fatal(err)
	}

	err = players.UpdatePlayerStats(db, rosters)
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/getpolls", getPolls)
	r.HandleFunc("/sixthman", sixManAward)
	r.HandleFunc("/dpoy", dpoyAward)
	r.HandleFunc("/mvp", mvpAward)
	r.HandleFunc("/mip", mipAward)
	r.HandleFunc("/roy", royAward)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go RunUpdate(db, ctx)

	fmt.Println("Starting server")
	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	cancel()
}
