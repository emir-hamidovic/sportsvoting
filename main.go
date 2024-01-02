package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/goatplayers"
	"sportsvoting/players"
	"sportsvoting/polls"
	"sportsvoting/teams"
	"sportsvoting/users"
	"sportsvoting/votes"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

func isSyncNeeded(db database.Database) (bool, error) {
	syncTime, err := db.GetLastSyncTime("Regular")
	if err != nil {
		log.Println(err)
		return true, err
	}

	threshold := 10 * 24 * time.Hour
	timeDiff := time.Since(syncTime)

	return timeDiff > threshold, nil
}

func isGOATSyncNeeded(db database.Database) (bool, error) {
	syncTime, err := db.GetLastSyncTime("GOAT")
	if err != nil {
		log.Println(err)
		return true, err
	}

	threshold := 10 * 24 * time.Hour
	timeDiff := time.Since(syncTime)

	return timeDiff > threshold, nil
}

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

		log.Printf("Couldn't connect to databse due to: %v\n", err)

		if i < retries {
			log.Printf("Retrying in %v ... \n", delayTime)
			time.Sleep(delayTime)
		}
	}

	if db == nil {
		log.Fatalf("Couldn't connect to database after %d retries, exiting!\n", retries)
	}

	defer db.CloseConnection()

	err = db.CreateAdminUser()
	if err != nil {
		log.Fatal(err)
	}

	isSyncNeeded, errSync := isSyncNeeded(db)
	if isSyncNeeded {
		err = InsertTeamAndPlayerInfo(db, players.GetEndYearOfTheSeason())
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

	isSyncNeeded, errSync = isGOATSyncNeeded(db)
	if isSyncNeeded {
		go func() {
			playerIDs := goatplayers.GetGoatPlayersList()
			goatplayers.InsertGoatPlayerStats(playerIDs, db)

			_, err = db.InsertSeasonEntered("All")
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

	// Schedule the function to run on the next 1st of November
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

	go func() {
		ticker := time.NewTicker(3 * 24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			goatplayers.UpdateActiveGOATStats(db)
		}
	}()

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

	usersHandler := users.UsersHandler{DB: db}
	votesHandler := votes.VotesHandler{DB: db}
	pollsHandler := polls.PollsHandler{DB: db}
	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", usersHandler.HandleRegister).Methods("POST")
	api.HandleFunc("/login", usersHandler.HandleLogin).Methods("POST")
	api.HandleFunc("/logout", usersHandler.HandleLogout)
	api.HandleFunc("/refresh", usersHandler.HandleRefresh)

	api.HandleFunc("/users/get/{id:[0-9]+}", usersHandler.HandleGetUserByID)
	api.HandleFunc("/users/get", usersHandler.HandleUserList)
	api.HandleFunc("/users/delete/{id:[0-9]+}", usersHandler.HandleUserDelete).Methods("DELETE")
	api.HandleFunc("/users/email/update", usersHandler.UpdateUserEmail).Methods("POST")
	api.HandleFunc("/users/username/update", usersHandler.UpdateUsername).Methods("POST")
	api.HandleFunc("/users/password/update", usersHandler.UpdatePassword).Methods("POST")
	api.HandleFunc("/users/admin/update", usersHandler.UpdateAdmin).Methods("POST")
	api.HandleFunc("/users/image/update", usersHandler.UploadProfilePicHandler).Methods("POST")
	api.HandleFunc("/users/admin/create", usersHandler.CreateUserAdmin).Methods("POST")

	api.HandleFunc("/polls/players/get/{pollid:[0-9]+}", pollsHandler.GetPlayerStatsForPoll)
	api.HandleFunc("/polls/get/{pollid:[0-9]+}", pollsHandler.GetPollById)
	api.HandleFunc("/polls/get", pollsHandler.GetPolls)
	api.HandleFunc("/polls/create", pollsHandler.CreatePoll).Methods("POST")
	api.HandleFunc("/polls/update", pollsHandler.UpdatePoll).Methods("POST")
	api.HandleFunc("/polls/delete/{pollid:[0-9]+}", pollsHandler.DeletePollByID).Methods("DELETE")
	api.HandleFunc("/polls/image/update", pollsHandler.UpdatePollImage).Methods("POST")
	api.HandleFunc("/polls/votes/reset", pollsHandler.ResetPollVotes).Methods("POST")
	api.HandleFunc("/polls/users/get/{userid}", pollsHandler.GetUserPolls)

	api.HandleFunc("/votes/users/get/{userid}", votesHandler.GetUserVotes)
	api.HandleFunc("/votes/players/{id:[0-9]+}", votesHandler.PlayerVotes).Methods("GET")
	api.HandleFunc("/votes/players", votesHandler.InsertPlayerVotes).Methods("POST")
	api.HandleFunc("/votes/teams/{id:[0-9]+}", votesHandler.TeamVotes)

	api.HandleFunc("/seasons/get", pollsHandler.GetSeasons)

	isDev := true
	if isdevEnv, exists := os.LookupEnv("IS_DEVELOPMENT"); exists {
		isDev, _ = strconv.ParseBool(isdevEnv)
	}

	var allowedOrigins []string
	if isDev {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost/"}
	} else {
		allowedOrigins = []string{"http://frontend:80"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedHeaders:   []string{"Content-type", "Authorization"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      c.Handler(r),
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
