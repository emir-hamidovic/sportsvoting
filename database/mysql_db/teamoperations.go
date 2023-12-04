package mysql_db

import (
	"database/sql"
	"sportsvoting/databasestructs"
)

func (m *MySqlDB) InsertTeam(info databasestructs.TeamInfo) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO teams(teamabbr, name, logo, winlosspct, playoffs, divisiontitles, conferencetitles, championships) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", info.TeamAbbr, info.Name, info.Logo, info.WinLossPct, info.Playoffs, info.DivisionTitles, info.ConferenceTitles, info.Championships)
}

func (m *MySqlDB) UpdateTeamForPlayer(teamabbr, playerid string) (sql.Result, error) {
	return m.db.Exec("UPDATE players set teamabbr=? WHERE playerid=?", teamabbr, playerid)
}

func (m *MySqlDB) SelectTeamByAbbrevation(teamabbr string) *sql.Row {
	return m.db.QueryRow("SELECT teamabbr from teams where teamabbr=?", teamabbr)
}
