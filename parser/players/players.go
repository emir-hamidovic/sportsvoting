package players

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"scraper/parser/stats"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name     string
	ID       string
	College  string
	TeamAbbr string
	Height   string
	Weight   string
	Age      int64
	stats.Stats
}

func InsertPlayers(db database.Database, players []Player) error {
	for _, player := range players {
		_, err := db.InsertPlayer(player.ID, player.Name, player.College, player.TeamAbbr, player.Height, player.Weight, player.Age)
		if err != nil {
			return err
		}
	}

	fmt.Println("Players added to database.")
	return nil
}

func UpdatePlayers(db database.Database, players []Player) ([]stats.Stats, error) {
	playerTraded := ""
	var stats []stats.Stats
	for _, player := range players {
		if player.TeamAbbr == "TOT" {
			playerTraded = player.ID
			_, err := db.UpdatePlayerAge(player.ID, player.Age)
			if err != nil {
				return nil, err
			}
			stats = append(stats, player.Stats)
			continue
		}

		if playerTraded == player.ID {
			fmt.Println("Update team for ", player.ID, player.TeamAbbr)
			_, err := db.UpdateTeamForPlayer(player.ID, player.TeamAbbr)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			_, err := db.UpdatePlayerAge(player.ID, player.Age)
			if err != nil {
				fmt.Println(err)
			}
			stats = append(stats, player.Stats)
		}
		fmt.Println(player)
	}

	fmt.Println("Players updated in database.")
	return stats, nil
}

func ParsePlayersCurrentSeason() ([]Player, error) {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_game.html#per_game_stats", getEndYearOfTheSeason())
	res, err := parser.SendRequest(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var players []Player
	rows := doc.Find("table#per_game_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		player := fillPlayerInfoForSeason(row)
		player.Stats = stats.FillPlayerStatsForSeason(row, getEndYearOfTheSeason())
		player.Stats.PlayerID = player.ID
		players = append(players, player)
	})

	return players, nil
}

func fillPlayerInfoForSeason(row *goquery.Selection) Player {
	var player Player
	player.Name = row.Find("td[data-stat='player']").Text()
	id, exists := row.Find("td[data-stat='player'] > a").Attr("href")
	if exists {
		idParts := strings.Split(id, "/")
		if len(idParts) > 3 {
			player.ID = strings.TrimSuffix(idParts[3], ".html")
		}
	}
	player.Age, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='age']").Text()), 10, 64)
	player.TeamAbbr = row.Find("td[data-stat='team_id']").Text()
	return player
}

func getEndYearOfTheSeason() string {
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
