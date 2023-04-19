package database

import (
	"database/sql"
	"errors"
	"scraper/database/mysql_db"
)

const (
	MYSQL string = "mysql"
)

type Database interface {
	InsertPlayer(playerid, name, college, teamabbr, height, weight string, age int64) (sql.Result, error)
	UpdatePlayerAge(playerid string, age int64) (sql.Result, error)
	InsertTeam(teamabbr, name, logo string, winlosspct float64, playoffs, divtitles, conftitles, championships int64) (sql.Result, error)
	UpdateStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error)
	InsertStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error)
	UpdateTradedPlayerStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid string) (sql.Result, error)
	UpdateAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, teamabbr, playerid, season string) (sql.Result, error)
	UpdateTradedPlayerAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, season, playerid string) (sql.Result, error)
	InsertAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, season, playerid, teamabbr string) (sql.Result, error)
	UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error)
	SelectPlayer(playerid string) *sql.Row
	SelectTeamByAbbrevation(teamabbr string) *sql.Row
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
