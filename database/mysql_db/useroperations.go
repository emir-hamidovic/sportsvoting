package mysql_db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sportsvoting/databasestructs"

	"golang.org/x/crypto/bcrypt"
)

func (m *MySqlDB) CreateAdminUser() error {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = 1").Scan(&count)
	if err != nil {
		return err
	}

	dbPass := os.Getenv("DBPASS")
	if count == 0 {
		bytes, err := bcrypt.GenerateFromPassword([]byte(dbPass), 14)
		if err != nil {
			return err
		}

		dbPass = string(bytes)
		_, err = m.db.Exec("INSERT INTO users(id, username, password) VALUES (1, 'admin', ?)", dbPass)
		if err != nil {
			return fmt.Errorf("error inserting admin user: %v", err)
		}

		_, err = m.db.Exec("INSERT INTO user_roles (user_id, role) VALUES (?, ?)", 1, "user,admin")
		if err != nil {
			return fmt.Errorf("error inserting user role for admin: %v", err)
		}
	}

	return nil
}

func (m *MySqlDB) DeleteUser(id int64) (sql.Result, error) {
	return m.db.Exec("DELETE FROM users WHERE id=?", id)
}

func (m *MySqlDB) GetAllUsers() (*sql.Rows, error) {
	return m.db.Query("SELECT id, username, email, password, refresh_token, profile_pic FROM users")
}

func (m *MySqlDB) GetUserByUsername(username string) *sql.Row {
	return m.db.QueryRow("SELECT id, username, email, password, refresh_token, profile_pic FROM users WHERE username=?", username)
}

func (m *MySqlDB) GetUserByID(id int64) *sql.Row {
	return m.db.QueryRow("SELECT username, email, profile_pic FROM users WHERE id=?", id)
}

func (m *MySqlDB) GetUserRolesByID(id int64) *sql.Row {
	return m.db.QueryRow("SELECT role FROM user_roles WHERE user_id=?", id)
}

func (m *MySqlDB) InsertUserRoles(role databasestructs.Role) (sql.Result, error) {
	return m.db.Exec("INSERT INTO user_roles (user_id, role) VALUES (?, ?)", role.UserID, role.Role)
}

func (m *MySqlDB) UpdateUserRoles(roles string, user_id int64) (sql.Result, error) {
	return m.db.Exec("UPDATE user_roles SET role=? WHERE user_id=?", roles, user_id)
}

func (m *MySqlDB) GetCurrentProfilePic(id int64) *sql.Row {
	return m.db.QueryRow("SELECT profile_pic FROM users WHERE id=?", id)
}

func (m *MySqlDB) GetUserByRefreshToken(refresh_token string) *sql.Row {
	return m.db.QueryRow("SELECT id, username, email FROM users WHERE refresh_token=?", refresh_token)
}

func (m *MySqlDB) InsertNewUser(user databasestructs.User) (sql.Result, error) {
	return m.db.Exec("INSERT INTO users(username, email, password, refresh_token) VALUES (?, ?, ?, ?)", user.Username, user.Email, user.Password, user.RefreshToken)
}

func (m *MySqlDB) UpdateUserRefreshToken(username, refresh_token string) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET refresh_token=? WHERE username=?", refresh_token, username)
}

func (m *MySqlDB) UpdateUserIsAdmin(username string, is_admin bool) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET is_admin=? WHERE username=?", is_admin, username)
}

func (m *MySqlDB) UpdateUserPassword(username, password string) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET password=? WHERE username=?", password, username)
}

func (m *MySqlDB) UpdateUserEmail(username, email string) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET email=? WHERE username=?", email, username)
}

func (m *MySqlDB) UpdateUserUsername(oldusername, username string) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET username=? WHERE username=?", username, oldusername)
}

func (m *MySqlDB) UpdateUserProfilePic(username, profile_pic string) (sql.Result, error) {
	return m.db.Exec("UPDATE users SET profile_pic=? WHERE username=?", profile_pic, username)
}

func (m *MySqlDB) GetVotesOfUser(ctx context.Context, userid int64) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT po.id, COALESCE(v.playerid, v.goatplayerid) AS playerid, COALESCE(p.name, gp.name) AS player_name, po.name, po.image FROM player_votes v INNER JOIN polls po ON v.pollid = po.id LEFT JOIN players p ON v.playerid = p.playerid LEFT JOIN   goat_players gp ON v.goatplayerid = gp.playerid WHERE v.userid=?", userid)
}
