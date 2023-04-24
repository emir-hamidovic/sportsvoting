package players

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"scraper/parser/advancedstats"
	"scraper/parser/stats"
	"strconv"
	"strings"
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

	playersList, err = getCurrentSeasonAdvancedStats(doc, playersList)
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
		name := row.Find("td[data-stat='player']").Text()
		id := parser.GetPlayerIDFromDocument(row)
		if id != "" {
			fmt.Printf("%s: %s\n", id, name)

			college := row.Find("td[data-stat='college']").Last().Text()
			height := row.Find("td[data-stat='height']").Text()
			weight := row.Find("td[data-stat='weight']").Text()
			position := row.Find("td[data-stat='pos']").Text()

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
func getCurrentSeasonAdvancedStats(doc *goquery.Document, player map[string]Player) (map[string]Player, error) {
	rows := doc.Find("table#advanced > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		var pl Player
		pl.ID = parser.GetPlayerIDFromDocument(row)
		if entry, ok := player[pl.ID]; ok {
			advancedstats.FillPlayerStatsForSeason(row, GetEndYearOfTheSeason(), &entry.AdvancedStats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}

func getCurrentSeasonOffAndDefRtg(player map[string]Player) (map[string]Player, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_poss.html", GetEndYearOfTheSeason())
	doc, err := parser.GetDocumentFromURL(url)
	if err != nil {
		return nil, err
	}

	rows := doc.Find("table#per_poss_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		id := parser.GetPlayerIDFromDocument(row)

		if entry, ok := player[id]; ok {
			entry.AdvancedStats.DefRtg, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='def_rtg']").Text()), 64)
			entry.AdvancedStats.OffRtg, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='off_rtg']").Text()), 64)
			player[id] = entry
		}
	})

	return player, nil
}

func getPlayerAge(row *goquery.Selection) Player {
	var player Player
	player.ID = parser.GetPlayerIDFromDocument(row)
	player.Age, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='age']").Text()), 10, 64)
	return player
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
