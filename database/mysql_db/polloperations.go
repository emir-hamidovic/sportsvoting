package mysql_db

import (
	"context"
	"database/sql"
	"sportsvoting/databasestructs"
	"strings"
)

func (m *MySqlDB) GetPlayerStatsForPoll(ctx context.Context, season string) (*sql.Rows, error) {
	query := `
        SELECT players.playerid, name, gamesplayed, minutespergame, pointspergame, reboundspergame, assistspergame, stealspergame, blockspergame, fgpercentage, threeptpercentage, ftpercentage, turnoverspergame, position, per, ows, dws, ws, obpm, dbpm, bpm, vorp, offrtg, defrtg
        FROM players
        INNER JOIN stats ON players.playerid = stats.playerid
        INNER JOIN advancedstats ON players.playerid = advancedstats.playerid
        WHERE advancedstats.season = ? AND stats.season = ? AND minutespergame > 20
        ORDER BY per DESC`
	return m.db.QueryContext(ctx, query, season, season)
}

func (m *MySqlDB) GetPolls(ctx context.Context) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT id, name, description, image, selected_stats, season, userid FROM polls")
}

func (m *MySqlDB) GetPollByID(id int64) *sql.Row {
	return m.db.QueryRow("SELECT name, description, image, selected_stats, season, userid FROM polls WHERE id=?", id)
}

func (m *MySqlDB) GetPollByUserID(userid int64) (*sql.Rows, error) {
	return m.db.Query("SELECT id, name, description, image, selected_stats, season FROM polls WHERE userid=?", userid)
}

func (m *MySqlDB) InsertPolls(poll databasestructs.Poll) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO polls(name, description, image, selected_stats, season, userid) VALUES (?, ?, ?, ?, ?, ?)", poll.Name, poll.Description, poll.Image, poll.SelectedStats, poll.Season, poll.UserID)
}

func (m *MySqlDB) InsertPollsWithId(poll databasestructs.Poll) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO polls(id, name, description, image, selected_stats, season, userid) VALUES (?, ?, ?, ?, ?, ?, ?)", poll.ID, poll.Name, poll.Description, poll.Image, poll.SelectedStats, poll.Season, poll.UserID)
}

func (m *MySqlDB) GetPlayerPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error) {
	var stats string
	err := m.db.QueryRow("SELECT selected_stats FROM polls WHERE id=?", pollid).Scan(&stats)
	if err != nil {
		return nil, err
	}

	if strings.Contains(stats, "GOAT") {
		return m.db.QueryContext(ctx, "SELECT p.name, COUNT(v.votes_for) as votes_for, po.name FROM player_votes v INNER JOIN goat_players p ON v.goatplayerid=p.playerid INNER JOIN polls po ON v.pollid=po.id WHERE v.pollid=? GROUP BY p.name, po.name ORDER BY COUNT(v.votes_for) DESC", pollid)
	}

	return m.db.QueryContext(ctx, "SELECT p.name, COUNT(v.votes_for) as votes_for, po.name FROM player_votes v INNER JOIN players p ON v.playerid=p.playerid INNER JOIN polls po ON v.pollid=po.id WHERE v.pollid=? GROUP BY p.name, po.name ORDER BY COUNT(v.votes_for) DESC", pollid)
}

func (m *MySqlDB) InsertPlayerVotes(pollid, userid int64, playerid string) (sql.Result, error) {
	var id int64
	var playerIdDB string

	var stats string
	err := m.db.QueryRow("SELECT selected_stats FROM polls WHERE id=?", pollid).Scan(&stats)
	if err != nil {
		return nil, err
	}

	if strings.Contains(stats, "GOAT") {
		rows := m.db.QueryRow("SELECT id, goatplayerid FROM player_votes WHERE pollid=? AND userid=?", pollid, userid)
		err := rows.Scan(&id, &playerIdDB)
		if err == sql.ErrNoRows {
			return m.db.Exec("INSERT IGNORE INTO player_votes(goatplayerid, pollid, userid, votes_for) VALUES (?, ?, ?, 1)", playerid, pollid, userid)
		} else if playerIdDB != playerid {
			m.db.Exec("DELETE FROM player_votes WHERE id=?", id)
			return m.db.Exec("INSERT IGNORE INTO player_votes(goatplayerid, pollid, userid, votes_for) VALUES (?, ?, ?, 1)", playerid, pollid, userid)
		}
	} else {
		rows := m.db.QueryRow("SELECT id, playerid FROM player_votes WHERE pollid=? AND userid=?", pollid, userid)
		err := rows.Scan(&id, &playerIdDB)
		if err == sql.ErrNoRows {
			return m.db.Exec("INSERT IGNORE INTO player_votes(playerid, pollid, userid, votes_for) VALUES (?, ?, ?, 1)", playerid, pollid, userid)
		} else if playerIdDB != playerid {
			m.db.Exec("DELETE FROM player_votes WHERE id=?", id)
			return m.db.Exec("INSERT IGNORE INTO player_votes(playerid, pollid, userid, votes_for) VALUES (?, ?, ?, 1)", playerid, pollid, userid)
		}
	}

	return nil, nil
}

func (m *MySqlDB) DeletePollByID(pollid int64) (sql.Result, error) {
	return m.db.Exec("DELETE FROM polls WHERE id=?", pollid)
}

func (m *MySqlDB) UpdatePollByID(poll databasestructs.Poll) (sql.Result, error) {
	return m.db.Exec("UPDATE polls SET name=?, description=?, selected_stats=?, season=? WHERE id=?", poll.Name, poll.Description, poll.SelectedStats, poll.Season, poll.ID)
}

func (m *MySqlDB) ResetPollVotes(pollid int64) (sql.Result, error) {
	return m.db.Exec("DELETE FROM player_votes WHERE pollid=?", pollid)
}

func (m *MySqlDB) GetTeamPollVotes(ctx context.Context, pollid int64) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, "SELECT t.name, v.votes_for FROM team_votes v INNER JOIN teams t ON v.teamabbr=t.teamabbr INNER JOIN polls po ON v.pollid=po.id WHERE v.pollid=? ORDER BY v.votes_for DESC", pollid)
}

func (m *MySqlDB) UpdatePollImage(image databasestructs.Image) (sql.Result, error) {
	return m.db.Exec("UPDATE polls SET image=? WHERE id=?", image.ImageURL, image.ID)
}

func (m *MySqlDB) InsertSeasonEntered(season string) (sql.Result, error) {
	return m.db.Exec("INSERT IGNORE INTO seasons_entered(season) VALUES (?)", season)
}

func (m *MySqlDB) SelectSeasonsAvailable() (*sql.Rows, error) {
	return m.db.Query("SELECT season FROM seasons_entered")
}

func (m *MySqlDB) SelectSeasonsForNonGOATStats() (*sql.Rows, error) {
	return m.db.Query("SELECT season FROM seasons_entered WHERE season NOT IN ('All', 'Playoffs', 'Career')")
}
