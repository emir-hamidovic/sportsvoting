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
	var err error
	db, err = database.NewDB(database.Config{DbType: "mysql", DbName: "nba", Addr: "localhost:3306"})
	if err != nil {
		log.Fatal(err)
	}

	err = InsertTeamAndPlayerInfo(db, players.GetEndYearOfTheSeason())
	if err != nil {
		log.Fatal(err)
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
				log.Fatal(err)
			}
		}
	}()

	polls := []Poll{
		{1, "MVP", "Description for MVP", "mvp-trophy.jpg", "All stats", "2023"},
		{2, "ROY", "Description for ROY", "roy-trophy.jpeg", "Rookie", "2023"},
		{3, "DPOY", "Description for DPOY", "dpoy-trophy.jpeg", "Defensive", "2023"},
		{4, "Sixth Man", "Description for 6-man", "6moy-trophy.jpeg", "Sixth man", "2023"},
	}

	for _, val := range polls {
		db.InsertPollsWithId(val.ID, val.Name, val.Description, val.Image, val.SelectedStats, val.Season)
	}

	r := mux.NewRouter()
	r.HandleFunc("/getpolls", getPolls)
	r.HandleFunc("/getseasons", getSeasons)
	r.HandleFunc("/quiz/{pollid:[0-9]+}", GetQuiz)
	r.HandleFunc("/teamvotes/{id:[0-9]+}", teamVotes)
	r.HandleFunc("/playervotes/{id:[0-9]+}", playerVotes).Methods("GET")
	r.HandleFunc("/playervotes/", insertPlayerVotes).Methods("POST")
	r.HandleFunc("/login", handleLogin).Methods("POST")
	r.HandleFunc("/register", handleRegister).Methods("POST")
	r.HandleFunc("/admin/createuser", createUserAdmin).Methods("POST")
	r.HandleFunc("/logout", handleLogout)
	r.HandleFunc("/refresh", handleRefresh)
	r.HandleFunc("/users/get", handleUserList)
	r.HandleFunc("/users/delete/{id:[0-9]+}", handleUserDelete).Methods("DELETE")
	r.HandleFunc("/api/get-user/{id:[0-9]+}", handleGetUserByID)
	r.HandleFunc("/api/update-email", updateUserEmail).Methods("POST")
	r.HandleFunc("/api/update-username", updateUsername).Methods("POST")
	r.HandleFunc("/api/update-password", updatePassword).Methods("POST")
	r.HandleFunc("/api/update-admin", updateAdmin).Methods("POST")
	r.HandleFunc("/api/upload-profile-pic", uploadProfilePicHandler).Methods("POST")
	r.HandleFunc("/api/create-quiz", createQuiz).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost/"},
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
