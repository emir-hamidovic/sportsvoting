package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sportsvoting/players"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Poll struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Endpoint    string `json:"endpoint"`
}

func getPolls(w http.ResponseWriter, r *http.Request) {
	var polls []Poll

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetPolls(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate over the rows and populate the polls slice
	for rows.Next() {
		var poll Poll
		err := rows.Scan(&poll.ID, &poll.Name, &poll.Description, &poll.Image, &poll.Endpoint)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		polls = append(polls, poll)
	}

	// Convert the polls slice to JSON
	pollsJSON, err := json.Marshal(polls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Write the JSON response
	w.Write(pollsJSON)
}

func dpoyAward(w http.ResponseWriter, r *http.Request) {
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
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Rebounds, &p.Steals, &p.Blocks, &p.Position, &p.DefWS, &p.DefBPM, &p.DefRtg)
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

func mipAward(w http.ResponseWriter, r *http.Request) {
	mvpAward(w, r)
}

func royAward(w http.ResponseWriter, r *http.Request) {
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
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.WS, &p.BPM, &p.OffRtg, &p.DefRtg)
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

func sixManAward(w http.ResponseWriter, r *http.Request) {
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
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
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

func mvpAward(w http.ResponseWriter, r *http.Request) {
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
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
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

type Votes struct {
	Name     string `json:"name"`
	Value    int64  `json:"value"`
	Pollname string `json:"pollname"`
}

func playerVotes(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024)
	pollid := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(pollid, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetPlayerPollVotes(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var playerList []Votes
	for rows.Next() {
		var v Votes
		err := rows.Scan(&v.Name, &v.Value, &v.Pollname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		playerList = append(playerList, v)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(playerList)
}

func teamVotes(w http.ResponseWriter, r *http.Request) {}

type VotePayload struct {
	PlayerID string `json:"playerid"`
	PollID   int64  `json:"pollid"`
}

func insertPlayerVotes(w http.ResponseWriter, r *http.Request) {
	var payload VotePayload

	// Decode the JSON request body into the VotePayload struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.InsertPlayerVotes(payload.PollID, payload.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Voted for player %s in poll %d", payload.PlayerID, payload.PollID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
}
