package stats

import (
	"fmt"
	"scraper/database"
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
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string) Stats {
	var stats Stats
	stats.Games, _ = strconv.ParseInt(strings.TrimSpace(row.Find("td[data-stat='g']").Text()), 10, 64)
	stats.Minutes, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='mp_per_g']").Text()), 64)
	stats.Points, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='pts_per_g']").Text()), 64)
	stats.Rebounds, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='trb_per_g']").Text()), 64)
	stats.Assists, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ast_per_g']").Text()), 64)
	stats.Blocks, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='blk_per_g']").Text()), 64)
	stats.Steals, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='stl_per_g']").Text()), 64)
	stats.Turnovers, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='tov_per_g']").Text()), 64)
	stats.FGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg_pct']").Text()), 64)
	stats.ThreeFGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg3_pct']").Text()), 64)
	stats.FTPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ft_pct']").Text()), 64)
	stats.Position = row.Find("td[data-stat='pos']").Text()
	stats.Season = season
	return stats
}

func UpdateStats(db database.Database, stats Stats) error {
	res, err := db.UpdateStats(stats.Games, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		_, err = db.InsertStats(stats.Games, stats.Minutes, stats.Points, stats.Rebounds, stats.Assists, stats.Steals, stats.Blocks, stats.Turnovers, stats.FGPercentage, stats.FTPercentage, stats.ThreeFGPercentage, stats.Season, stats.Position, stats.PlayerID)
		if err != nil {
			return err
		}
	}

	fmt.Println("Player stats added to database.")
	return nil
}

func getOldTeamAbbreviation(abbr string) string {
	switch abbr {
	case "NOP":
		return "NOH"
	case "CHO":
		return "CHA"
	case "BRK":
		return "NJN"
	}

	return abbr
}
