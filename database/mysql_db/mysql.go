package mysql_db

import (
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

func (m *MySqlDB) UpdateStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=?, teamabbr=? WHERE playerid=? AND season=?", gp, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, position, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTradedPlayerStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=? WHERE playerid=? AND season=?", gp, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, position, playerid, season)
}

func (m *MySqlDB) InsertStats(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, season, position, playerid, teamabbr string) (sql.Result, error) {
	return m.db.Exec("INSERT INTO stats (gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage, season, position, playerid, teamabbr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", gp, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, season, position, playerid, teamabbr)
}

func (m *MySqlDB) UpdateAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, teamabbr, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=?, teamabbr=? WHERE playerid=? AND season=?", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTradedPlayerAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=? WHERE playerid=? AND season=?", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, playerid, season)
}

func (m *MySqlDB) InsertAdvancedStats(per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp float64, teamabbr, playerid, season string) (sql.Result, error) {
	return m.db.Exec("INSERT INTO advancedstats (per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, teamabbr, playerid, season) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, teamabbr, playerid, season)
}

func (m *MySqlDB) UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE players set teamabbr=? WHERE playerid=?", teamabbr, playerid)
}

func (m *MySqlDB) SelectPlayer(playerid string) *sql.Row {
	return m.db.QueryRow("SELECT id, playerid, name, teamabbr, age from players where playerid=?", playerid)
}
func (m *MySqlDB) SelectTeamByAbbrevation(teamabbr string) *sql.Row {
	return m.db.QueryRow("SELECT teamabbr from teams where teamabbr=?", teamabbr)
}
