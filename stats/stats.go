package stats

import (
	"fmt"
	"scraper/database"
	"scraper/request"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Stats struct {
	PlayerID          string  `json:"playerid,omitempty"`
	Games             int64   `json:"g,omitempty"`
	GamesStarted      int64   `json:"gs,omitempty"`
	Minutes           float64 `json:"mpg,omitempty"`
	Points            float64 `json:"ppg,omitempty"`
	Rebounds          float64 `json:"rpg,omitempty"`
	Assists           float64 `json:"apg,omitempty"`
	Steals            float64 `json:"spg,omitempty"`
	Blocks            float64 `json:"bpg,omitempty"`
	Turnovers         float64 `json:"topg,omitempty"`
	FGPercentage      float64 `json:"fgpct,omitempty"`
	ThreeFGPercentage float64 `json:"threefgpct,omitempty"`
	FTPercentage      float64 `json:"ftpct,omitempty"`
	Season            string  `json:"season,omitempty"`
	Position          string  `json:"position,omitempty"`
	TeamAbbr          string  `json:"team,omitempty"`
	IsRookie          bool    `json:"rookie,omitempty"`
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *Stats) {
	stats.Games, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='g']").Text()), 10, 64)
	stats.GamesStarted, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='gs']").Text()), 10, 64)
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
	res, err := db.UpdateStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if rows == 0 && stats.Season != "" {
		_, err = db.InsertStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID, stats.TeamAbbr)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_game.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#per_game_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := row.Find("td[data-stat='team_id']").Text()
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			var stats Stats
			FillPlayerStatsForSeason(row, season, &stats)
			position := row.Find("td[data-stat='pos']").Text()
			_, err := db.UpdateTradedPlayerStats(stats.Games, stats.GamesStarted, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, position, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	return nil
}

func SetRookies(db database.Database, season string) error {
	fmt.Println("Setting rookies")
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_rookies.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#rookies > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		id := request.GetPlayerIDFromDocument(row)
		fmt.Println(id)
		_, err := db.SetRookieStatus(id)
		if err != nil {
			fmt.Println(err)
		}
	})

	return nil
}
