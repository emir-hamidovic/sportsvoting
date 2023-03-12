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

func (m *MySqlDB) InsertPlayer(playerid, name, college, teamabbr, height, weight string, age, gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, threept float64) (sql.Result, error) {
	return m.db.Exec("INSERT INTO players(playerid, name, college, teamabbr, height, weight, age, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", playerid, name, college, teamabbr, height, weight, age, gp, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, threept)
}

func (m *MySqlDB) InsertTeam(teamabbr, name, logo string, winlosspct float64, playoffs, divtitles, conftitles, championships int64) (sql.Result, error) {
	return m.db.Exec("INSERT INTO teams(teamabbr, name, logo, winlosspct, playoffs, divisiontitles, conferencetitles, championships) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", teamabbr, name, logo, winlosspct, playoffs, divtitles, conftitles, championships)
}

func (m *MySqlDB) UpdatePlayers(gp int64, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three float64, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE players SET gamesplayed=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=? WHERE playerid=?", gp, mpg, ppg, rpg, apg, spg, bpg, tpg, fg, ft, three, playerid)
}

func (m *MySqlDB) SelectPlayerID(abbr string) (*sql.Rows, error) {
	return m.db.Query("SELECT playerid FROM players WHERE teamabbr = ?", abbr)
}
