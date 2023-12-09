package mysql_db

import (
	"database/sql"
	"sportsvoting/databasestructs"
)

func (m *MySqlDB) InsertGOATPlayer(info databasestructs.GoatPlayers) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO goat_players(playerid, name, allstar, allnba, alldefense, championships, dpoy, sixman, roy, finalsmvp, mvp, isactive) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", info.ID, info.Name, info.AllStar, info.AllNba, info.AllDefense, info.Championships, info.Dpoy, info.SixMan, info.ROY, info.FMVP, info.MVP, info.IsActive)
}

func (m *MySqlDB) UpdateGOATPlayer(info databasestructs.GoatPlayers) (sql.Result, error) {
	return m.db.Exec("UPDATE goat_players SET allstar=?, allnba=?, alldefense=?, championships=?, dpoy=?, sixman=?, roy=?, finalsmvp=?, mvp=?, isactive=? WHERE playerid=?", info.AllStar, info.AllNba, info.AllDefense, info.Championships, info.Dpoy, info.SixMan, info.ROY, info.FMVP, info.MVP, info.IsActive, info.ID)
}

func (m *MySqlDB) UpdateGOATStats(stats databasestructs.GoatStats) (sql.Result, error) {
	return m.db.Exec("UPDATE goat_stats SET pointspergame=?, reboundspergame=?, assistspergame=?, stealspergame=?, blockspergame=?, turnoverspergame=?, fgpercentage=?, ftpercentage=?, threeptpercentage=?, per=?, ows=?, dws=?, ws=?, obpm=?, dbpm=?, bpm=?, vorp=?, offrtg=?, defrtg=?, totalpoints=?, totalrebounds=?, totalassists=?, totalsteals=?, totalblocks=?, position=? WHERE playerid=? AND isplayoffs=?", stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.PER, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TotalPoints, stats.TotalRebounds, stats.TotalAssists, stats.TotalSteals, stats.TotalBlocks, stats.Position, stats.PlayerID, stats.IsPlayoffs)
}

func (m *MySqlDB) InsertGOATStats(stats databasestructs.GoatStats) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO goat_stats (playerid, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, turnoverspergame, fgpercentage, ftpercentage, threeptpercentage, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg, totalpoints, totalrebounds, totalassists, totalsteals, totalblocks, isplayoffs, position, isactive) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", stats.PlayerID, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.PER, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TotalPoints, stats.TotalRebounds, stats.TotalAssists, stats.TotalSteals, stats.TotalBlocks, stats.IsPlayoffs, stats.Position, stats.IsActive)
}

func (m *MySqlDB) GetGOATStats(season string) (*sql.Rows, error) {
	return m.db.Query("SELECT p.playerid, p.name, sr.position, ROUND(AVG(sr.pointspergame), 1) as avg_pointspergame, ROUND(AVG(sr.reboundspergame), 1) as avg_reboundspergame, ROUND(AVG(sr.assistspergame), 1) as avg_assistspergame, ROUND(AVG(sr.stealspergame), 1) as avg_stealspergame, ROUND(AVG(sr.blockspergame), 1) as avg_blockspergame, p.allstar, p.allnba, p.alldefense, p.championships, p.dpoy, p.finalsmvp, p.mvp, ROUND(AVG(sr.per), 1) as avg_per, ROUND(AVG(sr.ows), 1) as avg_ows, ROUND(AVG(sr.dws), 1) as avg_dws, ROUND(AVG(sr.ws), 1) as avg_ws, ROUND(AVG(sr.dbpm), 1) as avg_dbpm, ROUND(AVG(sr.obpm), 1) as avg_obpm, ROUND(AVG(sr.bpm), 1) as avg_bpm, ROUND(AVG(sr.defrtg), 1) as avg_defrtg, ROUND(AVG(sr.offrtg), 1) as avg_offrtg, ROUND(AVG(sp.pointspergame), 1) as avg_playoff_pointspergame, ROUND(AVG(sp.reboundspergame), 1) as avg_playoff_reboundspergame, ROUND(AVG(sp.assistspergame), 1) as avg_playoff_assistspergame FROM goat_players p INNER JOIN goat_stats sr ON p.playerid=sr.playerid INNER JOIN goat_stats sp ON p.playerid=sp.playerid WHERE sp.isplayoffs=1 GROUP BY p.playerid, p.name, sr.position, p.allstar, p.allnba, p.alldefense, p.championships, p.dpoy, p.finalsmvp, p.mvp ORDER BY avg_per DESC")
}

func (m *MySqlDB) GetActivePlayers() (*sql.Rows, error) {
	return m.db.Query("SELECT playerid FROM goat_players WHERE isactive = 1")
}
