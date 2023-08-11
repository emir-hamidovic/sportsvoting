package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sportsvoting/players"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Poll struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Endpoint    string `json:"endpoint"`
}

type User struct {
	ID           int64  `json:"id,omitempty"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	IsAdmin      string `json:"is_admin"`
}

const privateKeyAccessPath = "auth.ed"
const privateKeyRefreshPath = "refresh.ed"
const publicKeyRefreshPath = "refresh.ed.pub"

func issueToken(user, privatekey string) (string, error) {
	keyBytes, err := ioutil.ReadFile(privatekey)
	if err != nil {
		panic(fmt.Errorf("unable to read private key file: %w", err))
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return "", fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	now := time.Now()
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		"aud":      "api",
		"nbf":      now.Unix(),
		"iat":      now.Unix(),
		"exp":      now.Add(10 * time.Minute).Unix(),
		"iss":      "http://localhost:8080",
		"username": user,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return tokenString, nil
}

func validateToken(tokenString, publickey string) (*jwt.Token, error) {
	keyBytes, err := ioutil.ReadFile(publickey)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %w", err)
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// return the single public key we trust
			return key, nil
		})
	if err != nil {
		return nil, fmt.Errorf("unable to parse token string: %w", err)
	}

	return token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	user, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing basic auth"))
		return
	}

	var u User
	var match bool
	err := db.GetUserByUsername(user).Scan(&u.Username, &u.Password, &u.RefreshToken, &u.IsAdmin)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
		return
	} else {
		match = checkPasswordHash(pwd, u.Password)
	}

	if !match {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
		return
	}

	accessToken, err := issueToken(user, privateKeyAccessPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to issue token:" + err.Error()))
		return
	}

	refreshToken, err := issueToken(user, privateKeyRefreshPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to issue token:" + err.Error()))
		return
	}

	_, err = db.UpdateUserRefreshToken(user, refreshToken)
	if err != nil {
		fmt.Println("Error updating refresh token for user ", user)
	}

	cookie := http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   86400,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)

	w.Write([]byte(accessToken))
	//json.NewEncoder(w).Encode(accessToken)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	user, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing user/pwd"))
		return
	}

	var u User
	err := db.GetUserByUsername(user).Scan(&u.Username, &u.Password, &u.RefreshToken, &u.IsAdmin)
	if err == sql.ErrNoRows {
		hash, _ := hashPassword(pwd)
		_, err := db.InsertNewUser(user, hash, "", false)
		if err != nil {
			fmt.Println("error inserting user")
		}
	} else if err != nil {
		fmt.Println("error inserting user", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if u.Username != "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("user already exists"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register successful"))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusNoContent)
		default:
			fmt.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	refreshToken := cookie.Value
	var u User
	err = db.GetUserByRefreshToken(refreshToken).Scan(&u.Username)
	if err != nil {
		fmt.Println("can't find user by refresh token", err)
	} else if u.Username != "" {
		_, err := db.UpdateUserRefreshToken(u.Username, "")
		if err != nil {
			fmt.Println("error deleting refresh token from db ", err)
		}
	}

	c := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, c)
	w.WriteHeader(http.StatusNoContent)
}

func handleRefresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusUnauthorized)
		default:
			fmt.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	refreshToken := cookie.Value
	var u User
	err = db.GetUserByRefreshToken(refreshToken).Scan(&u.Username)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		fmt.Println("can't find user by refresh token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if u.Username != "" {
		token, err := validateToken(refreshToken, publicKeyRefreshPath)
		if err != nil {
			fmt.Println(refreshToken)
			fmt.Println(err)
			http.Error(w, "token not valid", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["username"] == u.Username {
			accessToken, err := issueToken(u.Username, privateKeyAccessPath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("unable to issue token:" + err.Error()))
				return
			}

			json.NewEncoder(w).Encode(accessToken)
		} else {
			fmt.Println("Invalid JWT Token")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
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
