package players

import (
	"database/sql"
	"fmt"
	"sportsvoting/advancedstats"
	"sportsvoting/database"
	"sportsvoting/request"
	"sportsvoting/scraper"
	"sportsvoting/stats"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name                        string `json:"name,omitempty"`
	ID                          string `json:"playerid,omitempty"`
	College                     string `json:"college,omitempty"`
	TeamAbbr                    string `json:"team,omitempty"`
	Height                      string `json:"height,omitempty"`
	Weight                      string `json:"weight,omitempty"`
	Age                         int64  `json:"age,omitempty"`
	stats.Stats                 `json:"stats,omitempty"`
	advancedstats.AdvancedStats `json:"advstats,omitempty"`
}

func GetPlayerInfo(doc *goquery.Document, team string, playersList map[string]Player) (map[string]Player, error) {
	playersList, err := getRosterInfo(doc, team, playersList)
	if err != nil {
		return nil, err
	}

	playersList, err = getCurrentSeasonPerGameStats(doc, playersList)
	if err != nil {
		return nil, err
	}

	rows := doc.Find("table#advanced > tbody > tr")
	playersList, err = getCurrentSeasonAdvancedStats(rows, playersList)
	if err != nil {
		return nil, err
	}

	playersList, err = getCurrentSeasonOffAndDefRtg(playersList)
	if err != nil {
		return nil, err
	}

	return playersList, nil
}

func InsertPlayers(db database.Database, players map[string]Player) error {
	for _, player := range players {
		_, err := db.InsertPlayer(player.ID, player.Name, player.College, player.TeamAbbr, player.Height, player.Weight, player.Age)
		if err != nil {
			return err
		}
	}

	fmt.Println("Players added to database.")
	return nil
}
func getRosterInfo(doc *goquery.Document, team string, player map[string]Player) (map[string]Player, error) {
	rows := doc.Find("table#roster > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		name := scraper.GetTDDataStatString(row, "player")
		id := request.GetPlayerIDFromDocument(row)
		if id != "" {
			college := row.Find("td[data-stat='college']").Last().Text()
			height := scraper.GetTDDataStatString(row, "height")
			weight := scraper.GetTDDataStatString(row, "weight")
			position := scraper.GetTDDataStatString(row, "pos")

			player[id] = Player{Name: name, ID: id, College: college, Height: height, Weight: weight, TeamAbbr: team, Stats: stats.Stats{Position: position, PlayerID: id, TeamAbbr: team}, AdvancedStats: advancedstats.AdvancedStats{PlayerID: id, TeamAbbr: team}}
		}
	})

	return player, nil
}

func getCurrentSeasonPerGameStats(doc *goquery.Document, player map[string]Player) (map[string]Player, error) {
	rows := doc.Find("table#per_game > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		pl := getPlayerAge(row)
		if entry, ok := player[pl.ID]; ok {
			entry.Age = pl.Age
			stats.FillPlayerStatsForSeason(row, GetEndYearOfTheSeason(), &entry.Stats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}
func getCurrentSeasonAdvancedStats(rows *goquery.Selection, player map[string]Player) (map[string]Player, error) {
	rows.Each(func(i int, row *goquery.Selection) {
		var pl Player
		pl.ID = request.GetPlayerIDFromDocument(row)
		if entry, ok := player[pl.ID]; ok {
			advancedstats.FillPlayerStatsForSeason(row, GetEndYearOfTheSeason(), &entry.AdvancedStats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}

func getCurrentSeasonOffAndDefRtg(player map[string]Player) (map[string]Player, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_poss.html", GetEndYearOfTheSeason())
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

func getPlayerAge(row *goquery.Selection) Player {
	var player Player
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

	newplayers := make(map[string]Player, 500)
	updateplayers := make(map[string]Player, 500)
	table := doc.Find("table#per_game_stats > tbody > tr")
	table.Each(func(i int, row *goquery.Selection) {
		id := request.GetPlayerIDFromDocument(row)
		player := getPlayerAge(row)
		player.Stats.PlayerID = id
		player.AdvancedStats.PlayerID = id
		player.Name = scraper.GetTDDataStatString(row, "player")
		player.Stats.Position = scraper.GetTDDataStatString(row, "pos")
		player.TeamAbbr = scraper.GetTDDataStatString(row, "team_id")
		player.Stats.TeamAbbr = player.TeamAbbr
		player.AdvancedStats.TeamAbbr = player.TeamAbbr
		stats.FillPlayerStatsForSeason(row, season, &player.Stats)

		entry, ok := players[id]
		if !ok && id != "" {
			err = db.CheckPlayerExists(id).Scan()
			if err == sql.ErrNoRows {
				newplayers[id] = player
				players[id] = player.Games
			} else {
				err = stats.UpdateStats(db, player.Stats)
				if err != nil {
					fmt.Println(err)
				}
				updateplayers[id] = player
			}
		} else if player.Stats.Games > entry {
			err = stats.UpdateStats(db, player.Stats)
			if err != nil {
				fmt.Println(err)
			}
			updateplayers[id] = player
		}

		if entry, ok := newplayers[id]; ok {
			entry.TeamAbbr = player.TeamAbbr
			entry.Stats.TeamAbbr = player.TeamAbbr
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
	newplayers, err = getCurrentSeasonAdvancedStats(tableAdvanced, newplayers)
	if err != nil {
		fmt.Println(err)
	}
	updateplayers, err = getCurrentSeasonAdvancedStats(tableAdvanced, updateplayers)
	if err != nil {
		fmt.Println(err)
	}

	newplayers, err = getCurrentSeasonOffAndDefRtg(newplayers)
	if err != nil {
		fmt.Println(err)
	}
	updateplayers, err = getCurrentSeasonOffAndDefRtg(updateplayers)
	if err != nil {
		fmt.Println(err)
	}

	UpdatePlayerStats(db, newplayers)
	UpdatePlayerStats(db, updateplayers)

	return nil
}

func GetEndYearOfTheSeason() string {
	today := time.Now()
	year := today.Year()
	month := int(today.Month())
	var currentSeason string
	if month < 10 {
		currentSeason = fmt.Sprint(year)
	} else {
		currentSeason = fmt.Sprint(year + 1)
	}

	return currentSeason
}

func UpdatePlayerStats(db database.Database, rosters map[string]Player) error {
	fmt.Println("Updating stats")
	for _, player := range rosters {
		err := stats.UpdateStats(db, player.Stats)
		if err != nil {
			return err
		}

		err = advancedstats.UpdateStats(db, player.AdvancedStats)
		if err != nil {
			return err
		}
	}

	err := stats.UpdateTradedPlayerStats(db, GetEndYearOfTheSeason())
	if err != nil {
		return err
	}

	err = advancedstats.UpdateTradedPlayerStats(db, GetEndYearOfTheSeason())
	if err != nil {
		return err
	}

	err = stats.SetRookies(db, GetEndYearOfTheSeason())
	if err != nil {
		return err
	}

	return nil
}
