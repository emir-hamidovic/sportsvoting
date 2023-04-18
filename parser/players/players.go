package players

import (
	"fmt"
	"scraper/database"
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

func GetCurrentSeasonPerGameStats(doc *goquery.Document, player map[string]Player) (map[string]Player, error) {
	rows := doc.Find("table#per_game > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		pl := getPlayerAge(row)
		if entry, ok := player[pl.ID]; ok {
			entry.Age = pl.Age
			stats.FillPlayerStatsForSeason(row, getEndYearOfTheSeason(), &entry.Stats)
			player[pl.ID] = entry
		}
	})

	return player, nil
}

func getPlayerAge(row *goquery.Selection) Player {
	var player Player
	id, exists := row.Find("td[data-stat='player'] > a").Attr("href")
	if exists {
		idParts := strings.Split(id, "/")
		if len(idParts) > 3 {
			player.ID = strings.TrimSuffix(idParts[3], ".html")
		}
	}
	player.Age, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='age']").Text()), 10, 64)
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
