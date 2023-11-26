package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sportsvoting/players"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Poll struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	SelectedStats string `json:"selected_stats"`
	Season        string `json:"season"`
	UserID        int64  `json:"user_id,omitempty"`
}

type User struct {
	ID           int64  `json:"id,omitempty"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	ProfilePic   string `json:"profile_pic"`
}

const privateKeyAccessPath = "auth.ed"
const privateKeyRefreshPath = "refresh.ed"
const publicKeyRefreshPath = "refresh.ed.pub"

func issueToken(user, privatekey string) (string, error) {
	keyBytes, err := os.ReadFile(privatekey)
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
	err := db.GetUserByUsername(user).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &u.ProfilePic)
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

	var role string
	err = db.GetUserRolesByID(u.ID).Scan(&role)
	if err != nil {
		fmt.Println("error getting user roles", err)
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

	type Response struct {
		Id          int64  `json:"id"`
		AccessToken string `json:"access_token"`
		Role        string `json:"roles"`
	}

	var res Response
	res.Id = u.ID
	res.AccessToken = accessToken
	res.Role = role
	jsonRes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error marshaling json")
	}

	w.Write(jsonRes)
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	user, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing user/pwd"))
		return
	}

	type Register struct {
		Email string `json:"email"`
	}
	var register Register
	if err := json.NewDecoder(r.Body).Decode(&register); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	err := db.GetUserByUsername(user).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &u.ProfilePic)
	if err == sql.ErrNoRows {
		hash, _ := hashPassword(pwd)
		res, err := db.InsertNewUser(user, register.Email, hash, "")
		if err != nil {
			fmt.Println("error inserting user", err)
		}

		userid, _ := res.LastInsertId()

		_, err = db.InsertUserRoles(userid, UserRoleUser)
		if err != nil {
			fmt.Println("error inserting user roles", err)
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

func createUserAdmin(w http.ResponseWriter, r *http.Request) {
	var reqUser User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	err := db.GetUserByUsername(reqUser.Username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &u.ProfilePic)
	if err == sql.ErrNoRows {
		hash, _ := hashPassword(reqUser.Password)
		res, err := db.InsertNewUser(reqUser.Username, reqUser.Email, hash, "")
		if err != nil {
			fmt.Println("error inserting user", err)
		}

		userid, _ := res.LastInsertId()
		_, err = db.InsertUserRoles(userid, UserRoleUser)
		if err != nil {
			fmt.Println("error inserting user roles", err)
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
	w.Write([]byte("User created successfully"))
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
	err = db.GetUserByRefreshToken(refreshToken).Scan(&u.ID, &u.Username, &u.Email)
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
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	refreshToken := cookie.Value
	var u User
	err = db.GetUserByRefreshToken(refreshToken).Scan(&u.ID, &u.Username, &u.Email)
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

			var currentRoles string
			err = db.GetUserRolesByID(u.ID).Scan(&currentRoles)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			type Response struct {
				ID          int64  `json:"id"`
				Username    string `json:"user"`
				AccessToken string `json:"access_token"`
				Roles       string `json:"roles"`
			}
			resp := Response{ID: u.ID, Username: u.Username, AccessToken: accessToken, Roles: currentRoles}

			json.NewEncoder(w).Encode(resp)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func handleUserList(w http.ResponseWriter, r *http.Request) {
	rows, err := db.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.ProfilePic); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	var profilepic sql.NullString
	err = db.GetUserByID(id).Scan(&u.Username, &u.Email, &profilepic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.ProfilePic = profilepic.String

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getPollByID(id int64) (Poll, error) {
	var poll Poll
	err := db.GetPollByID(id).Scan(&poll.Name, &poll.Description, &poll.Image, &poll.SelectedStats, &poll.Season, &poll.UserID)
	if err != nil {
		return Poll{}, err
	}

	return poll, nil
}

func getPollByUserID(userid int64) ([]Poll, error) {
	rows, err := db.GetPollByUserID(userid)
	if err != nil {
		return nil, err
	}

	var polls []Poll
	for rows.Next() {
		var poll Poll
		err = rows.Scan(&poll.ID, &poll.Name, &poll.Description, &poll.Image, &poll.SelectedStats, &poll.Season)
		if err != nil {
			return nil, err
		}

		polls = append(polls, poll)
	}

	return polls, nil
}

func handleUserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	response := map[string]interface{}{
		"message":       "User deleted successfully",
		"rows_affected": rowsAffected,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func updateUserEmail(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEmail := u.Email
	err := db.GetUserByUsername(u.Username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &sql.NullString{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.UpdateUserEmail(u.Username, updatedEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateUsername(w http.ResponseWriter, r *http.Request) {
	var users struct {
		OldUser  string `json:"olduser"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.UpdateUserUsername(users.OldUser, users.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updatePassword(w http.ResponseWriter, r *http.Request) {
	var newPasswords struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
		Username    string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPasswords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var u User
	err := db.GetUserByUsername(newPasswords.Username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &sql.NullString{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !checkPasswordHash(newPasswords.OldPassword, u.Password) {
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	u.Password, err = hashPassword(newPasswords.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.UpdateUserPassword(u.Username, u.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateAdmin(w http.ResponseWriter, r *http.Request) {
	var id int64
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentRoles string
	err := db.GetUserRolesByID(id).Scan(&currentRoles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newRoles string
	// just reverse whether it is already an admin or not
	if strings.Contains(currentRoles, "admin") {
		newRoles = UserRoleUser
	} else {
		newRoles = fmt.Sprintf("%s,%s", UserRoleUser, UserRoleAdmin)
	}

	_, err = db.UpdateUserRoles(newRoles, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newRoles)
}

func updatePoll(w http.ResponseWriter, r *http.Request) {
	var poll Poll
	if err := json.NewDecoder(r.Body).Decode(&poll); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var p Poll
	err := db.GetPollByID(poll.ID).Scan(&p.Name, &p.Description, &p.Image, &p.SelectedStats, &p.Season, &p.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.UpdatePollByID(poll.Name, poll.Description, poll.SelectedStats, poll.Season, poll.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// if the season or stats changed for the poll, rest the votes
	if p.Season != poll.Season || p.SelectedStats != poll.SelectedStats {
		db.ResetPollVotes(poll.ID)
	}

	w.WriteHeader(http.StatusOK)
}

func createQuiz(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size
	var poll Poll
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
	insertRes, err := db.InsertPolls(poll.Name, poll.Description, poll.Image, poll.SelectedStats, poll.Season, poll.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	insertID, _ := insertRes.LastInsertId()
	response := fmt.Sprintf("Created quiz with ID: %d", insertID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}

func deletePollByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["pollid"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.DeletePollByID(id)
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

func resetPollVotes(w http.ResponseWriter, r *http.Request) {
	var id int64
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.ResetPollVotes(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func uploadProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size

	file, _, err := r.FormFile("profileImage")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	uploadDir := "public/"
	fileName := r.MultipartForm.File["profileImage"][0].Filename
	fmt.Println(fileName)
	username := r.FormValue("username")
	fmt.Println(username)

	var u User
	var profilepic sql.NullString
	err = db.GetUserByUsername(username).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.RefreshToken, &profilepic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u.ProfilePic = profilepic.String
	if u.ProfilePic != "" {
		err := os.Remove(uploadDir + u.ProfilePic)
		if err != nil {
			http.Error(w, "Error deleting old profile picture", http.StatusInternalServerError)
			return
		}
	}

	// Create a new file on the server and write the uploaded file to it
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

	_, err = db.UpdateUserProfilePic(username, fileName)
	if err != nil {
		http.Error(w, "Unable to update profile pic", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":  true,
		"fileName": fileName,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updatePollImage(w http.ResponseWriter, r *http.Request) {
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

	var p Poll
	var pollimage sql.NullString
	err = db.GetPollByID(pollIdInt).Scan(&p.Name, &p.Description, &pollimage, &p.SelectedStats, &p.Season, &p.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p.Image = pollimage.String
	if p.Image != "" {
		err := os.Remove(uploadDir + p.Image)
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

	_, err = db.UpdatePollImage(pollIdInt, fileName)
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

	for rows.Next() {
		var poll Poll
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

type MyVotesResponse struct {
	PlayerID   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	PollName   string `json:"poll_name"`
	PollImage  string `json:"poll_image"`
}

func getUserVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["userid"])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	rows, err := db.GetVotesOfUser(ctx, int64(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var allVotes []MyVotesResponse
	for rows.Next() {
		var votes MyVotesResponse
		err = rows.Scan(&votes.PlayerID, &votes.PlayerName, &votes.PollName, &votes.PollImage)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		allVotes = append(allVotes, votes)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allVotes)
}

func getSeasons(w http.ResponseWriter, r *http.Request) {
	var seasons []string

	rows, err := db.SelectSeasonsAvailable()
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

func GetQuiz(w http.ResponseWriter, r *http.Request) {
	pollId := mux.Vars(r)["pollid"]
	id, err := strconv.ParseInt(pollId, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	poll, err := getPollByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var players []players.Player
	switch poll.SelectedStats {
	case "Defensive":
		players, err = getDefensiveStats(ctx, poll.Season)
	case "Sixth man":
		players, err = getSixmanStats(ctx, poll.Season)
	case "Rookie":
		players, err = getRookieStats(ctx, poll.Season)
	case "All stats":
		players, err = getAllStats(ctx, poll.Season)
	default:
		players, err = getAllStats(ctx, poll.Season)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(players)
}

func GetQuizById(w http.ResponseWriter, r *http.Request) {
	pollId := mux.Vars(r)["pollid"]
	id, err := strconv.ParseInt(pollId, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	poll, err := getPollByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
	json.NewEncoder(w).Encode(poll)
}
func getUserPolls(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userid"]
	id, err := strconv.ParseInt(userId, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	polls, err := getPollByUserID(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(polls)
}

func getRookieStats(ctx context.Context, season string) ([]players.Player, error) {
	rows, err := db.GetROYStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []players.Player

	for rows.Next() {
		var p players.Player
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

func getAllStats(ctx context.Context, season string) ([]players.Player, error) {
	rows, err := db.GetPlayerStatsForQuiz(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []players.Player

	for rows.Next() {
		var p players.Player
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

func getSixmanStats(ctx context.Context, season string) ([]players.Player, error) {
	rows, err := db.GetSixManStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var playerList []players.Player

	for rows.Next() {
		var p players.Player
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

func getDefensiveStats(ctx context.Context, season string) ([]players.Player, error) {
	rows, err := db.GetDPOYStats(ctx, season)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playerList []players.Player
	for rows.Next() {
		var p players.Player
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
	UserID   int64  `json:"userid"`
}

func insertPlayerVotes(w http.ResponseWriter, r *http.Request) {
	var payload VotePayload

	// Decode the JSON request body into the VotePayload struct
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.InsertPlayerVotes(payload.PollID, payload.UserID, payload.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Voted for player %s in poll %d", payload.PlayerID, payload.PollID)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(response))
}
