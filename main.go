package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"scraper/database"
	"scraper/parser/players"
	"scraper/parser/teams"
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

/*
func RunUpdate(db database.Database, ctx context.Context, errCh chan<- error) {
	ticker := time.NewTicker(24 * time.Hour)
	for {
		select {
		case <-ctx.Done():
			return // if the context is cancelled, return without sending the error
		case <-ticker.C:
			now := time.Now().UTC()
			if now.Hour() == 8 && now.Minute() == 0 {
				err := UpdatePlayerStats(db)
				if err != nil {
					errCh <- err
					return
				}
			}
		}
	}
}*/

func DPOYAward(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetDPOYStats(ctx, players.GetEndYearOfTheSeason())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []players.Player
	for rows.Next() {
		var p players.Player
		err := rows.Scan(&p.Name, &p.Games, &p.Minutes, &p.Rebounds, &p.Steals, &p.Blocks, &p.Position, &p.DefWS, &p.DefBPM, &p.DefRtg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
}

func MIPAward(w http.ResponseWriter, r *http.Request) {
	MVPAward(w, r)
}

func ROYAward(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetROYStats(ctx, players.GetEndYearOfTheSeason())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []players.Player
	for rows.Next() {
		var p players.Player
		err := rows.Scan(&p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.WS, &p.BPM, &p.OffRtg, &p.DefRtg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
}

func SixManAward(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetSixManStats(ctx, players.GetEndYearOfTheSeason())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []players.Player
	for rows.Next() {
		var p players.Player
		err := rows.Scan(&p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
}

func MVPAward(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetMVPStats(ctx, players.GetEndYearOfTheSeason())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []players.Player
	for rows.Next() {
		var p players.Player
		err := rows.Scan(&p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
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
	}
	*/

	err = players.UpdatePlayersWhoPlayedAGame(db)
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	r.HandleFunc("/sixthman", SixManAward)
	r.HandleFunc("/dpoy", DPOYAward)
	r.HandleFunc("/mvp", MVPAward)
	r.HandleFunc("/mip", MIPAward)
	r.HandleFunc("/roy", ROYAward)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Starting server")
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	/*
		ctx, cancel := context.WithCancel(context.Background())
		errChan := make(chan error, 1)
		go RunUpdate(db, ctx, errChan)

		go func(errCh <-chan error) {
			for err := range errCh {
				if err != nil {
					log.Println(err)
					cancel() // cancel the context if an error is received
				}
			}
		}(errChan)

		<-ctx.Done()
		close(errChan)*/
}

// What if a player is added, we are doing an update, never an insert after the first insert
// what if he wasnt in that initial sync, trades, free agent signing etc.
