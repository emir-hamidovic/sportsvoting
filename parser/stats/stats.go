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
	stats.FGPercentage = stats.FGPercentage * 100
	stats.ThreeFGPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='fg3_pct']").Text()), 64)
	stats.ThreeFGPercentage = stats.ThreeFGPercentage * 100
	stats.FTPercentage, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ft_pct']").Text()), 64)
	stats.FTPercentage = stats.FTPercentage * 100
	stats.Position = row.Find("td[data-stat='pos']").Text()
	stats.Season = season
	return stats
}

func UpdateStats(db database.Database, stats []Stats) error {
	for _, stat := range stats {
		fmt.Println(stat)
		res, err := db.UpdateStats(stat.Games, stat.Minutes, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Turnovers, stat.FGPercentage, stat.FTPercentage, stat.ThreeFGPercentage, stat.Season, stat.Position, stat.PlayerID)
		if err != nil {
			fmt.Println("update stats err")
			fmt.Println(err)
		}

		rows, err := res.RowsAffected()
		if err != nil {
			fmt.Println("rows affected err")
			fmt.Println(err)
		}
		fmt.Println(rows)

		if rows == 0 {
			fmt.Println("Insert stats")
			_, err = db.InsertStats(stat.Games, stat.Minutes, stat.Points, stat.Rebounds, stat.Assists, stat.Steals, stat.Blocks, stat.Turnovers, stat.FGPercentage, stat.FTPercentage, stat.ThreeFGPercentage, stat.Season, stat.Position, stat.PlayerID)
			if err != nil {
				fmt.Println("insert stats err")
				fmt.Println(err)
			}
		}
	}

	fmt.Println("Player stats added to database.")
	return nil
}
