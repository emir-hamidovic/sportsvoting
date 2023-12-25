package database

import (
	"context"
	"database/sql"
	"errors"
	"sportsvoting/database/mysql_db"
	"sportsvoting/databasestructs"
	"time"
)

const (
	MYSQL string = "mysql"
)

type Database interface {
	PlayerOperations
	TeamOperations
	StatsOperations
	PollOperations
	UserOperations
	GoatOperations
	SyncOperations
}

type PlayerOperations interface {
	InsertPlayer(info databasestructs.PlayerInfo) (sql.Result, error)
	UpdatePlayerAge(playerid string, age int64) (sql.Result, error)
	SelectPlayerGamesPlayed(season string) (*sql.Rows, error)
	CheckPlayerExists(playerid string) *sql.Row
}

type TeamOperations interface {
	InsertTeam(info databasestructs.TeamInfo) (sql.Result, error)
	UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error)
	SelectTeamByAbbrevation(teamabbr string) *sql.Row
}

type StatsOperations interface {
	UpdateStats(stats databasestructs.PlayerStats) (sql.Result, error)
	InsertStats(stats databasestructs.PlayerStats) (sql.Result, error)
	UpdateTradedPlayerStats(stats databasestructs.PlayerStats) (sql.Result, error)
	UpdateAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error)
	UpdateTradedPlayerAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error)
	InsertAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error)
	UpdateOffAndDefRtg(offrtg, defrtg float64, playerid, season string) (sql.Result, error)
	GetSixManStats(ctx context.Context, season string) (*sql.Rows, error)
	GetDPOYStats(ctx context.Context, season string) (*sql.Rows, error)
	GetROYStats(ctx context.Context, season string) (*sql.Rows, error)
	SetRookieStatus(id string) (sql.Result, error)
}

type PollOperations interface {
	GetPolls(ctx context.Context) (*sql.Rows, error)
	GetPollByID(id int64) *sql.Row
	GetPollByUserID(userid int64) (*sql.Rows, error)
	InsertPolls(poll databasestructs.Poll) (sql.Result, error)
	InsertPollsWithId(poll databasestructs.Poll) (sql.Result, error)
	DeletePollByID(pollid int64) (sql.Result, error)
	ResetPollVotes(pollid int64) (sql.Result, error)
	UpdatePollByID(poll databasestructs.Poll) (sql.Result, error)
	InsertSeasonEntered(season string) (sql.Result, error)
	SelectSeasonsAvailable() (*sql.Rows, error)
	SelectSeasonsForNonGOATStats() (*sql.Rows, error)
	GetPlayerStatsForPoll(ctx context.Context, season string) (*sql.Rows, error)
	GetPlayerPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error)
	InsertPlayerVotes(pollid, userid int64, playerid string) (sql.Result, error)
	GetTeamPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error)
	UpdatePollImage(image databasestructs.Image) (sql.Result, error)
}

type UserOperations interface {
	GetUserByUsername(username string) *sql.Row
	GetUserByRefreshToken(refresh_token string) *sql.Row
	GetUserByID(id int64) *sql.Row
	GetUserRolesByID(id int64) *sql.Row
	InsertUserRoles(role databasestructs.Role) (sql.Result, error)
	UpdateUserRoles(roles string, user_id int64) (sql.Result, error)
	InsertNewUser(user databasestructs.User) (sql.Result, error)
	UpdateUserRefreshToken(username, refresh_token string) (sql.Result, error)
	UpdateUserIsAdmin(username string, is_admin bool) (sql.Result, error)
	UpdateUserPassword(username, password string) (sql.Result, error)
	UpdateUserEmail(username, email string) (sql.Result, error)
	UpdateUserUsername(oldusername, username string) (sql.Result, error)
	UpdateUserProfilePic(username, profile_pic string) (sql.Result, error)
	DeleteUser(id int64) (sql.Result, error)
	GetAllUsers() (*sql.Rows, error)
	GetVotesOfUser(ctx context.Context, userid int64) (*sql.Rows, error)
	CreateAdminUser() error
	GetCurrentProfilePic(id int64) *sql.Row
}

type GoatOperations interface {
	InsertGOATPlayer(info databasestructs.GoatPlayers) (sql.Result, error)
	UpdateGOATPlayer(info databasestructs.GoatPlayers) (sql.Result, error)
	UpdateGOATStats(stats databasestructs.GoatStats) (sql.Result, error)
	InsertGOATStats(stats databasestructs.GoatStats) (sql.Result, error)
	GetGOATStats() (*sql.Rows, error)
	GetActivePlayers() (*sql.Rows, error)
}

type SyncOperations interface {
	GetLastSyncTimeFromDB() (time.Time, error)
	UpdateLastSyncTimeInDB(newTime time.Time) error
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
