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

	"github.com/gorilla/handlers"
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

	/*rosters, err := InsertTeamAndPlayerInfo(db)
	if err != nil {
		log.Fatal(err)
	}

	err = players.UpdatePlayerStats(db, rosters)
	if err != nil {
		log.Fatal(err)
	}*/

	// Need a better solution for inserting pictures
	polls := []Poll{
		{1, "MVP", "Description for MVP", "https://arc-anglerfish-washpost-prod-washpost.s3.amazonaws.com/public/ELCVCAZWP5F5LPYIILNXKWO5FU.jpg", "mvp"},
		{2, "ROY", "Description for ROY", "https://pbs.twimg.com/media/Fj26PZrXkAAGPo4?format=jpg&name=4096x4096", "roy"},
		{3, "DPOY", "Description for DPOY", "https://pbs.twimg.com/media/Fj26CtNXwAESyDv?format=jpg&name=4096x4096", "dpoy"},
		{4, "Sixth Man", "Description for 6-man", "https://pbs.twimg.com/media/Fj26barWYAMeoZp?format=jpg&name=4096x4096", "sixthman"},
	}

	for _, val := range polls {
		db.InsertPolls(val.ID, val.Name, val.Description, val.Image, val.Endpoint)
	}

	r := mux.NewRouter()

	// Will need a custom endpoint
	r.HandleFunc("/getpolls", getPolls)
	r.HandleFunc("/sixthman", sixManAward)
	r.HandleFunc("/dpoy", dpoyAward)
	r.HandleFunc("/mvp", mvpAward)
	r.HandleFunc("/mip", mipAward)
	r.HandleFunc("/roy", royAward)
	r.HandleFunc("/teamvotes/{id:[0-9]+}", teamVotes)
	r.HandleFunc("/playervotes/{id:[0-9]+}", playerVotes).Methods("GET")
	r.HandleFunc("/playervotes/", insertPlayerVotes).Methods("POST")
	r.HandleFunc("/playervotes/", handleOptions).Methods("OPTIONS")
	r.HandleFunc("/playervotes/", insertPlayerVotes).Methods("POST")
	r.HandleFunc("/login", handleLogin).Methods("POST")
	r.HandleFunc("/register", handleRegister).Methods("POST")
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/refresh", handleRefresh)

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)

	http.Handle("/", corsHandler(r))

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
