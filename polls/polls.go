package polls

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type PollsHandler struct {
	DB database.Database
}

func parseID(r *http.Request, key string) (int64, error) {
	pollId := mux.Vars(r)[key]
	return strconv.ParseInt(pollId, 10, 64)
}

func (p PollsHandler) GetPlayerStatsForPoll(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "pollid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	poll, err := p.getPollByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var players []databasestructs.PlayerInfo
	switch poll.SelectedStats {
	case "Defensive":
		players, err = p.getDefensiveStats(ctx, poll.Season)
	case "Sixth man":
		players, err = p.getSixmanStats(ctx, poll.Season)
	case "Rookie":
		players, err = p.getRookieStats(ctx, poll.Season)
	case "All stats":
		players, err = p.getAllStats(ctx, poll.Season)
	case "GOAT stats":
		players, err = p.getGOATStats(poll.Season)
	default:
		players, err = p.getAllStats(ctx, poll.Season)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(players)
}

func (p PollsHandler) GetPollById(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "pollid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := p.getPollByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(poll)
}

func (p PollsHandler) getPollByID(id int64) (databasestructs.Poll, error) {
	var poll databasestructs.Poll
	err := p.DB.GetPollByID(id).Scan(&poll.Name, &poll.Description, &poll.Image, &poll.SelectedStats, &poll.Season, &poll.UserID)
	if err != nil {
		return databasestructs.Poll{}, err
	}

	return poll, nil
}

func (p PollsHandler) GetPolls(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := p.DB.GetPolls(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var polls []databasestructs.Poll
	for rows.Next() {
		var poll databasestructs.Poll
		err := rows.Scan(&poll.ID, &poll.Name, &poll.Description, &poll.Image, &poll.SelectedStats, &poll.Season, &poll.UserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		polls = append(polls, poll)
	}

	pollsJSON, err := json.Marshal(polls)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(pollsJSON)
}

func (p PollsHandler) CreatePoll(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
	var poll databasestructs.Poll
	poll.Name = r.FormValue("name")
	poll.Description = r.FormValue("description")
	poll.Season = r.FormValue("season")
	poll.SelectedStats = r.FormValue("selectedStats")
	userID, err := strconv.ParseInt(r.FormValue("userid"), 10, 64)
	if err != nil {
		http.Error(w, "Unable to parse user id", http.StatusBadRequest)
		return
	}

	poll.UserID = userID
	image, _, err := r.FormFile("photo")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer image.Close()

	isDev := true
	if isdevEnv, exists := os.LookupEnv("IS_DEVELOPMENT"); exists {
		isDev, _ = strconv.ParseBool(isdevEnv)
	}

	uploadDir := "public/"
	if !isDev {
		uploadDir = "/app/shared/"
	}
	poll.Image = r.MultipartForm.File["photo"][0].Filename
	newFile, err := os.Create(uploadDir + poll.Image)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, image)
	if err != nil {
		http.Error(w, "Unable to copy file", http.StatusInternalServerError)
		return
	}
	insertRes, err := p.DB.InsertPolls(poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertID, _ := insertRes.LastInsertId()
	response := fmt.Sprintf("Created poll with ID: %d", insertID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func (p PollsHandler) UpdatePoll(w http.ResponseWriter, r *http.Request) {
	var poll databasestructs.Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pollDB databasestructs.Poll
	err := p.DB.GetPollByID(poll.ID).Scan(&pollDB.Name, &pollDB.Description, &pollDB.Image, &pollDB.SelectedStats, &pollDB.Season, &pollDB.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = p.DB.UpdatePollByID(poll)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if the season or stats changed for the poll, rest the votes
	if pollDB.Season != poll.Season || pollDB.SelectedStats != poll.SelectedStats {
		p.DB.ResetPollVotes(poll.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func (p PollsHandler) DeletePollByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "pollid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := p.DB.DeletePollByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	response := map[string]interface{}{
		"message":       "Poll deleted successfully",
		"rows_affected": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (p PollsHandler) UpdatePollImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size

	file, _, err := r.FormFile("pollImage")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	uploadDir := "frontend/public/"
	fileName := r.MultipartForm.File["pollImage"][0].Filename
	pollId := r.FormValue("pollId")
	pollIdInt, _ := strconv.ParseInt(pollId, 10, 64)

	var poll databasestructs.Poll
	var pollimage sql.NullString
	err = p.DB.GetPollByID(pollIdInt).Scan(&poll.Name, &poll.Description, &pollimage, &poll.SelectedStats, &poll.Season, &poll.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	poll.Image = pollimage.String
	if poll.Image != "" {
		err := os.Remove(uploadDir + poll.Image)
		if err != nil {
			http.Error(w, "Error deleting old poll image", http.StatusInternalServerError)
			return
		}
	}

	newFile, err := os.Create(uploadDir + fileName)
	if err != nil {
		http.Error(w, "Unable to create file", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = io.Copy(newFile, file)
	if err != nil {
		http.Error(w, "Unable to copy file", http.StatusInternalServerError)
		return
	}

	img := databasestructs.Image{ID: pollIdInt, ImageURL: fileName}
	_, err = p.DB.UpdatePollImage(img)
	if err != nil {
		http.Error(w, "Unable to update poll image", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"fileName": fileName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (p PollsHandler) ResetPollVotes(w http.ResponseWriter, r *http.Request) {
	var id int64
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := p.DB.ResetPollVotes(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (p PollsHandler) GetUserPolls(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r, "userid")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	polls, err := p.getPollByUserID(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
}

func (p PollsHandler) getPollByUserID(userid int64) ([]databasestructs.Poll, error) {
	rows, err := p.DB.GetPollByUserID(userid)
	if err != nil {
		return nil, err
	}

	var polls []databasestructs.Poll
	for rows.Next() {
		var poll databasestructs.Poll
		err = rows.Scan(&poll.ID, &poll.Name, &poll.Description, &poll.Image, &poll.SelectedStats, &poll.Season)
		if err != nil {
			return nil, err
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func (p PollsHandler) GetSeasons(w http.ResponseWriter, r *http.Request) {
	var seasons []string

	rows, err := p.DB.SelectSeasonsAvailable()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var season string
		err := rows.Scan(&season)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		seasons = append(seasons, season)
	}

	seasonsJson, err := json.Marshal(seasons)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(seasonsJson)
}

func (p PollsHandler) getRookieStats(ctx context.Context, season string) ([]databasestructs.PlayerInfo, error) {
	rows, err := p.DB.GetROYStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []databasestructs.PlayerInfo

	for rows.Next() {
		var p databasestructs.PlayerInfo
		err = rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.WS, &p.BPM, &p.OffRtg, &p.DefRtg)
		if err != nil {
			return nil, err
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playerList, nil
}

func (p PollsHandler) getAllStats(ctx context.Context, season string) ([]databasestructs.PlayerInfo, error) {
	rows, err := p.DB.GetPlayerStatsForPoll(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []databasestructs.PlayerInfo

	for rows.Next() {
		var p databasestructs.PlayerInfo
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
		if err != nil {
			return nil, err
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playerList, nil
}

func (p PollsHandler) getGOATStats(season string) ([]databasestructs.PlayerInfo, error) {
	rows, err := p.DB.GetGOATStats(season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []databasestructs.PlayerInfo

	for rows.Next() {
		var p databasestructs.PlayerInfo
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
		if err != nil {
			return nil, err
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playerList, nil
}

func (p PollsHandler) getSixmanStats(ctx context.Context, season string) ([]databasestructs.PlayerInfo, error) {
	rows, err := p.DB.GetSixManStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []databasestructs.PlayerInfo

	for rows.Next() {
		var p databasestructs.PlayerInfo
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Points, &p.Rebounds, &p.Assists, &p.Steals, &p.Blocks, &p.FGPercentage, &p.ThreeFGPercentage, &p.FTPercentage, &p.Turnovers, &p.Position, &p.PER, &p.OffWS, &p.DefWS, &p.WS, &p.OffBPM, &p.DefBPM, &p.BPM, &p.VORP, &p.OffRtg, &p.DefRtg)
		if err != nil {
			return nil, err
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playerList, nil
}

func (p PollsHandler) getDefensiveStats(ctx context.Context, season string) ([]databasestructs.PlayerInfo, error) {
	rows, err := p.DB.GetDPOYStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playerList []databasestructs.PlayerInfo
	for rows.Next() {
		var p databasestructs.PlayerInfo
		err := rows.Scan(&p.ID, &p.Name, &p.Games, &p.Minutes, &p.Rebounds, &p.Steals, &p.Blocks, &p.Position, &p.DefWS, &p.DefBPM, &p.DefRtg)
		if err != nil {
			return nil, err
		}
		playerList = append(playerList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return playerList, nil
}
