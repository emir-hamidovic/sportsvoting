package syncer

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/goatplayers"
	"sportsvoting/players"
	"sportsvoting/teams"
	"time"
)

func InsertTeamAndPlayerInfo(db database.Database, season string) error {
	fmt.Println("Parsing teams")
	playerList, err := teams.ParseTeams(db, season)
	if err != nil {
		return err
	}

	fmt.Println("Inserting players")
	err = players.InsertPlayers(db, playerList)
	if err != nil {
		return err
	}

	err = players.UpdatePlayerStats(db, playerList, season)
	if err != nil {
		return err
	}

	_, err = db.InsertSeasonEntered(season)
	if err != nil {
		return err
	}

	return nil
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

func isSyncNeeded(db database.Database, syncType string) (bool, error) {
	syncTime, err := db.GetLastSyncTime(syncType)
	if err != nil {
		log.Println(err)
		return true, err
	}

	threshold := 10 * 24 * time.Hour
	timeDiff := time.Since(syncTime)

	return timeDiff > threshold, nil
}

func SyncRegular(db database.Database) {
	isSyncNeeded, errSync := isSyncNeeded(db, "Regular")
	if isSyncNeeded {
		err := InsertTeamAndPlayerInfo(db, players.GetEndYearOfTheSeason())
		if err != nil {
			log.Fatal(err)
			return
		}

		if errSync == sql.ErrNoRows {
			db.InsertLastSyncTime(time.Now(), "Regular")
		} else {
			db.UpdateLastSyncTime(time.Now(), "Regular")
		}
	}
}

func SyncGOAT(db database.Database) {
	isSyncNeeded, errSync := isSyncNeeded(db, "GOAT")
	if isSyncNeeded {
		go func() {
			playerIDs := goatplayers.GetGoatPlayersList()
			goatplayers.InsertGoatPlayerStats(playerIDs, db)

			_, err := db.InsertSeasonEntered("All")
			if err != nil {
				log.Println(err)
				return
			}

			_, err = db.InsertSeasonEntered("Playoff")
			if err != nil {
				log.Println(err)
				return
			}

			_, err = db.InsertSeasonEntered("Career")
			if err != nil {
				log.Println(err)
				return
			}

			if errSync == sql.ErrNoRows {
				db.InsertLastSyncTime(time.Now(), "GOAT")
			} else {
				db.UpdateLastSyncTime(time.Now(), "GOAT")
			}
		}()
	}
}

func ScheduleNewSeasonSync(db database.Database) {
	go func() {
		currentDate := time.Now()
		// Calculate the next 1st of November
		nextNovember1st := time.Date(currentDate.Year(), time.November, 1, 0, 0, 0, 0, currentDate.Location())

		for {
			// If the current date is past the next 1st of November, add one year
			if currentDate.After(nextNovember1st) {
				nextNovember1st = nextNovember1st.AddDate(1, 0, 0)
			}
			durationUntilNextNovember := time.Until(nextNovember1st)
			<-time.After(durationUntilNextNovember)
			err := InsertTeamAndPlayerInfo(db, players.GetEndYearOfTheSeason())
			if err != nil {
				log.Println(err)
			}
		}
	}()
}

func ScheduleGOATStatsUpdate(db database.Database) {
	go func() {
		ticker := time.NewTicker(3 * 24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			goatplayers.UpdateActiveGOATStats(db)
		}
	}()
}

func InsertDefaultPolls(db database.Database) {
	pollsInsert := []databasestructs.Poll{
		{ID: 1, Name: "MVP", Description: "Description for MVP", Image: "mvp-trophy.jpg", SelectedStats: "All stats", Season: "2024", UserID: 1},
		{ID: 2, Name: "ROY", Description: "Description for ROY", Image: "roy-trophy.jpeg", SelectedStats: "Rookie", Season: "2024", UserID: 1},
		{ID: 3, Name: "DPOY", Description: "Description for DPOY", Image: "dpoy-trophy.jpeg", SelectedStats: "Defensive", Season: "2024", UserID: 1},
		{ID: 4, Name: "Sixth Man", Description: "Description for 6-man", Image: "6moy-trophy.jpeg", SelectedStats: "Sixth man", Season: "2024", UserID: 1},
		{ID: 5, Name: "GOAT", Description: "Description for GOAT", Image: "6moy-trophy.jpeg", SelectedStats: "GOAT stats", Season: "All", UserID: 1},
	}

	for _, poll := range pollsInsert {
		db.InsertPollsWithId(poll)
	}
}

func SetupSyncSchedules(db database.Database) {
	ScheduleNewSeasonSync(db)
	ScheduleGOATStatsUpdate(db)
	InsertDefaultPolls(db)
}
