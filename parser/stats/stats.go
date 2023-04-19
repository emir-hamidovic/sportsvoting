package stats

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Stats struct {
	PlayerID          string
	Games             int64
	Minutes           float64
	Points            float64
	Rebounds          float64
	Assists           float64
	Steals            float64
	Blocks            float64
	Turnovers         float64
	FGPercentage      float64
	ThreeFGPercentage float64
	FTPercentage      float64
	Season            string
	Position          string
	TeamAbbr          string
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *Stats) {
	stats.Games, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='g']").Text()), 10, 64)
	stats.Minutes, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='mp_per_g']").Text()), 64)
	stats.Points, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='pts_per_g']").Text()), 64)
	stats.Rebounds, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='trb_per_g']").Text()), 64)
	stats.Assists, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ast_per_g']").Text()), 64)
	stats.Blocks, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='blk_per_g']").Text()), 64)
	stats.Steals, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='stl_per_g']").Text()), 64)
	stats.Turnovers, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='tov_per_g']").Text()), 64)
	stats.FGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg_pct']").Text()), 64)
	stats.FGPercentage = stats.FGPercentage * 100
	stats.ThreeFGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg3_pct']").Text()), 64)
	stats.ThreeFGPercentage = stats.ThreeFGPercentage * 100
	stats.FTPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ft_pct']").Text()), 64)
	stats.FTPercentage = stats.FTPercentage * 100
	stats.Season = season
}

func UpdateStats(db database.Database, stats Stats) error {
	fmt.Println(stats)
	res, err := db.UpdateStats(stats.Games, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if rows == 0 && stats.Season != "" {
		_, err = db.InsertStats(stats.Games, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_game.html", season)
	doc, err := parser.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#per_game_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := row.Find("td[data-stat='team_id']").Text()
		if team == "TOT" {
			id, exists := row.Find("td[data-stat='player'] > a").Attr("href")
			if exists {
				idParts := strings.Split(id, "/")
				if len(idParts) > 3 {
					id = strings.TrimSuffix(idParts[3], ".html")
				}
			}

			var stats Stats
			FillPlayerStatsForSeason(row, season, &stats)
			position := row.Find("td[data-stat='pos']").Text()
			res, err := db.UpdateTradedPlayerStats(stats.Games, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, position, id)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(res.RowsAffected())
		}
	})

	return nil
}
