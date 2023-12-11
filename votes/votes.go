package votes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sportsvoting/database"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Votes struct {
	Name     string `json:"name"`
	Value    int64  `json:"value"`
	Pollname string `json:"pollname"`
}

type MyVotesResponse struct {
	PollID     string `json:"poll_id"`
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	PollName   string `json:"poll_name"`
	PollImage  string `json:"poll_image"`
}

type VotePayload struct {
	PlayerID string `json:"playerid"`
	PollID   int64  `json:"pollid"`
	UserID   int64  `json:"userid"`
}

type VotesHandler struct {
	DB database.Database
}

func (v VotesHandler) GetUserVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userid"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := v.DB.GetVotesOfUser(ctx, int64(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allVotes []MyVotesResponse
	for rows.Next() {
		var votes MyVotesResponse
		err = rows.Scan(&votes.PollID, &votes.PlayerID, &votes.PlayerName, &votes.PollName, &votes.PollImage)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		allVotes = append(allVotes, votes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allVotes)
}

func (v VotesHandler) PlayerVotes(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	pollid := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(pollid, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := v.DB.GetPlayerPollVotes(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []Votes
	for rows.Next() {
		var votes Votes
		err := rows.Scan(&votes.Name, &votes.Value, &votes.Pollname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, votes)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
}

func (v VotesHandler) TeamVotes(w http.ResponseWriter, r *http.Request) {}

func (v VotesHandler) InsertPlayerVotes(w http.ResponseWriter, r *http.Request) {
	var payload VotePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = v.DB.InsertPlayerVotes(payload.PollID, payload.UserID, payload.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Voted for player %s in poll %d", payload.PlayerID, payload.PollID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}
