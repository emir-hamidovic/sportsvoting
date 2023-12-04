package players

import (
	"database/sql"
	"fmt"
	"sportsvoting/advancedstats"
	"sportsvoting/database"
	"sportsvoting/databasestructs"
	"sportsvoting/request"
	"sportsvoting/scraper"
	"sportsvoting/stats"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetPlayerInfo(doc *goquery.Document, team string, playersList map[string]databasestructs.PlayerInfo, season string) (map[string]databasestructs.PlayerInfo, error) {
	playersList, err := getRosterInfo(doc, team, playersList)
	if err != nil {
		return nil, err
	}

	playersList, err = getSeasonPerGameStats(doc, playersList, season)
	if err != nil {
		return nil, err
	}

	rows := doc.Find("table#advanced > tbody > tr")
	playersList, err = getSeasonAdvancedStats(rows, playersList, season)
	if err != nil {
		return nil, err
	}

	playersList, err = getSeasonOffAndDefRtg(playersList, season)
	if err != nil {
		return nil, err
	}

	return playersList, nil
}

func InsertPlayers(db database.Database, players map[string]databasestructs.PlayerInfo) error {
	for _, player := range players {
		_, err := db.InsertPlayer(player)
		if err != nil {
			return err
		}
	}

	fmt.Println("Players added to database.")
	return nil
}
func getRosterInfo(doc *goquery.Document, team string, player map[string]databasestructs.PlayerInfo) (map[string]databasestructs.PlayerInfo, error) {
	rows := doc.Find("table#roster > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		name := scraper.GetTDDataStatString(row, "player")
		id := request.GetPlayerIDFromDocument(row)
		if id != "" {
			college := row.Find("td[data-stat='college']").Last().Text()
			height := scraper.GetTDDataStatString(row, "height")
			weight := scraper.GetTDDataStatString(row, "weight")
			position := scraper.GetTDDataStatString(row, "pos")

			player[id] = databasestructs.PlayerInfo{Name: name, ID: id, College: college, Height: height, Weight: weight, TeamAbbr: team, PlayerStats: databasestructs.PlayerStats{Position: position, PlayerID: id, TeamAbbr: team}, AdvancedStats: databasestructs.AdvancedStats{PlayerID: id, TeamAbbr: team}}
		}
	})

	return player, nil
}

