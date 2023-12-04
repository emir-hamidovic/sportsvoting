package mysql_db

import (
	"context"
	"database/sql"
	"sportsvoting/databasestructs"
)

func (m *MySqlDB) UpdateStats(stats databasestructs.PlayerStats) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, gamesstarted=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=?, teamabbr=? WHERE playerid=? AND season=?", stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Position, stats.TeamAbbr, stats.PlayerID, stats.Season)
}

func (m *MySqlDB) UpdateTradedPlayerStats(stats databasestructs.PlayerStats) (sql.Result, error) {
	return m.db.Exec("UPDATE stats SET gamesplayed=?, gamesstarted=?, minutespergame=?, pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, position=? WHERE playerid=? AND season=?", stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Position, stats.PlayerID, stats.Season)
}

func (m *MySqlDB) InsertStats(stats databasestructs.PlayerStats) (sql.Result, error) {
	return m.db.Exec("INSERT INTO stats (gamesplayed, gamesstarted, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage, season, position, playerid, teamabbr) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
}

func (m *MySqlDB) UpdateAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=?, offrtg=?, defrtg=?, teamabbr=? WHERE playerid=? AND season=?", stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
}

func (m *MySqlDB) UpdateTradedPlayerAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET per=?, tspct=?, usgpct=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=? WHERE playerid=? AND season=?", stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.PlayerID, stats.Season)
}

func (m *MySqlDB) UpdateOffAndDefRtg(offrtg, defrtg float64, playerid, season string) (sql.Result, error) {
	return m.db.Exec("UPDATE advancedstats SET offrtg=?, defrtg=? WHERE playerid=? AND season=?", offrtg, defrtg, playerid, season)
}

func (m *MySqlDB) InsertAdvancedStats(stats databasestructs.AdvancedStats) (sql.Result, error) {
	return m.db.Exec("INSERT INTO advancedstats (per, tspct, usgpct, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg, teamabbr, playerid, season) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
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

func (m *MySqlDB) SetRookieStatus(id string) (sql.Result, error) {
	return m.db.Exec("UPDATE stats set rookieseason=1 WHERE playerid=?", id)
}
