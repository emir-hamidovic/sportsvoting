package database

import (
	"context"
	"database/sql"
	"errors"
	"sportsvoting/database/mysql_db"
)

const (
	MYSQL string = "mysql"
)

type Database interface {
	InsertSeasonEntered(season string) (sql.Result, error)
	SelectSeasonsAvailable() (*sql.Rows, error)
	InsertPlayer(playerid, name, college, teamabbr, height, weight string, age int64) (sql.Result, error)
	UpdatePlayerAge(playerid string, age int64) (sql.Result, error)
	InsertTeam(teamabbr, name, logo string, winlosspct float64, playoffs, divtitles, conftitles, championships int64) (sql.Result, error)
	UpdateStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error)
	InsertStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error)
	UpdateTradedPlayerStats(gp, gs int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid string) (sql.Result, error)
	UpdateAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg float64, teamabbr, playerid, season string) (sql.Result, error)
	UpdateTradedPlayerAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, season, playerid string) (sql.Result, error)
	InsertAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg float64, season, playerid, teamabbr string) (sql.Result, error)
	UpdateOffAndDefRtg(offrtg, defrtg float64, playerid, season string) (sql.Result, error)
	UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error)
	SelectPlayerGamesPlayed(season string) (*sql.Rows, error)
	CheckPlayerExists(playerid string) *sql.Row
	SelectTeamByAbbrevation(teamabbr string) *sql.Row
	GetPlayerStatsForQuiz(ctx context.Context, season string) (*sql.Rows, error)
	GetSixManStats(ctx context.Context, season string) (*sql.Rows, error)
	GetDPOYStats(ctx context.Context, season string) (*sql.Rows, error)
	GetROYStats(ctx context.Context, season string) (*sql.Rows, error)
	SetRookieStatus(id string) (sql.Result, error)
	GetPolls(ctx context.Context) (*sql.Rows, error)
	GetPollByID(id int64) *sql.Row
	InsertPolls(name, description, image, selected_stats, season string) (sql.Result, error)
	InsertPollsWithId(id int64, name, description, image, selected_stats, season string) (sql.Result, error)
	GetPlayerPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error)
	InsertPlayerVotes(pollid, userid int64, playerid string) (sql.Result, error)
	GetTeamPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error)
	GetUserByUsername(username string) *sql.Row
	GetUserByRefreshToken(refresh_token string) *sql.Row
	GetUserByID(id int64) *sql.Row
	GetUserRolesByID(id int64) *sql.Row
	InsertUserRoles(userid int64, role string) (sql.Result, error)
	UpdateUserRoles(roles string, user_id int64) (sql.Result, error)
	GetCurrentProfilePic(id int64) *sql.Row
	InsertNewUser(username, email, password, refresh_token string) (sql.Result, error)
	UpdateUserRefreshToken(username, refresh_token string) (sql.Result, error)
	UpdateUserIsAdmin(username string, is_admin bool) (sql.Result, error)
	UpdateUserPassword(username, password string) (sql.Result, error)
	UpdateUserEmail(username, email string) (sql.Result, error)
	UpdateUserUsername(oldusername, username string) (sql.Result, error)
	UpdateUserProfilePic(username, profile_pic string) (sql.Result, error)
	DeleteUser(id int64) (sql.Result, error)
	GetAllUsers() (*sql.Rows, error)
}

type Config struct {
	DbType string
	DbName string
	Addr   string
}

func NewDB(conf Config) (Database, error) {
	switch conf.DbType {
	case MYSQL:
		return mysql_db.NewDB(conf.DbName, conf.Addr)
	}

	return nil, errors.New("incorrect db type entered")
}
