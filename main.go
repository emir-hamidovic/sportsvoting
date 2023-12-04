package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sportsvoting/database"
	"sportsvoting/players"
	"sportsvoting/teams"
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

const (
	UserRoleAdmin string = "admin"
	UserRoleUser  string = "user"
)

var db database.Database

func main() {
	db_addr := os.Getenv("DBADDRESS")
	if db_addr == "" {
		db_addr = "localhost:3306"
	}

	var err error
	db, err = database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: db_addr})
	if err != nil {
		log.Fatal(err)
	}

	err = db.CreateAdminUser()
	if err != nil {
		log.Fatal(err)
	}

	/*
		err = InsertTeamAndPlayerInfo(db, players.GetEndYearOfTheSeason())
		if err != nil {
			log.Fatal(err)
		}
	*/
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
				log.Fatal(err)
			}
		}
	}()

	polls := []Poll{
		{1, "MVP", "Description for MVP", "mvp-trophy.jpg", "All stats", "2023", 1},
		{2, "ROY", "Description for ROY", "roy-trophy.jpeg", "Rookie", "2023", 1},
		{3, "DPOY", "Description for DPOY", "dpoy-trophy.jpeg", "Defensive", "2023", 1},
		{4, "Sixth Man", "Description for 6-man", "6moy-trophy.jpeg", "Sixth man", "2023", 1},
	}

	for _, val := range polls {
		db.InsertPollsWithId(val.ID, val.Name, val.Description, val.Image, val.SelectedStats, val.Season, val.UserID)
	}

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()

	api.HandleFunc("/register", handleRegister).Methods("POST")
	api.HandleFunc("/login", handleLogin).Methods("POST")
	api.HandleFunc("/logout", handleLogout)
	api.HandleFunc("/refresh", handleRefresh)

	api.HandleFunc("/users/get/{id:[0-9]+}", handleGetUserByID)
	api.HandleFunc("/users/get", handleUserList)
	api.HandleFunc("/users/delete/{id:[0-9]+}", handleUserDelete).Methods("DELETE")
	api.HandleFunc("/users/email/update", updateUserEmail).Methods("POST")
	api.HandleFunc("/users/username/update", updateUsername).Methods("POST")
	api.HandleFunc("/users/password/update", updatePassword).Methods("POST")
	api.HandleFunc("/users/admin/update", updateAdmin).Methods("POST")
	api.HandleFunc("/users/image/update", uploadProfilePicHandler).Methods("POST")
	api.HandleFunc("/users/admin/create", createUserAdmin).Methods("POST")

	api.HandleFunc("/polls/players/get/{pollid:[0-9]+}", GetQuiz)
	api.HandleFunc("/polls/get/{pollid:[0-9]+}", GetQuizById)
	api.HandleFunc("/polls/get", getPolls)
	api.HandleFunc("/polls/delete/{pollid:[0-9]+}", deletePollByID).Methods("DELETE")
	api.HandleFunc("/polls/create", createQuiz).Methods("POST")
	api.HandleFunc("/polls/image/update", updatePollImage).Methods("POST")
	api.HandleFunc("/polls/votes/reset", resetPollVotes).Methods("POST")
	api.HandleFunc("/polls/update", updatePoll).Methods("POST")
	api.HandleFunc("/polls/users/get/{userid}", getUserPolls)

	api.HandleFunc("/votes/user/get/{userid}", getUserVotes)
	api.HandleFunc("/votes/players/{id:[0-9]+}", playerVotes).Methods("GET")
	api.HandleFunc("/votes/players", insertPlayerVotes).Methods("POST")
	api.HandleFunc("/votes/teams/{id:[0-9]+}", teamVotes)

	api.HandleFunc("/seasons/get", getSeasons)

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
