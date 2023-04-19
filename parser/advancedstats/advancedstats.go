package advancedstats

import (
	"fmt"
	"scraper/database"
	"scraper/parser"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AdvancedStats struct {
	PlayerID string
	TeamAbbr string
	Season   string
	PER      float64
	TSPct    float64
	USGPCt   float64
	OffWS    float64
	DefWS    float64
	WS       float64
	OffBPM   float64
	DefBPM   float64
	BPM      float64
	VORP     float64
}

func FillPlayerStatsForSeason(row *goquery.Selection, season string, stats *AdvancedStats) {
	stats.PER, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='per']").Text()), 64)
	stats.TSPct, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ts_pct']").Text()), 64)
	stats.USGPCt, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='usg_pct']").Text()), 64)
	stats.OffWS, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ows']").Text()), 64)
	stats.DefWS, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='dws']").Text()), 64)
	stats.WS, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='ws']").Text()), 64)
	stats.OffBPM, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='obpm']").Text()), 64)
	stats.DefBPM, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='dbpm']").Text()), 64)
	stats.BPM, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='bpm']").Text()), 64)
	stats.VORP, _ = strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='vorp']").Text()), 64)
	stats.Season = season
}

func UpdateStats(db database.Database, stats AdvancedStats) error {
	fmt.Println(stats)
	res, err := db.UpdateAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.TeamAbbr, stats.PlayerID, stats.Season)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	if rows == 0 && stats.Season != "" {
		_, err = db.InsertAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.TeamAbbr, stats.PlayerID, stats.Season)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player advanced stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_advanced.html", season)
	doc, err := parser.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#advanced_stats > tbody > tr")
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

			var stats AdvancedStats
			FillPlayerStatsForSeason(row, season, &stats)
			res, err := db.UpdateTradedPlayerAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.Season, id)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(res.RowsAffected())
		}
	})

	return nil
}
