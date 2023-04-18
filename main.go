package main

import (
	"fmt"
	"log"
	"net/http"
	"scraper/database"
	"scraper/parser/players"
	"scraper/parser/stats"
	"scraper/parser/teams"
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

func UpdatePlayerStats(db database.Database, rosters map[string]players.Player) error {
	fmt.Println("Updating stats")
	for _, player := range rosters {
		err := stats.UpdateStats(db, player.Stats)
		if err != nil {
			return err
		}
	}

	return nil
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

func MVPAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func DPOYAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func MIPAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func COYAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func ROYAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func SixManAward(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Simple Server")
}

func main() {
	db, err := database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: "localhost:3306"})
	if err != nil {
		log.Fatal(err)
	}

	rosters, err := InsertTeamAndPlayerInfo(db)
	if err != nil {
		log.Fatal(err)
	}

	err = UpdatePlayerStats(db, rosters)
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

// Create HTTP handlers for each separate award available: lets start with regular awards like MVP, MIP, DPOY (need to get advanced stats as well),
// COY, 6MOY etc.

// What if a player is added, we are doing an update, never an insert after the first insert
// what if he wasnt in that initial sync, trades, free agent signing etc.

// need to insert stats
