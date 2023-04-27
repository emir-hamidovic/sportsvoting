package advancedstats

import (
	"fmt"
	"scraper/database"
	"scraper/request"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AdvancedStats struct {
	PlayerID string  `json:"stats,omitempty"`
	TeamAbbr string  `json:"team,omitempty"`
	Season   string  `json:"season,omitempty"`
	PER      float64 `json:"per,omitempty"`
	TSPct    float64 `json:"ts,omitempty"`
	USGPCt   float64 `json:"usg,omitempty"`
	OffWS    float64 `json:"ows,omitempty"`
	DefWS    float64 `json:"dws,omitempty"`
	WS       float64 `json:"ws,omitempty"`
	OffBPM   float64 `json:"obpm,omitempty"`
	DefBPM   float64 `json:"dbpm,omitempty"`
	BPM      float64 `json:"bpm,omitempty"`
	VORP     float64 `json:"vorp,omitempty"`
	DefRtg   float64 `json:"defrtg,omitempty"`
	OffRtg   float64 `json:"offrtg,omitempty"`
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
	// Problem with update, check this
	fmt.Println(stats)
	res, err := db.UpdateAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
	if err != nil {
		fmt.Println(err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rows, stats.Season)

	if rows == 0 && stats.Season != "" {
		fmt.Println("insert advanced", stats)
		_, err = db.InsertAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.OffRtg, stats.DefRtg, stats.TeamAbbr, stats.PlayerID, stats.Season)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("Player advanced stats added to database.")
	return nil
}

func UpdateTradedPlayerStats(db database.Database, season string) error {
	url := fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_advanced.html", season)
	doc, err := request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows := doc.Find("table#advanced_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := row.Find("td[data-stat='team_id']").Text()
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			var stats AdvancedStats
			FillPlayerStatsForSeason(row, season, &stats)
			_, err := db.UpdateTradedPlayerAdvancedStats(stats.PER, stats.TSPct, stats.USGPCt, stats.OffWS, stats.DefWS, stats.WS, stats.OffBPM, stats.DefBPM, stats.BPM, stats.VORP, stats.Season, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	url = fmt.Sprintf("https://www.basketball-reference.com/leagues/NBA_%s_per_poss.html", season)
	doc, err = request.GetDocumentFromURL(url)
	if err != nil {
		return err
	}

	rows = doc.Find("table#per_poss_stats > tbody > tr")
	rows.Each(func(i int, row *goquery.Selection) {
		team := row.Find("td[data-stat='team_id']").Text()
		if team == "TOT" {
			id := request.GetPlayerIDFromDocument(row)

			defrtg, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='def_rtg']").Text()), 64)
			offrtg, _ := strconv.ParseFloat(strings.TrimSpace(row.Find("td[data-stat='off_rtg']").Text()), 64)
			_, err := db.UpdateOffAndDefRtg(offrtg, defrtg, season, id)
			if err != nil {
				fmt.Println(err)
			}
		}
	})

	return nil
}
