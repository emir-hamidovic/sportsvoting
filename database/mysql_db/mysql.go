package mysql_db

import (
	"context"
	"database/sql"
	"os"

	"github.com/go-sql-driver/mysql"
)

type MySqlDB struct {
	db *sql.DB
}

func NewDB(dbname string, addr string) (*MySqlDB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 addr,
		DBName:               dbname,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, pingErr
	}

	return &MySqlDB{db: db}, nil
}

func (m *MySqlDB) InsertPlayer(playerid, name, college, teamabbr, height, weight string, age int64) (sql.Result, error) {
	return m.db.Exec("INSERT INTO players(playerid, name, college, teamabbr, height, weight, age) VALUES (?, ?, ?, ?, ?, ?, ?)", playerid, name, college, teamabbr, height, weight, age)
}

func (m *MySqlDB) UpdatePlayerAge(playerid string, age int64) (sql.Result, error) {
	return m.db.Exec("UPDATE players set age=? WHERE playerid=?", age, playerid)
}

func (m *MySqlDB) InsertTeam(teamabbr, name, logo string, winlosspct float64, playoffs, divtitles, conftitles, championships int64) (sql.Result, error) {
	return m.db.Exec("INSERT INTO teams(teamabbr, name, logo, winlosspct, playoffs, divisiontitles, conferencetitles, championships) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", teamabbr, name, logo, winlosspct, playoffs, divtitles, conftitles, championships)
}

func (m *MySqlDB) UpdateStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, gamesstarted=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=?, teamabbr=? WHERE playerid=? AND season=?", gp, gs, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, position, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTradedPlayerStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, gamesstarted=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=? WHERE playerid=? AND season=?", gp, gs, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, position, playerid, season)
}

func (m *MySqlDB) InsertStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error) {
	return m.db.Exec("INSERT INTO stats (gamesplayed, gamesstarted, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage, season, position, playerid, teamabbr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", gp, gs, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, season, position, playerid, teamabbr)
}

func (m *MySqlDB) SetRookieStatus(id string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats set rookieseason=1 WHERE playerid=?", id)
}

func (m *MySqlDB) UpdateAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg float64, teamabbr, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=?, offrtg=?, defrtg=?, teamabbr=? WHERE playerid=? AND season=?", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTradedPlayerAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=? WHERE playerid=? AND season=?", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, playerid, season)
}

func (m *MySqlDB) UpdateOffAndDefRtg(offrtg, defrtg float64, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET offrtg=?, defrtg=? WHERE playerid=? AND season=?", offrtg, defrtg, playerid, season)
}

func (m *MySqlDB) InsertAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg float64, teamabbr, playerid, season string) (sql.Result, error) {
	return m.db.Exec("INSERT INTO advancedstats (per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg, teamabbr, playerid, season) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE players set teamabbr=? WHERE playerid=?", teamabbr, playerid)
}

func (m *MySqlDB) SelectPlayerGamesPlayed(season string) (*sql.Rows, error) {
	return m.db.Query("SELECT playerid, gamesplayed FROM stats WHERE season=?", season)
}

func (m *MySqlDB) CheckPlayerExists(playerid string) *sql.Row {
	return m.db.QueryRow("SELECT 1 FROM players WHERE playerid=?", playerid)
}

func (m *MySqlDB) SelectTeamByAbbrevation(teamabbr string) *sql.Row {
	return m.db.QueryRow("SELECT teamabbr from teams where teamabbr=?", teamabbr)
}

func (m *MySqlDB) GetMVPStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT players.playerid, name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetDPOYStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT players.playerid, name, gamesplayed, minutespergame, reboundspergame, stealspergame, blockspergame, position, dws, dbpm, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY dws DESC", season, season)
}

func (m *MySqlDB) GetSixManStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT players.playerid, name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND gamesplayed - gamesstarted > gamesstarted ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetROYStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT players.playerid, name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ws, bpm, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND rookieseason=1 AND minutespergame > 10 ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetPolls(ctx context.Context) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT id, name, description, image, endpoint FROM polls")
}

func (m *MySqlDB) InsertPolls(id int64, name, description, image, endpoint string) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO polls(id, name, description, image, endpoint) VALUES (?, ?, ?, ?, ?)", id, name, description, image, endpoint)
}

func (m *MySqlDB) GetPlayerPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT p.name, v.votes_for, po.name FROM player_votes v INNER JOIN players p ON v.playerid=p.playerid INNER JOIN polls po ON v.pollid=po.id WHERE v.pollid=? ORDER BY v.votes_for DESC", pollid)
}

func (m *MySqlDB) InsertPlayerVotes(pollid int64, playerid string) (sql.Result, error) {
	var id int64
	rows := m.db.QueryRow("SELECT id FROM player_votes WHERE pollid=? AND playerid=?", pollid, playerid)
	err := rows.Scan(&id)
	if err == sql.ErrNoRows {
		return m.db.Exec("INSERT IGNORE INTO player_votes(playerid, pollid, votes_for) VALUES (?, ?, 1)", playerid, pollid)
	} else {
		return m.db.Exec("UPDATE player_votes SET votes_for = votes_for + 1 WHERE id=?", id)
	}
}

func (m *MySqlDB) GetTeamPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT t.name, v.votes_for FROM team_votes v INNER JOIN teams t ON v.teamabbr=t.teamabbr INNER JOIN polls po ON v.pollid=po.id WHERE v.pollid=? ORDER BY v.votes_for DESC", pollid)
}

func (m *MySqlDB) GetUserByUsername(username string) *sql.Row {
	return m.db.QueryRow("SELECT username, email, password, refresh_token, profile_pic, is_admin FROM users WHERE username=?", username)
}

func (m *MySqlDB) GetUserByID(id int64) *sql.Row {
	return m.db.QueryRow("SELECT username, email, profile_pic, is_admin FROM users WHERE id=?", id)
}

func (m *MySqlDB) GetUserByRefreshToken(refresh_token string) *sql.Row {
	return m.db.QueryRow("SELECT username FROM users WHERE refresh_token=?", refresh_token)
}

func (m *MySqlDB) InsertNewUser(username, email, password, refresh_token string, is_admin bool) (sql.Result, error) {
	return m.db.Exec("INSERT INTO users(username, email, password, refresh_token, is_admin) VALUES (?, ?, ?, ?)", username, email, password, refresh_token, is_admin)
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

func (m *MySqlDB) DeleteUser(id int64) (sql.Result, error) {
	return m.db.Exec("DELETE FROM users WHERE id=?", id)
}

func (m *MySqlDB) GetAllUsers() (*sql.Rows, error) {
	return m.db.Query("SELECT id, username, email, password, refresh_token, profile_pic, is_admin FROM users")
}
