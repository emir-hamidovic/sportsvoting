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

func (m *MySqlDB) SelectPlayer(playerid string) *sql.Row {
	return m.db.QueryRow("SELECT id, playerid, name, teamabbr, age from players where playerid=?", playerid)
}

func (m *MySqlDB) SelectTeamByAbbrevation(teamabbr string) *sql.Row {
	return m.db.QueryRow("SELECT teamabbr from teams where teamabbr=?", teamabbr)
}

func (m *MySqlDB) GetMVPStats(ctx context.Context, season string) (*sql.Rows, error) {
	// need pictures for each player
	return m.db.QueryContext(ctx, "SELECT name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetDPOYStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT name, gamesplayed, minutespergame, reboundspergame, stealspergame, blockspergame, position, dws, dbpm, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY dws DESC", season, season)
}

func (m *MySqlDB) GetSixManStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND gamesplayed - gamesstarted > gamesstarted ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetCOYStats(season string) (*sql.Rows, error) {
	// need team stats, team record etc.
	return m.db.Query("SELECT name, gamesplayed, minutespergame, reboundspergame, stealspergame, blockspergame, position, dws, dbpm FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY dws DESC", season, season)
}

func (m *MySqlDB) GetROYStats(ctx context.Context, season string) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ws, bpm, offrtg, defrtg FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND rookieseason=1 AND minutespergame > 10 ORDER BY per DESC", season, season)
}

func (m *MySqlDB) GetMIPStats(season string) (*sql.Rows, error) {
	// need to get previous year stats
	return m.db.Query("SELECT name, gamesplayed, minutespergame, reboundspergame, stealspergame, blockspergame, position, dws, dbpm FROM players INNER JOIN stats ON players.playerid=stats.playerid INNER JOIN advancedstats ON players.playerid=advancedstats.playerid WHERE advancedstats.season=? AND stats.season=? AND minutespergame > 20 ORDER BY dws DESC", season, season)
}
