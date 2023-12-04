package mysql_db

import (
	"database/sql"
	"sportsvoting/databasestructs"
)

func (m *MySqlDB) InsertPlayer(info databasestructs.PlayerInfo) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO players(playerid, name, college, teamabbr, height, weight, age) VALUES (?, ?, ?, ?, ?, ?, ?)", info.ID, info.Name, info.College, info.TeamAbbr, info.Height, info.Weight, info.Age)
}

func (m *MySqlDB) UpdatePlayerAge(playerid string, age int64) (sql.Result, error) {
	return m.db.Exec("UPDATE players set age=? WHERE playerid=?", age, playerid)
}

func (m *MySqlDB) SelectPlayerGamesPlayed(season string) (*sql.Rows, error) {
	return m.db.Query("SELECT playerid, gamesplayed FROM stats WHERE season=?", season)
}

func (m *MySqlDB) CheckPlayerExists(playerid string) *sql.Row {
	return m.db.QueryRow("SELECT 1 FROM players WHERE playerid=?", playerid)
}
