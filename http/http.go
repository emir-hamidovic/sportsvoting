package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sportsvoting/database"
	"sportsvoting/polls"
	"sportsvoting/syncer"
	"sportsvoting/users"
	"sportsvoting/votes"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func SetupHandlers(db database.Database) *mux.Router {
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

	return r
}

func StartServer(db database.Database) {
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
		Handler:      c.Handler(SetupHandlers(db)),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	go syncer.RunUpdate(db, ctx)

	fmt.Println("Starting server")
	err := srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

	cancel()
}
