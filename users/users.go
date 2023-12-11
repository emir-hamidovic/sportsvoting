package users

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

const (
	UserRoleAdmin string = "admin"
	UserRoleUser  string = "user"
)

type UsersHandler struct {
	DB database.Database
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
		"exp":      now.Add(1 * time.Hour).Unix(),
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

func (u UsersHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
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

	err := u.createNewUser(username, password, register.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Register successful"))
}

func (u UsersHandler) createNewUser(username, password, email string) error {
	var user databasestructs.User
	err := u.DB.GetUserByUsername(username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.ProfilePic)
	if err == sql.ErrNoRows {
		hash, _ := hashPassword(password)
		userDb := databasestructs.User{Username: username, Email: email, Password: hash}
		res, err := u.DB.InsertNewUser(userDb)
		if err != nil {
			return err
		}

		userid, _ := res.LastInsertId()

		roles := databasestructs.Role{UserID: userid, Role: UserRoleUser}
		_, err = u.DB.InsertUserRoles(roles)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if user.Username != "" {
		return errors.New("user already exists")
	}

	return nil
}

func (u UsersHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	user, pwd, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing basic auth"))
		return
	}

	var userDB databasestructs.User
	var match bool
	err := u.DB.GetUserByUsername(user).Scan(&userDB.ID, &userDB.Username, &userDB.Email, &userDB.Password, &userDB.RefreshToken, &userDB.ProfilePic)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials"))
		return
	} else {
		match = checkPasswordHash(pwd, userDB.Password)
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
	err = u.DB.GetUserRolesByID(userDB.ID).Scan(&role)
	if err != nil {
		fmt.Println("error getting user roles", err)
	}

	_, err = u.DB.UpdateUserRefreshToken(user, refreshToken)
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

	jsonResp, err := u.returnTokenAndRoleOfUser(userDB.ID, accessToken, role, "")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println("Error marshaling json")
	}

	w.Write(jsonResp)
}

func (u UsersHandler) returnTokenAndRoleOfUser(id int64, accessToken, roles, username string) ([]byte, error) {
	type Response struct {
		ID          int64  `json:"id"`
		AccessToken string `json:"access_token"`
		Username    string `json:"user"`
		Roles       string `json:"roles"`
	}
	resp := Response{ID: id, Username: username, AccessToken: accessToken, Roles: roles}

	return json.Marshal(resp)
}

func (u UsersHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
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
	var user databasestructs.User
	err = u.DB.GetUserByRefreshToken(refreshToken).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		fmt.Println("can't find user by refresh token", err)
	} else if user.Username != "" {
		_, err := u.DB.UpdateUserRefreshToken(user.Username, "")
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

func (u UsersHandler) HandleRefresh(w http.ResponseWriter, r *http.Request) {
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
	var user databasestructs.User
	err = u.DB.GetUserByRefreshToken(refreshToken).Scan(&user.ID, &user.Username, &user.Email)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if err != nil {
		fmt.Println("can't find user by refresh token", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else if user.Username != "" {
		token, err := validateToken(refreshToken, publicKeyRefreshPath)
		if err != nil {
			http.Error(w, "token not valid", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["username"] == user.Username {
			accessToken, err := issueToken(user.Username, privateKeyAccessPath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("unable to issue token:" + err.Error()))
				return
			}

			var currentRoles string
			err = u.DB.GetUserRolesByID(user.ID).Scan(&currentRoles)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resp, err := u.returnTokenAndRoleOfUser(user.ID, accessToken, currentRoles, user.Username)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(resp)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func (u UsersHandler) HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user databasestructs.User
	var profilepic sql.NullString
	err = u.DB.GetUserByID(id).Scan(&user.Username, &user.Email, &profilepic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ProfilePic = profilepic.String

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (u UsersHandler) HandleUserList(w http.ResponseWriter, r *http.Request) {
	rows, err := u.DB.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []databasestructs.User

	for rows.Next() {
		var user databasestructs.User
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

func (u UsersHandler) HandleUserDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := u.DB.DeleteUser(id)
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

func (u UsersHandler) UpdateUserEmail(w http.ResponseWriter, r *http.Request) {
	var user databasestructs.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedEmail := user.Email
	err := u.DB.GetUserByUsername(user.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &sql.NullString{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = u.DB.UpdateUserEmail(user.Username, updatedEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u UsersHandler) UpdateUsername(w http.ResponseWriter, r *http.Request) {
	var users struct {
		OldUser  string `json:"olduser"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&users); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := u.DB.UpdateUserUsername(users.OldUser, users.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (u UsersHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	var newPasswords struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
		Username    string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&newPasswords); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user databasestructs.User
	err := u.DB.GetUserByUsername(newPasswords.Username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &sql.NullString{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !checkPasswordHash(newPasswords.OldPassword, user.Password) {
		http.Error(w, "Incorrect old password", http.StatusUnauthorized)
		return
	}

	user.Password, err = hashPassword(newPasswords.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = u.DB.UpdateUserPassword(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u UsersHandler) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	var id int64
	if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var currentRoles string
	err := u.DB.GetUserRolesByID(id).Scan(&currentRoles)
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

	_, err = u.DB.UpdateUserRoles(newRoles, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newRoles)
}

func (u UsersHandler) UploadProfilePicHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB limit for file size

	file, _, err := r.FormFile("profileImage")
	if err != nil {
		http.Error(w, "Unable to parse file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	uploadDir := "frontend/public/"
	fileName := r.MultipartForm.File["profileImage"][0].Filename
	username := r.FormValue("username")

	var user databasestructs.User
	var profilepic sql.NullString
	err = u.DB.GetUserByUsername(username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &profilepic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.ProfilePic = profilepic.String
	if user.ProfilePic != "" {
		err := os.Remove(uploadDir + user.ProfilePic)
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

	_, err = u.DB.UpdateUserProfilePic(username, fileName)
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

func (u UsersHandler) CreateUserAdmin(w http.ResponseWriter, r *http.Request) {
	var reqUser databasestructs.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := u.createNewUser(reqUser.Username, reqUser.Password, reqUser.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User created successfully"))
}