func getSeasonPerGameStats(doc *goquery.Document, player map[string]databasestructs.PlayerInfo, season string) (map[string]databasestructs.PlayerInfo, error) {
	rows := doc.Find("table#per_game > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		pl := getPlayerAge(row)
		if entry, ok := player[pl.ID]; ok {
			entry.Age = pl.Age
			stats.FillPlayerStatsForSeason(row, season, &entry.PlayerStats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}

func getSeasonAdvancedStats(rows *goquery.Selection, player map[string]databasestructs.PlayerInfo, season string) (map[string]databasestructs.PlayerInfo, error) {
	rows.Each(func(i int, row *goquery.Selection) {
		var pl databasestructs.PlayerInfo
		pl.ID = request.GetPlayerIDFromDocument(row)
		if entry, ok := player[pl.ID]; ok {
			advancedstats.FillPlayerStatsForSeason(row, season, &entry.AdvancedStats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}

func getSeasonOffAndDefRtg(player map[string]databasestructs.PlayerInfo, season string) (map[string]databasestructs.PlayerInfo, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_poss.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return nil, err
	}

	rows := doc.Find("table#per_poss_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		id := request.GetPlayerIDFromDocument(row)

		if entry, ok := player[id]; ok {
			entry.AdvancedStats.DefRtg = scraper.GetTDDataStatFloat(row, "def_rtg")
			entry.AdvancedStats.OffRtg = scraper.GetTDDataStatFloat(row, "off_rtg")
			player[id] = entry
		}
	})

	return player, nil
}

func getPlayerAge(row *goquery.Selection) databasestructs.PlayerInfo {
	var player databasestructs.PlayerInfo
	player.ID = request.GetPlayerIDFromDocument(row)
	player.Age = scraper.GetTDDataStatInt(row, "age")
	return player
}

func UpdatePlayersWhoPlayedAGame(db database.Database) error {
	fmt.Println("Updating players who played")
	season := GetEndYearOfTheSeason()

	rows, err := db.SelectPlayerGamesPlayed(season)
	if err != nil {
		return err
	}
	defer rows.Close()

	players := make(map[string]int64, 600)
	for rows.Next() {
		var id string
		var gamesplayed int64
		err := rows.Scan(&id, &gamesplayed)
		if err != nil {
			fmt.Println(err)
		}
		players[id] = gamesplayed
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_game.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	newplayers := make(map[string]databasestructs.PlayerInfo, 500)
	updateplayers := make(map[string]databasestructs.PlayerInfo, 500)
	table := doc.Find("table#per_game_stats > tbody > tr")
	table.Each(func(i int, row *goquery.Selection) {
		id := request.GetPlayerIDFromDocument(row)
		player := getPlayerAge(row)
		player.PlayerStats.PlayerID = id
		player.AdvancedStats.PlayerID = id
		player.Name = scraper.GetTDDataStatString(row, "player")
		player.PlayerStats.Position = scraper.GetTDDataStatString(row, "pos")
		player.TeamAbbr = scraper.GetTDDataStatString(row, "team_id")
		player.PlayerStats.TeamAbbr = player.TeamAbbr
		player.AdvancedStats.TeamAbbr = player.TeamAbbr
		stats.FillPlayerStatsForSeason(row, season, &player.PlayerStats)

		entry, ok := players[id]
		if !ok && id != "" {
			err = db.CheckPlayerExists(id).Scan()
			if err == sql.ErrNoRows {
				newplayers[id] = player
				players[id] = player.Games
			} else {
				err = stats.UpdateStats(db, player.PlayerStats)
				if err != nil {
					fmt.Println(err)
				}
				updateplayers[id] = player
			}
		} else if player.PlayerStats.Games > entry {
			err = stats.UpdateStats(db, player.PlayerStats)
			if err != nil {
				fmt.Println(err)
			}
			updateplayers[id] = player
		}

		if entry, ok := newplayers[id]; ok {
			entry.TeamAbbr = player.TeamAbbr
			entry.PlayerStats.TeamAbbr = player.TeamAbbr
			entry.AdvancedStats.TeamAbbr = player.TeamAbbr
			newplayers[id] = entry
		}
	})

	err = InsertPlayers(db, newplayers)
	if err != nil {
		return err
	}

	url = fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_advanced.html", season)
	docAdvanced, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	tableAdvanced := docAdvanced.Find("table#advanced_stats > tbody > tr")
	newplayers, err = getSeasonAdvancedStats(tableAdvanced, newplayers, season)
	if err != nil {
		fmt.Println(err)
	}
	updateplayers, err = getSeasonAdvancedStats(tableAdvanced, updateplayers, season)
	if err != nil {
		fmt.Println(err)
	}

	newplayers, err = getSeasonOffAndDefRtg(newplayers, season)
	if err != nil {
		fmt.Println(err)
	}
	updateplayers, err = getSeasonOffAndDefRtg(updateplayers, season)
	if err != nil {
		fmt.Println(err)
	}

	UpdatePlayerStats(db, newplayers, season)
	UpdatePlayerStats(db, updateplayers, season)

	return nil
}

func GetEndYearOfTheSeason() string {
	today := time.Now()
	year := today.Year()
	month := int(today.Month())
	var currentSeason string
	if month < 11 {
		currentSeason = fmt.Sprint(year)
	} else {
		currentSeason = fmt.Sprint(year + 1)
	}

	return currentSeason
}

func UpdatePlayerStats(db database.Database, rosters map[string]databasestructs.PlayerInfo, season string) error {
	fmt.Println("Updating stats")
	for _, player := range rosters {
		err := stats.UpdateStats(db, player.PlayerStats)
		if err != nil {
			return err
		}

		err = advancedstats.UpdateStats(db, player.AdvancedStats)
		if err != nil {
			return err
		}
	}

	err := stats.UpdateTradedPlayerStats(db, season)
	if err != nil {
		return err
	}

	err = advancedstats.UpdateTradedPlayerStats(db, season)
	if err != nil {
		return err
	}

	err = stats.SetRookies(db, season)
	if err != nil {
		return err
	}

	return nil
}
